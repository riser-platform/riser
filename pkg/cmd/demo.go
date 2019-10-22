package cmd

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"regexp"
	"riser/pkg/logger"
	"riser/pkg/rc"
	"riser/pkg/steps"
	"riser/pkg/ui"
	"riser/pkg/ui/style"
	"strings"

	"github.com/pkg/errors"

	"github.com/AlecAivazis/survey/v2"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/shurcooL/httpfs/vfsutil"
	"github.com/spf13/cobra"
)

const ApiKeySizeBytes = 20
const demoStageName = "demo"

func newDemoCommand(config *rc.RuntimeConfiguration, assets http.FileSystem) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demo",
		Short: "Commands for the riser demo",
	}

	cmd.AddCommand(newInstallDemoCommand(config, assets))
	cmd.AddCommand(newDemoStatusCommand(config))

	return cmd
}

func newInstallDemoCommand(config *rc.RuntimeConfiguration, assets http.FileSystem) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Installs a self-contained riser demo to a k8s cluster (minikube recommended)",
		Long:  "Install a self-contained riser demo to a k8s cluster (minikube recommended) along with all required infrastructure (istio, postgresql, etc)",
		Run: func(cmd *cobra.Command, args []string) {
			demoInstall(config, assets)
		},
	}

	return cmd
}

func newDemoStatusCommand(config *rc.RuntimeConfiguration) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Provides status and help information for the riser demo",
		Run: func(cmd *cobra.Command, args []string) {
			demoStatus(config)
		},
	}
}

func outputAssetsToTempDir(assets http.FileSystem, targetDir string) error {
	walkFn := func(assetsPath string, fi os.FileInfo, r io.ReadSeeker, err error) error {
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("can't stat file %q", assetsPath))
		}
		if !fi.IsDir() {
			b, err := ioutil.ReadAll(r)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("can't read file %q", assetsPath))
			}

			targetPath := path.Join(targetDir, assetsPath)
			err = os.MkdirAll(path.Dir(targetPath), 0777)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("can't create dir %q", path.Dir(targetPath)))
			}
			err = ioutil.WriteFile(targetPath, b, 0777)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("can't write file %q", targetPath))
			}
		}
		return nil
	}

	return vfsutil.WalkFiles(assets, "/demo", walkFn)
}

