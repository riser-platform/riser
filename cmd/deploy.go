package cmd

import (
	"fmt"
	"riser/config"
	"riser/rc"
	"riser/ui"

	"github.com/wzshiming/ctc"

	"github.com/tshak/riser-server/api/v1/model"
	"github.com/tshak/riser/sdk"

	"github.com/spf13/cobra"
)

func newDeployCommand(currentContext *rc.Context) *cobra.Command {
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

			deployment := &model.DeploymentRequest{
				DeploymentMeta: model.DeploymentMeta{
					Name:   deploymentName,
					Stage:  stage,
					Docker: model.DeploymentDocker{Tag: dockerTag},
				},
				App: *app,
			}

			apiClient, err := sdk.NewClient(currentContext.ServerURL, currentContext.Apikey)
			ui.ExitIfError(err)

			deployResult, err := apiClient.Deployments.Save(deployment, dryRun)
			ui.ExitIfError(err)

			fmt.Println(deployResult.Message)
			if dryRun && deployResult.DryRunCommits != nil {
				for _, commit := range deployResult.DryRunCommits {
					fmt.Print(ctc.ForegroundBrightCyan)
					fmt.Printf("Commit: %s\n", commit.Message)
					for _, file := range commit.Files {
						fmt.Print(ctc.ForegroundBrightWhite)
						fmt.Printf("File: %s\n", file.Name)
						fmt.Print(ctc.ForegroundBrightBlack)
						fmt.Println(file.Contents)
					}
					fmt.Print(ctc.Reset)
				}
			}
		},
	}

	cmd.Flags().StringVarP(&deploymentName, "name", "n", "", "Optionally name the deployment. When specified the full deployment name will be <APP>-name (e.g. myapp-mydeployment).")
	addAppFilePathFlag(cmd.Flags(), &appFilePath)
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "Prints the deployment but does not create it")

	return cmd
}
