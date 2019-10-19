package cmd

import (
	"crypto/rand"
	"fmt"
	"net/url"
	"os/exec"
	"path"
	"riser/logger"
	"riser/ui"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

const ApiKeySizeBytes = 20

func newOpsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ops",
		Short: "Commands for operational tasks. These are not typically needed for day-to-day usage of riser.",
	}

	cmd.AddCommand(newGenerateApikeyCommand())
	cmd.AddCommand(newInstallDemoCommand())

	return cmd
}

func newGenerateApikeyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "generate-apikey",
		Short: "Generates a riser compliant API KEY. This is typically used for bootrapping.",
		Long:  "Generates a riser compliant API KEY. This is typically used for bootrapping the riser server. For user creation, see \"riser users\" for creating new users with API KEYS.",
		Run: func(cmd *cobra.Command, args []string) {
			var key = make([]byte, ApiKeySizeBytes)
			_, err := rand.Read(key)
			ui.ExitIfErrorMsg(err, "Error generating API KEY")

			fmt.Printf("%x", key)
		},
	}
}

type DemoConfig struct {
}

func newInstallDemoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install-demo",
		Short: "Installs a self-contained riser demo to a k8s cluster",
		Long:  "Install a self-contained riser demo along with all required infrastructure (istio, postgresql, etc)",
		Run: func(cmd *cobra.Command, args []string) {
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
				Message: fmt.Sprintf("Are you sure you wish to install the demo to the current context %q?", strings.TrimSpace(string(kcOutput))),
			}
			err = survey.AskOne(prompt, &useKc)
			ui.ExitIfError(err)

			if !useKc {
				ui.ExitErrorMsg("Please change to the desired kube context and try again")
			}

			var gitUrl string
			gitUrlPrompt := &survey.Input{
				Message: "Enter the GitHub URL as the riser state repo (e.g. https://github.com/your/repo). It's recommended to use an empty repo.",
			}
			err = survey.AskOne(gitUrlPrompt, &gitUrl)
			ui.ExitIfError(err)

			err = validation.Validate(&gitUrl, is.URL, validation.Required)
			// TODO: reprompt instead of exit on validation failure
			ui.ExitIfError(err)

			gitUsername := "oauthtoken"
			gitUsernamePrompt := &survey.Input{
				Message: "Enter a username that has write access to the repo. Use \"oauthtoken\" if using a Github Personal Access Token (recommended).",
				Help:    "See https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line for more information on creating a GitHub Personal Access Token",
				Default: gitUsername,
			}
			err = survey.AskOne(gitUsernamePrompt, &gitUsername)
			ui.ExitIfError(err)

			var gitPassword string
			gitPasswordPrompt := &survey.Password{
				Message: "Enter the password or Personal Access Token",
			}
			err = survey.AskOne(gitPasswordPrompt, &gitPassword)
			ui.ExitIfError(err)

			gitUrlParsed, err := url.Parse(gitUrl)
			ui.ExitIfErrorMsg(err, "Unable to parse git URL")
			gitUrlParsed.User = url.UserPassword(gitUsername, gitPassword)

			// TODO: Wrap with a timeout
			_, err = exec.Command("git", "ls-remote", gitUrlParsed.String(), "HEAD").Output()
			ui.ExitIfErrorMsg(err, fmt.Sprintf("Error validating git remote."))

			logger.Log().Info("Installing kubeapplier...")
			_, err = exec.Command("kubectl", "apply", "-f", path.Join(demoPath, "kubeapplier")).Output()
			ui.ExitIfError(err)
		},
	}

	return cmd
}
