package cmd

import (
	"fmt"
	"riser/rc"
	"riser/ui"
	"riser/ui/table"

	"github.com/spf13/cobra"
	"github.com/tshak/riser/sdk"
)

func newStagesCommand(currentContext *rc.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stages",
		Short: "Commands for stages",
	}

	cmd.AddCommand(newStagesListCommand(currentContext))

	return cmd
}

func newStagesListCommand(currentContext *rc.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists all available stages",
		Run: func(cmd *cobra.Command, args []string) {
			apiClient, err := sdk.NewClient(currentContext.ServerURL, currentContext.Apikey)
			ui.ExitIfError(err)
			stages, err := apiClient.ListStages()
			ui.ExitIfError(err)

			table := table.Default().Header("Name")

			for _, stage := range stages {
				table.AddRow(stage.Name)
			}

			fmt.Println(table)
		},
	}
}
