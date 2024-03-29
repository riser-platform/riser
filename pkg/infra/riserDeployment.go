package infra

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"riser/pkg/rc"
	"riser/pkg/steps"
	"riser/pkg/ui"

	"github.com/pkg/errors"
	"github.com/shurcooL/httpfs/vfsutil"
)

const (
	DefaultEnvironmentName = "demo"
	DefaultServerImage     = "ghcr.io/riser-platform/riser-server:v0.0.47"
	DefaultControllerImage = "ghcr.io/riser-platform/riser-controller:v0.0.19"
)

type RiserDeployment struct {
	Assets      fs.FS
	GitUrl      string
	GitBranch   string
	RiserConfig *rc.RuntimeConfiguration
	// Optional
	GitSSHKeyPath   string
	EnvironmentName string
	KubeDeployer    KubeDeployer
	ServerImage     string
	ControllerImage string
}

func NewRiserDeployment(assets fs.FS, riserConfig *rc.RuntimeConfiguration, gitUrl string, gitBranch string) *RiserDeployment {
	return &RiserDeployment{
		Assets:          assets,
		RiserConfig:     riserConfig,
		GitUrl:          gitUrl,
		GitBranch:       gitBranch,
		EnvironmentName: DefaultEnvironmentName,
		KubeDeployer:    &NoopDeployer{},
		ServerImage:     DefaultServerImage,
		ControllerImage: DefaultControllerImage,
	}
}

