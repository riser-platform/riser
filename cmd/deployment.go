package cmd

import (
	"fmt"
	"riser/client"
	"riser/config"

	"github.com/sanity-io/litter"

	"github.com/spf13/cobra"
)

func newDeploymentCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deployment",
		Short: "Commands for deployments",
	}

	cmd.AddCommand(newDeploymentStatusCommand())
	cmd.AddCommand(newDeploymentNewCommand())

	return cmd
}

func newDeploymentStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status (deployment or app name)",
		Short: "Gets the status for a deployment.",
	}
}

func newDeploymentNewCommand() *cobra.Command {
	var appFilePath string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "new (stage) (dockerTag) [deployment name (appName)]",
		Short: "Creates a new deployment",
		Args:  cobra.RangeArgs(2, 3),
		Run: func(cmd *cobra.Command, args []string) {
			stage := args[0]
			dockerTag := args[1]

			app, err := config.LoadApp(appFilePath)
			if err != nil {
				panic(err)
			}

			if dryRun {
				println("DRY RUN MODE")
			}

			deployment := client.Deployment{
				// TODO: Support optional deploymentName
				Name:   app.Name,
				Stage:  stage,
				App:    *app,
				Docker: client.DeploymentDocker{Tag: dockerTag},
			}

			fmt.Println(litter.Sdump(deployment))

			apiClient, err := client.NewClient("http://localhost:8000")
			if err != nil {
				panic(err)
			}
			err = apiClient.PutDeployment(deployment, dryRun)
			if err != nil {
				panic(err)
			}
		},
	}

	cmd.Flags().StringVarP(&appFilePath, "file", "f", "./app.yml", "Path to the application config file")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "Prints the deployment but does not create it.")

	return cmd
}
