package cmd

import (
	"github.com/spf13/cobra"
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
	}
}

func newAppsNewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "new",
		Short: "Creates a new app",
	}
}