// Deploy deploys Riser k8s manifests for the demo and for e2e tests.
func (deployment *RiserDeployment) Deploy() error {
	templateVars := map[string]string{
		"RISER_SERVER_IMAGE":     deployment.ServerImage,
		"RISER_CONTROLLER_IMAGE": deployment.ControllerImage,
	}
	assetPath, err := outputDeployAssetsToTempDir(http.FS(deployment.Assets), templateVars)
	if err != nil {
		return errors.Wrap(err, "Error writing assets to temp dir")
	}
	defer os.RemoveAll(assetPath)

	err = deployment.KubeDeployer.Deploy()
	ui.ExitIfErrorMsg(err, "Error deploying kubernetes")

	getApiKeyFromRiserSecretStep := steps.NewShellExecStep("Check for existing Riser API key",
		"kubectl get secret riser-server -n riser-system -o jsonpath='{.data.RISER_BOOTSTRAP_APIKEY}' || echo ''")
	apiKeyGenStep := steps.NewExecStep("Generate Riser API key", exec.Command("riser", "ops", "generate-apikey"))
	err = steps.Run(
		// Install namespaces and some CRDs separately due to ordering issues (declarative infra... not quite!)
		steps.NewExecStep("Apply prerequisites", exec.Command("kubectl", "apply",
			"-f", path.Join(assetPath, "kube-resources/riser-server/namespaces.yaml"),
			"-f", path.Join(assetPath, "flux/namespace.yaml"),
		)),
		steps.NewExecStep("Validate Git remote", exec.Command("git", "ls-remote", deployment.GitUrl, "HEAD")),
		getApiKeyFromRiserSecretStep,
	)
	ui.ExitIfError(err)

	var apiKey string

	// If we've already bootstrapped the API key, get that key instead of generating a new one.
	if getApiKeyFromRiserSecretStep.State("stdout") != "" {
		apiKeyBytes, err := base64.StdEncoding.DecodeString(getApiKeyFromRiserSecretStep.State("stdout").(string))
		if err == nil {
			apiKey = string(apiKeyBytes)
		}
	}

	if apiKey == "" {
		ui.ExitIfError(steps.Run(apiKeyGenStep))
		apiKey = apiKeyGenStep.State("stdout").(string)
	}

	gitDeployKeyArg := ""
	if deployment.GitSSHKeyPath != "" {
		gitDeployKeyArg = fmt.Sprintf("--from-file=identity=%s", deployment.GitSSHKeyPath)
	}

	// Run another group of steps since we rely on the state of previous steps (step runner could support deferred state but this is simpler for now)
	err = steps.Run(
		steps.NewShellExecStep("Configure riser-server",
			"kubectl create configmap riser-server --namespace=riser-system --from-literal=RISER_GIT_SSH_KEY_PATH=/etc/riser/ssh/identity "+
				"--dry-run=client -o yaml | kubectl apply -f -"),
		steps.NewShellExecStep("Create secret for riser-server",
			"kubectl create secret generic riser-server --namespace=riser-system "+
				fmt.Sprintf("--from-literal=RISER_BOOTSTRAP_APIKEY=%s ", apiKey)+
				fmt.Sprintf("--from-literal=RISER_GIT_URL=%s ", deployment.GitUrl)+
				"--from-literal=RISER_POSTGRES_USERNAME=riseradmin "+
				"--from-literal=RISER_POSTGRES_PASSWORD=riserpw "+
				" --dry-run=client -o yaml | kubectl apply -f -"),
		// Empty secret must exist since there's a volume mount that expects it
		steps.NewShellExecStep("Create secret for git",
			fmt.Sprintf("kubectl create secret generic riser-git-deploy %s --namespace=riser-system --dry-run=client -o yaml | kubectl apply -f -", gitDeployKeyArg)),
		steps.NewShellExecStep("Configure riser-controller", fmt.Sprintf(
			`kubectl create configmap riser-controller --namespace=riser-system \
					--from-literal=RISER_SERVER_URL=http://riser-server.riser-system.svc.cluster.local  \
					--from-literal=RISER_ENVIRONMENT=%s \
					--dry-run=client -o yaml | kubectl apply -f -`, deployment.EnvironmentName)),
		steps.NewShellExecStep("Create secret for riser-controller",
			"kubectl create secret generic riser-controller --namespace=riser-system "+
				fmt.Sprintf("--from-literal=RISER_SERVER_APIKEY=%s ", apiKey)+
				" --dry-run=client -o yaml | kubectl apply -f -"),
		steps.NewShellExecStep("Create secret for flux",
			"kubectl create secret generic flux-git --namespace=flux "+
				fmt.Sprintf("--from-literal=GIT_URL=%s ", deployment.GitUrl)+
				fmt.Sprintf("--from-literal=GIT_BRANCH=%s ", deployment.GitBranch)+
				"--from-literal=GIT_PATH=state "+
				" --dry-run=client -o yaml | kubectl apply -f -"),
		steps.NewShellExecStep("Create flux git key secret",
			fmt.Sprintf("kubectl create secret generic flux-git-deploy %s --namespace=flux --dry-run=client -o yaml | kubectl apply -f -", gitDeployKeyArg)),
		steps.NewShellExecStep("Install flux",
			fmt.Sprintf("kubectl apply -f %s --namespace flux", path.Join(assetPath, "flux"))),
		steps.NewRetryStep(
			func() steps.Step {
				return steps.NewExecStep("Apply other resources", exec.Command("kubectl", "apply", "-R", "-f", path.Join(assetPath, "kube-resources")))
			},
			180,
			// Eventually CRDs will converge. Hopefully kubectl apply will handle this in the future.
			steps.AlwaysRetry(),
		),
		steps.NewFuncStep(fmt.Sprintf("Save riser context %q", deployment.EnvironmentName),
			func() error {
				secure := false
				newRiserContext := &rc.Context{
					Name:      deployment.EnvironmentName,
					ServerURL: "https://riser-server.riser-system.demo.riser",
					Apikey:    apiKey,
					Secure:    &secure}
				deployment.RiserConfig.SetContext(newRiserContext)
				return rc.SaveRc(deployment.RiserConfig)
			}),
		steps.NewShellExecStep("Wait for riser-server", "kubectl wait --for=condition=ready --timeout=120s ksvc/riser-server -n riser-system"),
		// This allows for faster e2e runs as the crash backoff is too slow. Init containers could be used here too.
		steps.NewShellExecStep("Install riser-controller",
			fmt.Sprintf("kubectl apply -f %s", path.Join(assetPath, "riser-controller"))),
		steps.NewShellExecStep("Wait for riser-controller", "kubectl wait --for=condition=available --timeout=120s deployment/riser-controller-manager -n riser-system"),
	)

	return err
}

func outputDeployAssetsToTempDir(assets http.FileSystem, templateVars map[string]string) (assetPath string, err error) {
	baseDir, err := ioutil.TempDir(os.TempDir(), "riser-deploy")
	if err != nil {
		return "", errors.Wrap(err, "Error creating temp dir")
	}

	walkFn := func(assetsPath string, fi os.FileInfo, r io.ReadSeeker, err error) error {
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("can't stat file %q", assetsPath))
		}
		if !fi.IsDir() {
			b, err := ioutil.ReadAll(r)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("can't read file %q", assetsPath))
			}

			targetDir := path.Join(baseDir, assetsPath)
			err = os.MkdirAll(path.Dir(targetDir), 0777)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("can't create dir %q", path.Dir(targetDir)))
			}

			if path.Ext(targetDir) == ".yaml" {
				b = []byte(os.Expand(string(b), func(v string) string {
					return templateVars[v]
				}))
			}
			err = ioutil.WriteFile(targetDir, b, 0777)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("can't write file %q", targetDir))
			}
		}
		return nil
	}

	return baseDir + "/deploy", vfsutil.WalkFiles(assets, "/deploy", walkFn)
}
