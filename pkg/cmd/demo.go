package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"riser/pkg/infra"
	"riser/pkg/logger"
	"riser/pkg/rc"
	"riser/pkg/steps"
	"riser/pkg/ui"
	"riser/pkg/ui/style"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/spf13/cobra"
)

const ApiKeySizeBytes = 20
const demoEnvironmentName = "demo"

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
		Long:  "Install a self-contained riser demo to a k8s cluster (minikube recommended) along with all required infrastructure (knative, istio, postgresql, etc)",
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

func demoInstall(config *rc.RuntimeConfiguration, assets http.FileSystem) {
	_, err := exec.LookPath("kubectl")
	ui.ExitIfErrorMsg(err, "kubectl must exist in path")

	_, err = exec.LookPath("git")
	ui.ExitIfErrorMsg(err, "git must exist in path")

	kcOutput, err := exec.Command("kubectl", "config", "current-context").Output()
	ui.ExitIfErrorMsg(err, "Error getting current kube context. Maybe the current context is not set?")

	logger.Log().Warn("The riser demo installs infrastructure that may collide with existing infrastructure. It is highly recommended that you install the demo into an empty Kubernetes cluster.")

	useKc := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Are you sure you wish to install to the %q kubernetes context?", strings.TrimSpace(string(kcOutput))),
	}
	err = survey.AskOne(prompt, &useKc)
	ui.ExitIfError(err)

	if !useKc {
		ui.ExitErrorMsg("Please change to the desired kubernetes context and try again")
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

	logger.Log().Info("Installing demo")

	err = infra.Deploy(&infra.DeployConfig{
		Assets:          assets,
		GitUrl:          gitUrlParsed,
		EnvironmentName: "demo",
		RiserConfig:     config,
	})
	ui.ExitIfError(err)

	logger.Log().Info(style.Good("Installation Complete!"))
	logger.Log().Info("Executing \"riser demo status\"...")

	demoStatus(config)
}

func demoStatus(config *rc.RuntimeConfiguration) {
	logger.Log().Warn(`If you're using minikube be sure that "minikube tunnel" is running.`)
	err := config.SetCurrentContext(demoEnvironmentName)
	ui.ExitIfErrorMsg(err, "Error loading demo config. Please run \"riser demo install\".")

	ingressGatewayStep := steps.NewRetryStep(func() steps.Step {
		return steps.NewShellExecStep(
			"Check istio ingress gateway",
			"kubectl get service istio-ingressgateway -n istio-system -o jsonpath='{.status.loadBalancer.ingress[0].ip}' | grep ^")
	},
		60,
		steps.AlwaysRetry())

	err = steps.Run(ingressGatewayStep)

	if err != nil {
		logger.Log().Error(err.Error())
		ui.ExitErrorMsg(`Tips:
• If you're using minikube be sure that "minikube tunnel" is running and run "riser demo status" again.
• Ensure that your kubernetes context is set to the cluster with the demo installed.
• Check the service status and pod logs for "istio-ingressgateway" in the "istio-system" namespace.
		`)
	}

	err = steps.Run(
		steps.NewRetryStep(func() steps.Step {
			return steps.NewShellExecStep(
				"Check riser-server status (this could take a few minutes after installation)",
				"kubectl get ksvc riser-server -n riser-system -o jsonpath=\"{.status.conditions[?(@.type==\\\"Ready\\\")].status}\" | grep True")
		},
			300,
			func(err error) bool {
				// Abort right away if we can't reach the kube API
				return !regexp.MustCompile("The connection to the server .+ was refused").MatchString(err.Error())
			}))

	if err != nil {
		logger.Log().Error(err.Error())
		ui.ExitErrorMsg(`Tips:
• On slower systems this can take longer than expected right after an installation. You may try running "riser demo status" again after a few minutes.
• Ensure that your kubernetes context is set to the cluster with the demo installed.
• Ensure that the riser demo is installed using "riser demo install".
• Ensure that riser is set to the "demo" context using "riser context current demo"
		`)
	}

	gatewayIp := ui.StripNewLines(ingressGatewayStep.State("stdout").(string))

	logger.Log().Info("\n" + style.Good("🚀 Everything checks out!") + "\n")

	logger.Log().Info("Gateway IP:\t" + style.Emphasis(gatewayIp))
	logger.Log().Info("API Host:\t" + style.Emphasis("riser-server.riser-system.demo.riser"))
	logger.Log().Info("Apps host:\t" + style.Emphasis("*.apps.demo.riser"))

	logger.Log().Info("\nInstructions:")
	logger.Log().Info(fmt.Sprintf("• In your hosts file (e.g. /etc/hosts on OSX) or local DNS server set the IP for the host %s to the gateway IP: %s", style.Emphasis("riser-server.riser-system.demo.riser"), style.Emphasis(gatewayIp)))
	logger.Log().Info(fmt.Sprintf("  Example /etc/hosts entry:\n  %s", style.Muted(fmt.Sprintf("%s riser-server.riser-system.demo.riser", gatewayIp))))
	logger.Log().Info("• For easier access to your apps, you may wish to add additional host entries for each app using the format <YOUR-APP>.apps.demo.riser to the same gateway IP, or create a wildcard DNS record for *.apps.demo.riser.")
	logger.Log().Info("• You may also access your apps by passing a host header. For example, with curl:")
	logger.Log().Info(style.Muted(fmt.Sprintf("  curl -k -H \"Host: <YOUR-APP>.apps.demo.riser\" https://%s", gatewayIp)))
	logger.Log().Info(`• Try out the testdummy app!
  - In an empty folder create the app with a default config using "riser apps init testdummy"
  - Edit "app.yaml" and specify "tshak/testdummy" as the docker image
  - Deploy using "riser deploy latest demo"
  - Use "riser status" to check the status of your deployment
  - Once deployed access using "curl -k https://testdummy.apps.demo.riser". If all went well you should receive a HTTP 200 response with the text "pong".`)
	// TODO: Link to docs for further info

	logger.Log().Info("\nExecute \"riser demo status\" to see this message again.")
}
