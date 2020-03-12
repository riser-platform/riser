package cmd

import (
	"fmt"
	"riser/pkg/logger"
	"riser/pkg/rc"
	"riser/pkg/ui"
	"riser/pkg/ui/table"

	"github.com/spf13/cobra"
)

func newNamespacesCommand(runtimeConfig *rc.RuntimeConfiguration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "namespaces",
		Short: "Commands for managing namespaces",
	}

	cmd.AddCommand(newNamespacesCreateCommand(runtimeConfig))
	cmd.AddCommand(newNamespacesListCommand(runtimeConfig))
	return cmd
}

func newNamespacesCreateCommand(runtimeConfig *rc.RuntimeConfiguration) *cobra.Command {
	return &cobra.Command{
		Use:   "create (namespace name)",
		Short: "Create a new namespace",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			currentContext := safeCurrentContext(runtimeConfig)
			namespaceName := args[0]
			riserClient := getRiserClient(currentContext)
			err := riserClient.Namespaces.Create(namespaceName)
			ui.ExitIfError(err)
			logger.Log().Info("Namespace created")
		},
	}
}

func newNamespacesListCommand(runtimeConfig *rc.RuntimeConfiguration) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all namespaces",
		Run: func(*cobra.Command, []string) {
			currentContext := safeCurrentContext(runtimeConfig)
			riserClient := getRiserClient(currentContext)
			namespaces, err := riserClient.Namespaces.List()
			ui.ExitIfErrorMsg(err, "error listing namespaces")

			table := table.Default().Header("Name")

			for _, ns := range namespaces {
				table.AddRow(string(ns.Name))
			}

			fmt.Println(table)
		},
	}
}
