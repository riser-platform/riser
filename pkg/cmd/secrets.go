package cmd

import (
	"fmt"
	"riser/pkg/rc"
	"riser/pkg/ui"
	"riser/pkg/ui/table"
	"time"

	"github.com/spf13/cobra"
)

func newSecretsCommand(currentContext *rc.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secrets",
		Short: "Commands for secrets",
	}

	cmd.AddCommand(newSecretsListCommand(currentContext))
	cmd.AddCommand(newSecretsSaveCommand(currentContext))
	return cmd
}

func newSecretsSaveCommand(currentContext *rc.Context) *cobra.Command {
	var appName string
	cmd := &cobra.Command{
		Use:   "save (stage) (name) (plaintextsecret)",
		Short: "Creates a new secret or updates an existing one",
		Long:  "Creates a new secret or updates an existing one. Secrets are stored seperately per app and stage.",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			stageName := args[0]
			secretName := args[1]
			plainTextSecret := args[2]

			riserClient := getRiserClient(currentContext)

			// TODO: Prompt to confirm first
			err := riserClient.Secrets.Save(appName, stageName, secretName, plainTextSecret)
			ui.ExitIfErrorMsg(err, "Error saving secret")

			fmt.Printf("Secret %q saved. New values will take affect the next time %q in stage %q is deployed\n", secretName, appName, stageName)
		},
	}
	addAppFlag(cmd.Flags(), &appName)

	return cmd
}

func newSecretsListCommand(currentContext *rc.Context) *cobra.Command {
	var appName string
	cmd := &cobra.Command{
		Use:   "list (stage)",
		Short: "Lists secrets configured for a given stage",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			stageName := args[0]
			riserClient := getRiserClient(currentContext)

			secretMetas, err := riserClient.Secrets.List(appName, stageName)
			ui.ExitIfError(err)

			table := table.Default().Header("Name", "Last Updated")
			for _, secretMeta := range secretMetas {
				table.AddRow(
					secretMeta.Name,
					secretMeta.LastUpdated.In(time.Now().Location()).Format(time.RFC1123))
			}

			fmt.Println(table)
		},
	}

	addAppFlag(cmd.Flags(), &appName)

	return cmd
}
