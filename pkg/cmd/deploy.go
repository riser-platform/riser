package cmd

import (
	"fmt"
	"riser/pkg/config"
	"riser/pkg/rc"
	"riser/pkg/ui"
	"riser/pkg/ui/style"

	"github.com/wzshiming/ctc"

	"github.com/riser-platform/riser-server/api/v1/model"

	"github.com/spf13/cobra"
)

func newDeployCommand(currentContext *rc.Context) *cobra.Command {
	var appFilePath string
	var dryRun bool
	var deploymentName string
	var manualRollout bool
	cmd := &cobra.Command{
		Use:   "deploy (docker tag) (stage)",
		Short: "Creates or updates a deployment",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			dockerTag := args[0]
			stage := args[1]

			app, err := config.LoadApp(appFilePath)
			ui.ExitIfErrorMsg(err, "Error loading app config")

			deployment := &model.DeploymentRequest{
				DeploymentMeta: model.DeploymentMeta{
					Name:          deploymentName,
					Stage:         stage,
					Docker:        model.DeploymentDocker{Tag: dockerTag},
					ManualRollout: manualRollout,
				},
				App: app,
			}

			riserClient := getRiserClient(currentContext)

			deployResult, err := riserClient.Deployments.Save(deployment, dryRun)
			ui.ExitIfError(err)

			fmt.Println(deployResult.Message)

			if manualRollout {
				fmt.Println(style.Emphasis("Manual rollout specified. You must use \"riser rollout\" to route traffic to the new deployment"))
			}

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

	cmd.Flags().StringVarP(&deploymentName, "name", "n", config.SafeLoadDefaultAppName(), "Optionally name the deployment. The name must follow the format <APP_NAME>-<SUFFIX> (e.g. myapp-mydeployment).")
	addAppFilePathFlag(cmd.Flags(), &appFilePath)
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "Prints the deployment but does not create it")
	cmd.Flags().BoolVarP(&manualRollout, "manual-rollout", "m", false, "When set no traffic routes to the new deployment. Use \"riser rollout\" to manually route traffic.")

	return cmd
}
