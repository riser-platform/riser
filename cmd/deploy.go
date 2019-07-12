package cmd

import (
	"fmt"
	"riser/config"
	"riser/rc"
	"riser/ui"

	"github.com/sanity-io/litter"
	"github.com/tshak/riser-server/api/v1/model"
	"github.com/tshak/riser/sdk"

	"github.com/spf13/cobra"
)

func newDeployCommand(currentContext *rc.RuntimeContext) *cobra.Command {
	var appFilePath string
	var dryRun bool
	var deploymentName string
	cmd := &cobra.Command{
		Use:   "deploy (docker tag) (stage)",
		Short: "Creates a new deployment",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			dockerTag := args[0]
			stage := args[1]

			app, err := config.LoadApp(appFilePath)
			ui.ExitIfErrorMsg(err, "Error loading app config")

			if dryRun {
				println("DRY RUN MODE")
			}

			deployment := &model.RawDeployment{
				DeploymentMeta: model.DeploymentMeta{
					Name:   deploymentName,
					Stage:  stage,
					Docker: model.DeploymentDocker{Tag: dockerTag},
				},
				App: *app,
			}

			if dryRun {
				fmt.Println(litter.Sdump(deployment))
			} else {
				apiClient, err := sdk.NewClient(currentContext.ServerURL)
				ui.ExitIfError(err)

				message, err := apiClient.PutDeployment(deployment, dryRun)
				ui.ExitIfError(err)

				fmt.Println(message)
			}
		},
	}

	cmd.Flags().StringVarP(&deploymentName, "name", "n", "", "Optionally name the deployment. When specified the full deployment name will be <APP>-name (e.g. myapp-mydeployment).")
	addAppFilePathFlag(cmd.Flags(), &appFilePath)
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "Prints the deployment but does not create it")

	return cmd
}
