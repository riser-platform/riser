package cmd

import (
	"errors"
	"fmt"
	"io"
	"riser/pkg/config"
	"riser/pkg/deploy"
	"riser/pkg/rc"
	"riser/pkg/ui"
	"riser/pkg/ui/style"
	"time"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/wzshiming/ctc"

	"github.com/spf13/cobra"
)

func newDeployCommand(runtimeConfig *rc.RuntimeConfiguration) *cobra.Command {
	var appFilePath string
	var dryRun bool
	var deploymentName string
	var manualRollout bool
	var wait bool
	var waitSeconds int
	cmd := &cobra.Command{
		Use:   "deploy (docker tag) (targetEnvironment)",
		Short: "Creates a new deployment or revision",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			currentContext := safeCurrentContext(runtimeConfig)
			dockerTag := args[0]
			environment := args[1]

			ui.ExitIfError(validateNewDeployCommand(manualRollout, wait))

			app, err := config.LoadAppFromConfig(appFilePath)
			ui.ExitIfErrorMsg(err, "Error loading app config")

			deployment := &model.DeploymentRequest{
				DeploymentMeta: model.DeploymentMeta{
					Name:          deploymentName,
					Environment:   environment,
					Docker:        model.DeploymentDocker{Tag: dockerTag},
					ManualRollout: manualRollout,
				},
				App: app,
			}

			riserClient := getRiserClient(currentContext)

			deployResult, err := riserClient.Deployments.Save(deployment, dryRun)
			ui.ExitIfError(err)

			if wait {
				err = deploy.WaitForReady(
					riserClient.Apps,
					model.App{Id: app.Id, Name: app.Name, Namespace: app.Namespace},
					deploymentName,
					environment,
					deployResult.RiserRevision,
					time.Duration(waitSeconds)*time.Second)
				ui.ExitIfError(err)
			} else {
				view := &newDeployView{
					result:        deployResult,
					manualRollout: manualRollout,
					dryRun:        dryRun,
				}
				ui.RenderView(view)
			}
		},
	}

	addDeploymentNameFlag(cmd.Flags(), &deploymentName)
	addAppFilePathFlag(cmd.Flags(), &appFilePath)
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "Prints the deployment but does not create it")
	cmd.Flags().BoolVarP(&manualRollout, "manual-rollout", "m", false, "When set no traffic routes to the new deployment. Use \"riser rollout\" to manually route traffic")
	cmd.Flags().BoolVar(&wait, "wait", false, "Blocks until the new deployment is ready to receive traffic or until --wait-seconds is reached. Cannot be used with --manual-rollout")
	cmd.Flags().IntVar(&waitSeconds, "wait-seconds", 60, "Sets the number of seconds for --wait")
	addOutputFlag(cmd.Flags())

	return cmd
}

func validateNewDeployCommand(manualRollout, wait bool) error {
	if manualRollout && wait {
		return errors.New(`You cannot specify both "--wait" and "--manual-rollout"`)
	}
	return nil
}

type newDeployView struct {
	result        *model.DeploymentResponse
	manualRollout bool
	dryRun        bool
}

func (view *newDeployView) RenderHuman(writer io.Writer) error {
	outStr := fmt.Sprintf("%s\n", view.result.Message)

	if view.manualRollout {
		outStr += style.Emphasis("Manual rollout specified. You must use \"riser rollout\" to route traffic to the new deployment\n")
	}

	if view.dryRun && view.result.DryRunCommits != nil {
		for _, commit := range view.result.DryRunCommits {
			outStr += ctc.ForegroundBrightCyan.String()
			outStr += fmt.Sprintf("Commit: %s\n", commit.Message)
			for _, file := range commit.Files {
				outStr += ctc.ForegroundBrightWhite.String()
				outStr += fmt.Sprintf("File: %s\n", file.Name)
				outStr += ctc.ForegroundBrightBlack.String()
				outStr += fmt.Sprintln(file.Contents)
			}
			outStr += ctc.Reset.String()
		}
	}

	_, err := writer.Write([]byte(outStr))
	return err
}

func (view *newDeployView) RenderJson(writer io.Writer) error {
	return ui.RenderJson(view.result, writer)
}
