package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"riser/pkg/config"
	"riser/pkg/infra"
	"riser/pkg/logger"
	"riser/pkg/rc"
	"riser/pkg/steps"
	"riser/pkg/ui"
	"riser/pkg/ui/style"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
	giturls "github.com/whilp/git-urls"
)

const ApiKeySizeBytes = 20
const demoEnvironmentName = "demo"
const kubectlVersionConstraint = ">=1.18"

type kubectlVersion struct {
	ClientVersion kubectlClientVersion `json:"clientVersion"`
}

type kubectlClientVersion struct {
	GitVersion string `json:"gitVersion"`
}

func newDemoCommand(config *rc.RuntimeConfiguration, assets fs.FS) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demo",
		Short: "Commands for the riser demo",
	}

	cmd.AddCommand(newInstallDemoCommand(config, assets))
	cmd.AddCommand(newDemoStatusCommand(config))
	cmd.AddCommand(newCurlCommand(config))

	return cmd
}

func newCurlCommand(runtimeConfig *rc.RuntimeConfiguration) *cobra.Command {
	var deploymentName string
	var appFilePath string
	cmd := &cobra.Command{
		Use:   "curl [path]",
		Short: "Generates a curl command for access to demo deployments",
		Long:  "Generates a curl command for access to demo deployments. This is helpful as typical demo environments don't have DNS servers or valid certificates, requiring more complicated curl commands.",
		Example: `
All examples assume that you are in the same directory of your app.yaml.

Print the curl command for the default deployment of your app:
  riser demo curl
curl the default deployment for your app:
  riser demo curl | sh
curl the /version path:
  riser demo curl /version | sh
curl a named deployment:
  riser demo curl --name mydeployment | sh`,
		Run: func(cmd *cobra.Command, args []string) {
			path := "/"

			if len(args) > 0 {
				path = args[0]
			}

			if runtimeConfig.CurrentContextName != demoEnvironmentName {
				ui.ExitErrorMsg(fmt.Sprintf("This command is only supported with the Riser context %q", demoEnvironmentName))
			}

			riserContext, err := runtimeConfig.CurrentContext()
			ui.ExitIfError(err)

			if riserContext.DemoGatewayIP == "" {
				ui.ExitErrorMsg("The Gateway IP has not been found. Use \"riser demo status\" to ensure that the demo is running properly.")
			}

			app, err := config.LoadAppFromConfig(appFilePath)
			ui.ExitIfErrorMsg(err, "Error loading app config")

			if deploymentName == "" {
				deploymentName = string(app.Name)
			}

			hostName := fmt.Sprintf("%s.%s.demo.riser", deploymentName, app.Namespace)

			fmt.Printf("curl -k https://%s%s --resolve %s:443:%s", hostName, path, hostName, riserContext.DemoGatewayIP)
		},
	}
	addDeploymentNameFlag(cmd.Flags(), &deploymentName)
	addAppFilePathFlag(cmd.Flags(), &appFilePath)
	return cmd
}

