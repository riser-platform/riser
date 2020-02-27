package cmd

import (
	"fmt"
	"os"
	"riser/pkg/logger"
	"riser/pkg/rc"
	"riser/pkg/ui"
	"riser/pkg/ui/style"
	"riser/pkg/ui/table"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser/sdk"
	"github.com/spf13/cobra"
)

const AppConfigPath = "./app.yaml"

func newAppsCommand(currentContext *rc.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apps",
		Short: "Commands for managing apps",
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
			riserClient := getRiserClient(currentContext)
			apps, err := riserClient.Apps.List()
			ui.ExitIfError(err)

			table := table.Default().Header("Name", "Id")

			for _, app := range apps {
				table.AddRow(app.Name, app.Id.String())
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
			logger.Log().Info(fmt.Sprintf("App %s created with a default app config file %q. Please review the TODO's before deploying your app.", style.Emphasis(appName), AppConfigPath))
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
	riserClient := getRiserClient(currentContext)
	app, err := riserClient.Apps.Create(&model.NewApp{Name: appName})
	ui.ExitIfError(err)
	return app
}