func demoInstall(config *rc.RuntimeConfiguration, assets http.FileSystem) {
	assetDir, err := ioutil.TempDir(os.TempDir(), "riser-demo-installer")
	defer os.RemoveAll(assetDir)
	ui.ExitIfErrorMsg(err, "Error creating temp dir")

	err = outputAssetsToTempDir(assets, assetDir)
	ui.ExitIfErrorMsg(err, "Error writing assets to temp dir")

	demoPath := path.Join(assetDir, "demo")

	_, err = exec.LookPath("kubectl")
	ui.ExitIfErrorMsg(err, "kubectl must exist in path")

	_, err = exec.LookPath("git")
	ui.ExitIfErrorMsg(err, "git must exist in path")

	kcOutput, err := exec.Command("kubectl", "config", "current-context").Output()
	ui.ExitIfErrorMsg(err, fmt.Sprintf("Error getting current kube context. Maybe the current context is not set?"))

	logger.Log().Warn("The riser demo installs infrastructure that may collide with existing infrastructure (e.g. istio). It is highly recommended that you install the demo to an empty cluster (e.g. a new minikube project)")

	useKc := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Are you sure you wish to install the demo to the %q context?", strings.TrimSpace(string(kcOutput))),
	}
	err = survey.AskOne(prompt, &useKc)
	ui.ExitIfError(err)

	if !useKc {
		ui.ExitErrorMsg("Please change to the desired kube context and try again")
	}

	var gitUrl string
	gitUrlPrompt := &survey.Input{
		Message: "Enter the GitHub URL (including auth if private) for the riser state repo.",
		Help: ui.StripNewLines(`
The riser state repo contains all kubernetes state for riser apps and infrastructure. This repo should never hold any secrets, but you may still wish for it to be private.
For private repos it's recommended that you use a Personal Access Token with repo full access. For example: https://oauthtoken:YOUR-TOKEN-HERE@github.com/your/repo.
See https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line for more information on creating a GitHub Personal Access Token
`),
	}

	for {
		err = survey.AskOne(gitUrlPrompt, &gitUrl, survey.WithValidator(func(ans interface{}) error {
			return validation.Validate(ans,
				validation.Required,
				is.URL,
				// This pattern assumes the <host>/<org>/<repoName> format which works with GitHub, BitBucket, and GitLab.
				validation.Match(regexp.MustCompile("https://.+/.+/.+")).Error("Repo URL must be in the format: https://<host>/<org>/<repoName>"))
		}))
		if err == nil {
			break
		}
		if err.Error() == "interrupt" {
			ui.ExitErrorMsg("aborted")
		}
	}

	gitUrlParsed, err := url.Parse(gitUrl)
	// Should never happen since the URL is validated above
	ui.ExitIfError(err)

	gitUrlPassword, _ := gitUrlParsed.User.Password()

	// Riser-server takes the repo URL w/o auth.
	gitUrlNoAuthParsed, _ := url.Parse(gitUrlParsed.String())
	gitUrlNoAuthParsed.User = nil

	gitRepoName := strings.Split(gitUrlParsed.Path, "/")[2]

	logger.Log().Info("Installing demo...")
	getApiKeyFromRiserSecretStep := steps.NewShellExecStep("Check for existing Riser API key",
		"kubectl get secret riser-server -n riser-system -o jsonpath='{.data.RISER_BOOTSTRAP_APIKEY}' || echo ''")
	apiKeyGenStep := steps.NewExecStep("Generate Riser API key", exec.Command("riser", "ops", "generate-apikey"))
	err = steps.Run(
		steps.NewExecStep("Validate Git remote", exec.Command("git", "ls-remote", gitUrlParsed.String(), "HEAD")),
		// Install namespaces and istio CRDs separately due to ordering issues (declarative infra... not quite!)
		steps.NewExecStep("Apply namespaces and CRDs", exec.Command("kubectl", "apply",
			"-f", path.Join(demoPath, "kube-resources/riser-server/namespaces.yaml"),
			"-f", path.Join(demoPath, "kube-resources/istio/0_namespace.yaml"),
			"-f", path.Join(demoPath, "kube-resources/istio/1_init.yaml"),
		)),
		steps.NewRetryStep(
			func() steps.Step {
				// We don't wait for each specific CRD. In testing we've found these two are the most common ones that aren't immediately ready
				// May have to adjust over time.
				return steps.NewShellExecStep("Wait for CRDs",
					"kubectl wait --for condition=established crd/gateways.networking.istio.io && kubectl wait --for condition=established crd/clusterissuers.certmanager.k8s.io")
			},
			120,
			func(stepErr error) bool {
				return strings.Contains(stepErr.Error(), "Error from server (NotFound)")
			}),
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

	// Run another group of steps since we rely on the state of previous steps (step runner could support deferred state but this is simpler for now)
	err = steps.Run(
		steps.NewShellExecStep("Create riser-server configuration",
			"kubectl create configmap riser-server --namespace=riser-system "+
				fmt.Sprintf("--from-literal=RISER_GIT_URL=%s", gitUrlNoAuthParsed.String())+
				" --dry-run=true -o yaml | kubectl apply -f -"),
		steps.NewShellExecStep("Create secret for riser-server",
			"kubectl create secret generic riser-server --namespace=riser-system "+
				fmt.Sprintf("--from-literal=RISER_BOOTSTRAP_APIKEY=%s ", apiKey)+
				fmt.Sprintf("--from-literal=RISER_GIT_USERNAME=%s ", gitUrlParsed.User.Username())+
				fmt.Sprintf("--from-literal=RISER_GIT_PASSWORD=%s ", gitUrlPassword)+
				"--from-literal=RISER_POSTGRES_USERNAME=riseradmin "+
				"--from-literal=RISER_POSTGRES_PASSWORD=riserpw "+
				" --dry-run=true -o yaml | kubectl apply -f -"),
		// TODO: See idempotent message from riser-server secret
		steps.NewShellExecStep("Create secret for riser-controller",
			"kubectl create secret generic riser-controller --namespace=riser-system "+
				fmt.Sprintf("--from-literal=RISER_SERVER_APIKEY=%s ", apiKey)+
				" --dry-run=true -o yaml | kubectl apply -f -"),
		steps.NewExecStep("Apply static demo resources", exec.Command("kubectl", "apply", "-R", "-f", path.Join(demoPath, "kube-resources"))),
		steps.NewShellExecStep("Create secret for kube-applier",
			"kubectl create secret generic kube-applier --namespace=kube-applier "+
				fmt.Sprintf("--from-literal=GIT_SYNC_REPO=%s", gitUrlParsed.String())+
				" --dry-run=true -o yaml | kubectl apply -f -"),
		steps.NewShellExecStep("Create kube-applier configuration",
			"kubectl create configmap kube-applier --namespace kube-applier "+
				fmt.Sprintf("--from-literal=REPO_PATH=/git-repo/%s/stages/%s/kube-resources", gitRepoName, demoStageName)+
				" --dry-run=true -o yaml | kubectl apply -f -"),
		steps.NewFuncStep(fmt.Sprintf("Save riser context %q", demoStageName),
			func() error {
				secure := false
				demoContext := &rc.Context{Name: demoStageName, ServerURL: "https://api.demo.riser", Apikey: apiKey, Secure: &secure}
				config.SaveContext(demoContext)
				return rc.SaveRc(config)
			}),
	)

	ui.ExitIfError(err)

	logger.Log().Info(style.Good("Installation Complete!"))
	logger.Log().Info("Executing \"riser demo status\"...")

	demoStatus(config)
}

