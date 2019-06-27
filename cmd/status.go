package cmd

import (
	"github.com/sanity-io/litter"
	"github.com/spf13/cobra"
	"github.com/tshak/riser-server/sdk"
)

func newStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status (app name)",
		Short: "Gets the status for a deployment.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			appName := args[0]
			apiClient, err := sdk.NewClient("http://localhost:8000")
			if err != nil {
				panic(err)
			}

			statuses, err := apiClient.GetStatus(appName)

			if err != nil {
				panic(err)
			}

			litter.Dump(statuses)
		},
	}
}
