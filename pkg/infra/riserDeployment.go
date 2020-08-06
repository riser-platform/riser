package infra

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"riser/pkg/rc"
	"riser/pkg/steps"
	"riser/pkg/ui"
	"strings"

	"github.com/pkg/errors"
	"github.com/shurcooL/httpfs/vfsutil"
)

const (
	DefaultEnvironmentName = "demo"
	DefaultServerImage     = "riserplatform/riser-server:0.0.37"
	DefaultControllerImage = "riserplatform/riser-controller:0.0.13"
)

type RiserDeployment struct {
	Assets        http.FileSystem
	GitUrl        string
	GitSSHKeyPath string
	RiserConfig   *rc.RuntimeConfiguration
	// Optional
	EnvironmentName string
	KubeDeployer    KubeDeployer
	ServerImage     string
	ControllerImage string
}

func NewRiserDeployment(assets http.FileSystem, riserConfig *rc.RuntimeConfiguration, gitUrl string) *RiserDeployment {
	return &RiserDeployment{
		Assets:          assets,
		RiserConfig:     riserConfig,
		GitUrl:          gitUrl,
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
	assetPath, err := outputDeployAssetsToTempDir(deployment.Assets, templateVars)
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
		steps.NewExecStep("Validate Git remote", exec.Command("git", "ls-remote", deployment.GitUrl, "HEAD")),
		// Install namespaces and some CRDs separately due to ordering issues (declarative infra... not quite!)
		steps.NewExecStep("Apply prerequisites", exec.Command("kubectl", "apply",
			"-f", path.Join(assetPath, "kube-resources/istio/istio_operator.yaml"),
			"-f", path.Join(assetPath, "kube-resources/riser-server/namespaces.yaml"),
			"-f", path.Join(assetPath, "knative/serving-crds.yaml"),
			"-f", path.Join(assetPath, "cert-manager/cert-manager.yaml"),
			"-f", path.Join(assetPath, "flux/namespace.yaml"),
		)),
		steps.NewRetryStep(
			func() steps.Step {
				// We don't wait for each specific CRD. In testing we've found these are the most common ones that aren't immediately ready
				// May have to adjust over time.
				return steps.NewShellExecStep("Wait for CRDs",
					`kubectl wait --for condition=established crd/clusterissuers.cert-manager.io && \
					kubectl wait --for condition=established crd/istiooperators.install.istio.io
					`)
			},
			120,
			func(stepErr error) bool {
				return strings.Contains(stepErr.Error(), "Error from server (NotFound)")
			}),
		getApiKeyFromRiserSecretStep,
		// Knative installation is very order dependant, must install istio first.
		steps.NewExecStep("Apply istio resources", exec.Command("kubectl", "apply",
			"-f", path.Join(assetPath, "kube-resources/istio"),
		)),
		steps.NewRetryStep(func() steps.Step {
			return steps.NewShellExecStep("Wait for istio", "kubectl get deployment istiod -n istio-system -o jsonpath='{.status.availableReplicas}' | grep ^1$")
		},
			180, steps.AlwaysRetry()),
		steps.NewExecStep("Apply knative resources", exec.Command("kubectl", "apply", "-R", "-f", path.Join(assetPath, "knative"))),
		// Due to race condition with applying ksvc too early: https://github.com/knative/serving/issues/7576
		// Also, on slower systems the cert-manager-webhook may not be available yet which causes problems down the line
		steps.NewShellExecStep("Wait for knative",
			`kubectl wait --timeout=300s --for=condition=available --namespace knative-serving deployment --all && \
			kubectl wait --for condition=available deployment/cert-manager-webhook -n cert-manager`),
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
					--from-literal=RISER_SERVER_URL=http://riser-server.riser-system.svc.cluster.local \
					--from-literal=RISER_ENVIRONMENT=%s \
					--dry-run=client -o yaml | kubectl apply -f -`, deployment.EnvironmentName)),
		steps.NewShellExecStep("Create secret for riser-controller",
			"kubectl create secret generic riser-controller --namespace=riser-system "+
				fmt.Sprintf("--from-literal=RISER_SERVER_APIKEY=%s ", apiKey)+
				" --dry-run=client -o yaml | kubectl apply -f -"),
		steps.NewShellExecStep("Create secret for flux",
			"kubectl create secret generic flux-git --namespace=flux "+
				fmt.Sprintf("--from-literal=GIT_URL=%s ", deployment.GitUrl)+
				fmt.Sprintf("--from-literal=GIT_PATH=state/%s ", deployment.EnvironmentName)+
				" --dry-run=client -o yaml | kubectl apply -f -"),
		steps.NewShellExecStep("Create flux git key secret",
			fmt.Sprintf("kubectl create secret generic flux-git-deploy %s --namespace=flux --dry-run=client -o yaml | kubectl apply -f -", gitDeployKeyArg)),
		steps.NewShellExecStep("Install flux",
			fmt.Sprintf("kubectl apply -f %s --namespace flux", path.Join(assetPath, "flux"))),
		steps.NewExecStep("Apply other resources", exec.Command("kubectl", "apply", "--validate=false", "-R", "-f", path.Join(assetPath, "kube-resources"))),
		steps.NewFuncStep(fmt.Sprintf("Save riser context %q", deployment.EnvironmentName),
			func() error {
				secure := false
				newRiserContext := &rc.Context{Name: deployment.EnvironmentName, ServerURL: "https://riser-server.riser-system.demo.riser", Apikey: apiKey, Secure: &secure}
				deployment.RiserConfig.SaveContext(newRiserContext)
				return rc.SaveRc(deployment.RiserConfig)
			}),
		steps.NewShellExecStep("Wait for riser-server", "kubectl wait --for=condition=ready --timeout=120s ksvc/riser-server -n riser-system"),
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
