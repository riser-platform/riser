package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tshak/riser-server/api/v1/model"
	"github.com/tshak/riser/sdk"
)

func newAppsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apps",
		Short: "Commands for apps",
	}

	cmd.AddCommand(newAppsListCommand())
	cmd.AddCommand(newAppsNewCommand())

	return cmd
}

func newAppsListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists available apps",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Not implemented")
		},
	}
}

func newAppsNewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "new (app name)",
		Short: "Creates a new app",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			appName := args[0]

			apiClient, err := sdk.NewClient("http://localhost:8000")
			if err != nil {
				panic(err)
			}
			app, err := apiClient.PostApp(&model.NewApp{Name: appName})
			if err != nil {
				panic(err)
			}

			fmt.Printf("App %s created. Please add the following id to your manifest: %s\n", app.Name, app.Id)
		},
	}
}
