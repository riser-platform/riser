package cmd

import (
	"fmt"
	"riser/config"

	"github.com/sanity-io/litter"
	"github.com/tshak/riser-server/sdk"
	"github.com/tshak/riser-server/api/v1/model"

	"github.com/spf13/cobra"
)

func newDeployCommand() *cobra.Command {
	var appFilePath string
	var dryRun bool
	var deploymentName string
	cmd := &cobra.Command{
		Use:   "deploy (dockerTag) (stage)",
		Short: "Creates a new deployment",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			dockerTag := args[0]
			stage := args[1]

			app, err := config.LoadApp(appFilePath)
			if err != nil {
				panic(err)
			}

			if dryRun {
				println("DRY RUN MODE")
			}

			deployment := model.RawDeployment{
				DeploymentMeta: model.DeploymentMeta {
				// TODO: Support optional deploymentName
				Name:   app.Name,
				Stage:  stage,
				Docker: model.DeploymentDocker{Tag: dockerTag},
				},
				App:    *app,
			}

			fmt.Println(litter.Sdump(deployment))

			apiClient, err := sdk.NewClient("http://localhost:8000")
			if err != nil {
				panic(err)
			}
			err = apiClient.PutDeployment(deployment, dryRun)
			if err != nil {
				panic(err)
			}
		},
	}

	cmd.Flags().StringVarP(&deploymentName, "name", "n", "", "Optionally name the deployment. When specified the full deployment name will be <APP>-name (e.g. myapp-mydeployment).")
	cmd.Flags().StringVarP(&appFilePath, "file", "f", "./app.yml", "Path to the application config file")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "Prints the deployment but does not create it")

	return cmd
}
