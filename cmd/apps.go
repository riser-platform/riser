package cmd

import (
	"fmt"
	"riser/rc"
	"riser/ui"
	"riser/ui/table"

	"github.com/spf13/cobra"
	"github.com/tshak/riser-server/api/v1/model"
	"github.com/tshak/riser/sdk"
)

func newAppsCommand(currentContext *rc.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apps",
		Short: "Commands for apps",
	}

	cmd.AddCommand(newAppsListCommand(currentContext))
	cmd.AddCommand(newAppsNewCommand(currentContext))

	return cmd
}

func newAppsListCommand(currentContext *rc.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists all apps",
		Run: func(cmd *cobra.Command, args []string) {
			apiClient, err := sdk.NewClient(currentContext.ServerURL, currentContext.Apikey)
			ui.ExitIfError(err)
			apps, err := apiClient.ListApps()
			ui.ExitIfError(err)

			table := table.Default().Header("Name", "Id")

			for _, app := range apps {
				table.AddRow(app.Name, app.Id)
			}

			fmt.Println(table)
		},
	}
}

func newAppsNewCommand(currentContext *rc.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "new (app name)",
		Short: "Creates a new app",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			appName := args[0]

			apiClient, err := sdk.NewClient(currentContext.ServerURL, currentContext.Apikey)
			ui.ExitIfError(err)
			app, err := apiClient.PostApp(&model.NewApp{Name: appName})
			ui.ExitIfError(err)

			fmt.Printf("App %s created. Please add the following id to your manifest: %s", app.Name, app.Id)
		},
	}
}
