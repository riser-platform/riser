package cmd

import (
	"fmt"
	"net/url"
	"os/exec"
	"path"
	"riser/logger"
	"riser/steps"
	"riser/ui"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/spf13/cobra"
)

const ApiKeySizeBytes = 20
const demoStageName = "demo"

func newDemoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demo",
		Short: "Commands for the riser demo",
	}

	cmd.AddCommand(newInstallDemoCommand())
	return cmd
}

func newInstallDemoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Installs a self-contained riser demo to a k8s cluster (minikube recommended)",
		Long:  "Install a self-contained riser demo to a k8s cluster (minikube recommended) along with all required infrastructure (istio, postgresql, etc)",
		Run: func(cmd *cobra.Command, args []string) {
			installDemo()
		},
	}

	return cmd
}

func installDemo() {
	// TODO: Bake resources into riser binary
	demoPath := "assets/demo"

	_, err := exec.LookPath("kubectl")
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
		Message: "Enter the GitHub URL (including auth) for the riser state repo.",
		Help: `
It's recommended that you use a Personal Access Token with repo full access. For example: https://oauthtoken:YOUR-TOKEN-HERE@github.com/your/repo.
See https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line for more information on creating a GitHub Personal Access Token
`,
	}
	err = survey.AskOne(gitUrlPrompt, &gitUrl)
	ui.ExitIfError(err)

	err = validation.Validate(&gitUrl, is.URL, validation.Required)
	// TODO: reprompt instead of exit on validation failure
	ui.ExitIfError(err)

	gitUrlParsed, err := url.Parse(gitUrl)
	ui.ExitIfErrorMsg(err, "Unable to parse git URL")
	gitUrlPassword, hasPassword := gitUrlParsed.User.Password()
	if !hasPassword {
		ui.ExitErrorMsg("The repo URL must include a password or Personal Access Token. For example: https://oauthtoken:YOUR-TOKEN-HERE@github.com/your/repo")
	}

	// Riser-server takes the repo URL w/o auth.
	gitUrlNoAuthParsed, _ := url.Parse(gitUrlParsed.String())
	gitUrlNoAuthParsed.User = nil

	gitRepoName := strings.Split(gitUrlParsed.Path, "/")[2]

	logger.Log().Info("Installing demo...")
	apiKeyGenStep := steps.NewExecStep("Generate Riser API key", exec.Command("riser", "ops", "generate-apikey"))
	err = steps.Run(
		steps.NewExecStep("Validate Git remote", exec.Command("git", "ls-remote", gitUrlParsed.String(), "HEAD")),
		// Install namespaces and istio CRDs separately due to ordering issues (declarative infra... not quite!)
		steps.NewExecStep("Apply namespaces and CRDs", exec.Command("kubectl", "apply",
			"-f", path.Join(demoPath, "kube-resources/riser-server/namespaces.yaml"),
			"-f", path.Join(demoPath, "kube-resources/istio/0_namespace.yaml"),
			"-f", path.Join(demoPath, "kube-resources/istio/1_init.yaml"),
		)),
		steps.NewWaitStep(
			steps.NewExecStep("Wait for istio CRDs",
				exec.Command("kubectl", "wait", "--for", "condition=established", "crd/gateways.networking.istio.io")),
			10,
			func(stepErr error) bool {
				return strings.Contains(stepErr.Error(), "Error from server (NotFound)")
			}),
		apiKeyGenStep,
	)

	ui.ExitIfError(err)

	// Run another group of steps since we rely on the state of previous steps (step runner could support deffered state but this is simpler for now)
	err = steps.Run(
		steps.NewShellExecStep("Create riser-server configuration",
			"kubectl create configmap riser-server --namespace=riser-system "+
				fmt.Sprintf("--from-literal=RISER_GIT_URL=%s", gitUrlNoAuthParsed.String())+
				" --dry-run=true -o yaml | kubectl apply -f -"),
		//TODO: This is not  idempotent as a new APIKEY will be generated but only the bootstrapped key will be used.
		// While this does not affect the server, the installer will overwrite the riser login with the wrong APIKEY
		steps.NewShellExecStep("Create secret for riser-server",
			"kubectl create secret generic riser-server --namespace=riser-system "+
				fmt.Sprintf("--from-literal=RISER_BOOTSTRAP_APIKEY=%s ", apiKeyGenStep.State("stdout"))+
				fmt.Sprintf("--from-literal=RISER_GIT_USERNAME=%s ", gitUrlParsed.User.Username())+
				fmt.Sprintf("--from-literal=RISER_GIT_PASSWORD=%s ", gitUrlPassword)+
				"--from-literal=RISER_POSTGRES_USERNAME=riseradmin "+
				"--from-literal=RISER_POSTGRES_PASSWORD=riserpw "+
				" --dry-run=true -o yaml | kubectl apply -f -"),
		steps.NewExecStep("Apply all demo resources", exec.Command("kubectl", "apply", "-R", "-f", path.Join(demoPath, "kube-resources"))),
		steps.NewShellExecStep("Create secret for kube-applier",
			"kubectl create secret generic kube-applier --namespace=kube-applier "+
				fmt.Sprintf("--from-literal=GIT_SYNC_REPO=%s", gitUrlParsed.String())+
				" --dry-run=true -o yaml | kubectl apply -f -"),
		steps.NewShellExecStep("Create kube-applier configuration",
			"kubectl create configmap kube-applier --namespace kube-applier "+
				fmt.Sprintf("--from-literal=REPO_PATH=/git-repo/%s/stages/%s/kube-resources", gitRepoName, demoStageName)+
				" --dry-run=true -o yaml | kubectl apply -f -"),
	)

	ui.ExitIfError(err)
}