func demoStatus(config *rc.RuntimeConfiguration) {
	logger.Log().Warn(`If you're using minikube be sure that "minikube tunnel" is running.`)
	err := config.SetCurrentContext(demoStageName)
	ui.ExitIfErrorMsg(err, "Error loading demo config. Please run \"riser demo install\".")

	err = steps.Run(
		steps.NewRetryStep(func() steps.Step {
			return steps.NewShellExecStep(
				"Check riser-server status (this could take a few minutes after installation)",
				"kubectl get po riser-server-0 -n riser-system -o jsonpath='{.status.conditions[?(@.type==\"Ready\")].status}' | grep True")
		},
			300,
			steps.AlwaysRetry()))

	if err != nil {
		logger.Log().Error(err.Error())
		ui.ExitErrorMsg(`Tips:
• On slower systems this can take longer than expected. You may try running \"riser demo status\" again.
• Ensure that your kubernetes context is set to the cluster with the demo installed.
• Ensure that the riser demo is installed using "riser demo install".
• Check the pod logs for pods in the "riser-system" namespace.
		`)
	}

	ingressGatewayStep := steps.NewRetryStep(func() steps.Step {
		return steps.NewShellExecStep(
			"Check Istio ingress gateway",
			"kubectl get service istio-ingressgateway -n istio-system -o jsonpath='{.status.loadBalancer.ingress[0].ip}' | grep ^")
	},
		120,
		// If there's an error always retry - Could improve this as obviously there are some errors where we'd want to abort e.g. kubernetes is not reachable
		steps.AlwaysRetry())

	err = steps.Run(ingressGatewayStep)

	if err != nil {
		logger.Log().Error(err.Error())
		ui.ExitErrorMsg(`Tips:
• If you're using minikube be sure that "minikube tunnel" is running.
• Ensure that your kubernetes context is set to the cluster with the demo installed.
• Ensure that the riser demo is installed using "riser demo install".
• Check the service status and pod logs for "istio-ingressgateway" in the "istio-system" namespace.
		`)
	}

	ingressIp := ui.StripNewLines(ingressGatewayStep.State("stdout").(string))

	logger.Log().Info("\n" + style.Good("🚀 Everything checks out!") + "\n")

	logger.Log().Info("Gateway IP:\t" + style.Emphasis(ingressIp))
	logger.Log().Info("API Host:\t" + style.Emphasis("api.demo.riser"))
	logger.Log().Info("Apps host:\t" + style.Emphasis("*.apps.demo.riser"))

	logger.Log().Info("\nInstructions:")
	logger.Log().Info(fmt.Sprintf("• In your hosts file (e.g. /etc/hosts on OSX) or local DNS server set the IP for the host %s to the ingress IP: %s", style.Emphasis("api.demo.riser"), style.Emphasis(ingressIp)))
	logger.Log().Info(fmt.Sprintf("  Example /etc/hosts entry:\n  %s", style.Muted(fmt.Sprintf("%s api.demo.riser", ingressIp))))
	logger.Log().Info("• For easier access to your apps, you may wish to add additional host entries for each app using the format <YOUR-APP>.apps.demo.riser to the same ingress IP, or create a wildcard DNS record for *.apps.demo.riser.")
	logger.Log().Info("• You may also access your apps by passing a host header. For example, with curl:")
	logger.Log().Info(style.Muted(fmt.Sprintf("  curl -k -H \"Host: <YOUR-APP>.apps.demo.riser\" https://%s", ingressIp)))
	logger.Log().Info(`• Try out the testdummy app!
  - In an empty folder create the app with a default config using "riser apps init testdummy"
  - Edit "app.yaml" and specify "tshak/testdummy" as the docker image
  - Deploy using "riser deploy latest demo"
  - Use "riser status" to check the status of your deployment
  - Once deployed access using "curl -k https://testdummy.apps.demo.riser". If all went well you should receive a HTTP 200 response with the text "pong".`)
	// TODO: Link to docs for further info

	logger.Log().Info("\nExecute \"riser demo status\" to see this message again.")
}
