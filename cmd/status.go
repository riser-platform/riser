package cmd

import (
	"github.com/spf13/cobra"
)

func newStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status (deployment or app name)",
		Short: "Gets the status for a deployment.",
	}
}
