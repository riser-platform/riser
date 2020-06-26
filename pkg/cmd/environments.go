package cmd

import (
	"riser/pkg/rc"
	"riser/pkg/ui"

	"github.com/spf13/cobra"
)

func newEnvironmentsCommand(runtimeConfig *rc.RuntimeConfiguration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "environments",
		Short: "Commands for environments.",
		Long:  "Commands for environments. A environment represents a single kubernetes cluster. Environments may have names like \"dev\", \"test\", or \"prod\". Environments are created automatically after installing the riser controller in a cluster.",
	}

	cmd.AddCommand(newEnvironmentsListCommand(runtimeConfig))

	return cmd
}

func newEnvironmentsListCommand(runtimeConfig *rc.RuntimeConfiguration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all available environments",
		Run: func(cmd *cobra.Command, args []string) {
			currentContext := safeCurrentContext(runtimeConfig)
			riserClient := getRiserClient(currentContext)
			environments, err := riserClient.Environments.List()
			ui.ExitIfError(err)

			view := &ui.BasicTableView{}
			view.Header("Name")

			for _, environment := range environments {
				view.AddRow(environment.Name)
			}

			ui.RenderView(view)
		},
	}

	addOutputFlag(cmd.Flags())

	return cmd
}
