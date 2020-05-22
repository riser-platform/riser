package cmd

import (
	"fmt"
	"riser/pkg/rc"
	"riser/pkg/ui"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

func newDeploymentsCommand(runtimeConfig *rc.RuntimeConfiguration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deployments",
		Short: "Commands for managing deployments",
		Long:  "Commands for managing deployments. Use \"riser deploy\" to create a new deployment or revision.",
	}

	cmd.AddCommand(newDeploymentsDeleteCommand(runtimeConfig))

	return cmd
}

func newDeploymentsDeleteCommand(runtimeConfig *rc.RuntimeConfiguration) *cobra.Command {
	var namespace string
	noPrompt := false
	cmd := &cobra.Command{
		Use:   "delete (deploymentName) (targetEnvironment)",
		Short: "Permanentally deletes a deployment and all of its revisions in the specified environment",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			currentContext := safeCurrentContext(runtimeConfig)
			deleteConfirmed := false
			deploymentName := args[0]
			environmentName := args[1]
			prompt := &survey.Confirm{
				Message: fmt.Sprintf("Are you sure you wish to delete the deployment %q in environment %q?", deploymentName, environmentName),
			}

			if !noPrompt {
				err := survey.AskOne(prompt, &deleteConfirmed)
				ui.ExitIfError(err)
				if !deleteConfirmed {
					return
				}
			}

			riserClient := getRiserClient(currentContext)
			result, err := riserClient.Deployments.Delete(deploymentName, namespace, environmentName)
			ui.ExitIfError(err)

			fmt.Println(result.Message)
		},
	}

	cmd.Flags().BoolVar(&noPrompt, "no-prompt", false, "do not prompt for a confirmation")
	addNamespaceFlag(cmd.Flags(), &namespace)

	return cmd
}
