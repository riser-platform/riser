package cmd

import (
	"fmt"
	"riser/rc"
	"riser/ui"
	"riser/ui/table"
	"time"

	"github.com/tshak/riser/sdk"

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
		Use:   "save (name) (plaintextsecret) (stage)",
		Short: "Creates a new secret or updates an existing one",
		Long:  "Creates a new secret or updates an existing one. Secrets are maintained seperately per stage.",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			secretName := args[0]
			plainTextSecret := args[1]
			stageName := args[2]

			apiClient, err := sdk.NewClient(currentContext.ServerURL, currentContext.Apikey)
			ui.ExitIfError(err)

			// TODO: Prompt to confirm first
			err = apiClient.SaveSecret(appName, stageName, secretName, plainTextSecret)
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
			apiClient, err := sdk.NewClient(currentContext.ServerURL, currentContext.Apikey)
			ui.ExitIfError(err)

			secretMetas, err := apiClient.ListSecretMetas(appName, stageName)
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
