package cmd

import (
	"crypto/rand"
	"fmt"
	"riser/pkg/ui"

	"github.com/spf13/cobra"
)

func newOpsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ops",
		Short: "Commands for operational tasks. These are not typically needed for day-to-day usage of riser.",
	}

	cmd.AddCommand(newGenerateApikeyCommand())

	return cmd
}

func newGenerateApikeyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "generate-apikey",
		Short: "Generates a riser compliant API KEY. This is typically used for bootrapping.",
		Long:  "Generates a riser compliant API KEY. This is typically used for bootrapping the riser server. For user creation, see \"riser users\" for creating new users with API KEYS.",
		Run: func(cmd *cobra.Command, args []string) {
			var key = make([]byte, ApiKeySizeBytes)
			_, err := rand.Read(key)
			ui.ExitIfErrorMsg(err, "Error generating API KEY")

			fmt.Printf("%x", key)
		},
	}
}
