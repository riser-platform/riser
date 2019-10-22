package cmd

import (
	"fmt"
	"riser/pkg/rc"
	"riser/pkg/ui"
	"riser/pkg/ui/table"

	"github.com/spf13/cobra"
)

func newStagesCommand(currentContext *rc.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stages",
		Short: "Commands for stages.",
		Long:  "Commands for stages. A stage represents a single kubernetes cluster. Stages are commonly have names like \"dev\", \"test\", or \"prod\". Stages are created automatically after installing the riser controller in a cluster.",
	}

	cmd.AddCommand(newStagesListCommand(currentContext))

	return cmd
}

func newStagesListCommand(currentContext *rc.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists all available stages",
		Run: func(cmd *cobra.Command, args []string) {
			riserClient := getRiserClient(currentContext)
			stages, err := riserClient.Stages.List()
			ui.ExitIfError(err)

			table := table.Default().Header("Name")

			for _, stage := range stages {
				table.AddRow(stage.Name)
			}

			fmt.Println(table)
		},
	}
}
