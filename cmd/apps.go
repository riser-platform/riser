package cmd

import (
	"fmt"
	"os"
	"riser/rc"
	"riser/ui"
	"riser/ui/table"

	"github.com/spf13/cobra"
	"github.com/tshak/riser-server/api/v1/model"
	"github.com/tshak/riser/sdk"
)

const AppConfigPath = "./app.yaml"

func newAppsCommand(currentContext *rc.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apps",
		Short: "Commands for apps",
	}

	cmd.AddCommand(newAppsListCommand(currentContext))
	cmd.AddCommand(newAppsNewCommand(currentContext))
	cmd.AddCommand(newAppsInitCommand(currentContext))

	return cmd
}

func newAppsListCommand(currentContext *rc.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists all apps",
		Run: func(cmd *cobra.Command, args []string) {
			apiClient, err := sdk.NewClient(currentContext.ServerURL, currentContext.Apikey)
			ui.ExitIfError(err)
			apps, err := apiClient.Apps.List()
			ui.ExitIfError(err)

			table := table.Default().Header("Name", "Id")

			for _, app := range apps {
				table.AddRow(app.Name, app.Id)
			}

			fmt.Println(table)
		},
	}
}

func newAppsInitCommand(currentContext *rc.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init (app name)",
		Short: "Creates a new app with a default app.yaml file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			appName := args[0]
			app := createNewApp(currentContext, appName)

			file, err := os.OpenFile(AppConfigPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
			ui.ExitIfErrorMsg(err, "Error creating default app config")
			defer file.Close()

			err = sdk.DefaultAppConfig(file, appName, app.Id)
			ui.ExitIfErrorMsg(err, "Error creating default app config")
			fmt.Printf("App %s created with a default app config file %q. Please review the TODOs before deploying your app.", appName, AppConfigPath)
		},
	}

	return cmd
}

func newAppsNewCommand(currentContext *rc.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "new (app name)",
		Short: "Creates a new app",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			appName := args[0]
			app := createNewApp(currentContext, appName)

			fmt.Printf("App %s created. Please add the following id to your manifest: %s", app.Name, app.Id)
		},
	}
}

func createNewApp(currentContext *rc.Context, appName string) *model.App {
	apiClient, err := sdk.NewClient(currentContext.ServerURL, currentContext.Apikey)
	ui.ExitIfError(err)
	app, err := apiClient.Apps.Create(&model.NewApp{Name: appName})
	ui.ExitIfError(err)
	return app
}