func newInstallDemoCommand(config *rc.RuntimeConfiguration, assets fs.FS) *cobra.Command {
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

func demoInstall(config *rc.RuntimeConfiguration, assets fs.FS) {
	out, err := exec.Command("kubectl", "version", "--client=true", "-o=json").Output()
	ui.ExitIfErrorMsg(err, "Error validating kubectl")

	versionOutput := kubectlVersion{}
	err = json.Unmarshal(out, &versionOutput)
	ui.ExitIfErrorMsg(err, "Unable to parse kubectl version")

	parsedVersion, err := version.NewVersion(versionOutput.ClientVersion.GitVersion)
	ui.ExitIfErrorMsg(err, "Unable to parse kubectl version")

	constraint, err := version.NewConstraint(kubectlVersionConstraint)
	ui.ExitIfErrorMsg(err, "Invalid kubectl version constraint")

	if !constraint.Check(parsedVersion) {
		ui.ExitErrorMsg(fmt.Sprintf("Invalid kubectl version. Must be %q", kubectlVersionConstraint))
	}

	ui.ExitIfError(err)
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
	var gitUrlParsed *url.URL
	gitUrlPrompt := &survey.Input{
		Message: "Enter the GitHub URL for the Riser state repo. Riser requires write access to this repo. If using an HTTPS url, you must include auth.\n",
		Help: ui.StripNewLines(`
The Riser state repo contains all kubernetes state for Riser apps and infrastructure.
It's recommended that you use a deploy key (e.g. git@github.com:your/repo) with write access. You may also use a Personal Access Token with repo write access (e.g. https://YOUR-TOKEN@github.com/your/repo).
`),
	}

	for {
		err = survey.AskOne(gitUrlPrompt, &gitUrl, survey.WithValidator(func(ans interface{}) error {
			return validation.Validate(ans,
				validation.Required,
				validation.By(func(v interface{}) error {
					gitUrlParsed, err = giturls.Parse(v.(string))
					return err
				}))
		}))
		if err == nil {
			break
		}
		if err.Error() == "interrupt" {
			ui.ExitErrorMsg("aborted")
		}
	}

	var gitSshKeyPath string
	if gitUrlParsed.Scheme != "https" {
		gitSshPrompt := &survey.Input{
			Message: "Enter the path to your git ssh or deploy private key.",
			Help:    ui.StripNewLines("If using a deploy key, it must have write access to the repo."),
		}
		for {
			err = survey.AskOne(gitSshPrompt, &gitSshKeyPath, survey.WithValidator(func(ans interface{}) error {
				return validation.Validate(ans,
					validation.Required,
					validation.By(func(v interface{}) error {
						_, err := os.Stat(expandTildeInPath(v.(string)))
						return err
					}))
			}))
			if err == nil {
				break
			}
			if err.Error() == "interrupt" {
				ui.ExitErrorMsg("aborted")
			}
		}
	}

	logger.Log().Info("Installing demo")

	deployment := infra.NewRiserDeployment(assets, config, gitUrl, demoEnvironmentName)

	deployment.GitSSHKeyPath = expandTildeInPath(gitSshKeyPath)
	err = deployment.Deploy()
	ui.ExitIfError(err)

	logger.Log().Info(style.Good("Installation Complete!"))
	logger.Log().Info("Executing \"riser demo status\"...")

	demoStatus(config)
}

func demoStatus(config *rc.RuntimeConfiguration) {
	logger.Log().Warn(`If you're using minikube be sure that "minikube tunnel" is running (Note: minikube tunnel may ask for password).`)
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
• You may run "riser demo status" at any time to check if the issue is resolved.
• If you're using minikube be sure that "minikube tunnel" is running (Note: minikube tunnel may ask for password).
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

	gatewayIp := strings.TrimSpace(ui.StripNewLines(ingressGatewayStep.State("stdout").(string)))

	riserContext, err := config.CurrentContext()
	ui.ExitIfError(err)
	riserContext.DemoGatewayIP = gatewayIp
	config.SetContext(riserContext)
	err = rc.SaveRc(config)
	ui.ExitIfErrorMsg(err, "Error saving Riser config")

	logger.Log().Info("\n" + style.Good("🚀 Everything checks out!") + "\n")

	logger.Log().Info("Environment:\t" + style.Emphasis(demoEnvironmentName))
	logger.Log().Info("Gateway IP:\t" + style.Emphasis(gatewayIp))
	logger.Log().Info("API Host:\t" + style.Emphasis("riser-server.riser-system.demo.riser"))

	logger.Log().Info("\nInstructions:")
	logger.Log().Info(fmt.Sprintf("• In your hosts file (e.g. /etc/hosts on OSX) or local DNS server set the IP for the host %s to the gateway IP: %s", style.Emphasis("riser-server.riser-system.demo.riser"), style.Emphasis(gatewayIp)))
	logger.Log().Info(fmt.Sprintf("  Example /etc/hosts entry:\n  %s", style.Muted(fmt.Sprintf("%s riser-server.riser-system.demo.riser", gatewayIp))))
	logger.Log().Info("• For easier access to your apps, you may wish to add additional host entries for each app using the format <YOUR-APP>.apps.demo.riser to the same gateway IP, or create a wildcard DNS record for *.apps.demo.riser.")
	logger.Log().Info(`• Try out the testdummy app!
  - In an empty folder create the app with a default config using "riser apps init testdummy"
  - Edit "app.yaml" and specify "tshak/testdummy" as the docker image
  - Deploy using "riser deploy latest demo"
  - Use "riser status" to check the status of your deployment
  - Once deployed access using "riser demo curl | sh" (or "curl -k https://testdummy.apps.demo.riser" if you have your own DNS server). If all went well you should receive a HTTP 200 response with the text "pong".`)

	logger.Log().Info("\nExecute \"riser demo status\" to see this message again.")
}
