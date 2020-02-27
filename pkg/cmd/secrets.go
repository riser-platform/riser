package cmd

import (
	"fmt"
	"riser/pkg/rc"
	"riser/pkg/ui"
	"riser/pkg/ui/table"

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
	var appIdOrName string
	cmd := &cobra.Command{
		Use:   "save (name) (plaintextsecret) (stage)",
		Short: "Creates a new secret or updates an existing one",
		Long:  "Creates a new secret or updates an existing one. Secrets are stored seperately per app and stage.",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			secretName := args[0]
			plainTextSecret := args[1]
			stageName := args[2]

			riserClient := getRiserClient(currentContext)

			app, err := riserClient.Apps.Get(appIdOrName)
			ui.ExitIfError(err)

			err = riserClient.Secrets.Save(app.Id, stageName, secretName, plainTextSecret)
			ui.ExitIfErrorMsg(err, "Error saving secret")

			fmt.Printf("Secret %q saved in stage %q. Changes will take affect for new deployments.\n", secretName, stageName)
		},
	}
	addAppFlag(cmd.Flags(), &appIdOrName)

	return cmd
}

func newSecretsListCommand(currentContext *rc.Context) *cobra.Command {
	var appIdOrName string
	cmd := &cobra.Command{
		Use:   "list (stage)",
		Short: "Lists secrets configured for a given stage",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			stageName := args[0]
			riserClient := getRiserClient(currentContext)

			app, err := riserClient.Apps.Get(appIdOrName)
			ui.ExitIfError(err)

			secretMetas, err := riserClient.Secrets.List(app.Id, stageName)
			ui.ExitIfError(err)

			table := table.Default().Header("Name", "Rev")
			for _, secretMeta := range secretMetas {
				table.AddRow(
					secretMeta.Name,
					fmt.Sprintf("%d", secretMeta.Revision))
			}

			fmt.Println(table)
		},
	}

	addAppFlag(cmd.Flags(), &appIdOrName)

	return cmd
}
