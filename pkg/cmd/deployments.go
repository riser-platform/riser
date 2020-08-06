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
	cmd.AddCommand(newDeploymentsDescribeCommand(runtimeConfig))

	return cmd
}

func newDeploymentsDescribeCommand(runtimeConfig *rc.RuntimeConfiguration) *cobra.Command {
	var appName string
	var namespace string
	cmd := &cobra.Command{
		Use:   "describe (deploymentName) (targetEnvironment)",
		Short: "Display details of a deployment in a specific environment",
		Long:  "Display details of a deployment in a specific environment. Use \"riser status\" to display a summary of all deployments in all environments for an app",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			deploymentName := args[0]
			environmentName := args[1]
			currentContext := safeCurrentContext(runtimeConfig)
			riserClient := getRiserClient(currentContext)

			app, err := riserClient.Apps.Get(appName, namespace)
			ui.ExitIfErrorMsg(err, "Error getting App")
			status, err := riserClient.Apps.GetStatus(appName, namespace)
			ui.ExitIfErrorMsg(err, "Error getting App status")
			envConfig, err := riserClient.Environments.GetConfig(environmentName)
			ui.ExitIfErrorMsg(err, "Error getting Environment config")

			view, err := newDeploymentsDescribeView(app, status, deploymentName, environmentName, envConfig.PublicGatewayHost)
			ui.ExitIfError(err)

			ui.RenderView(view)
		},
	}
	addAppFlag(cmd.Flags(), &appName)
	addNamespaceFlag(cmd.Flags(), &namespace)
	addOutputFlag(cmd.Flags())
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
				Message: fmt.Sprintf("Are you sure you wish to delete the deployment %q in namespace %q in environment %q?", deploymentName, namespace, environmentName),
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
