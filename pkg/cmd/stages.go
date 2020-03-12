package cmd

import (
	"fmt"
	"riser/pkg/rc"
	"riser/pkg/ui"
	"riser/pkg/ui/table"

	"github.com/spf13/cobra"
)

func newStagesCommand(runtimeConfig *rc.RuntimeConfiguration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stages",
		Short: "Commands for stages.",
		Long:  "Commands for stages. A stage represents a single kubernetes cluster. Stages are commonly have names like \"dev\", \"test\", or \"prod\". Stages are created automatically after installing the riser controller in a cluster.",
	}

	cmd.AddCommand(newStagesListCommand(runtimeConfig))

	return cmd
}

func newStagesListCommand(runtimeConfig *rc.RuntimeConfiguration) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists all available stages",
		Run: func(cmd *cobra.Command, args []string) {
			currentContext := safeCurrentContext(runtimeConfig)
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
