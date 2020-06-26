package cmd

import (
	"fmt"
	"riser/pkg/rc"
	"riser/pkg/ui"

	"github.com/spf13/cobra"
)

func newSecretsCommand(runtimeConfig *rc.RuntimeConfiguration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secrets",
		Short: "Commands for secrets",
	}

	cmd.AddCommand(newSecretsListCommand(runtimeConfig))
	cmd.AddCommand(newSecretsSaveCommand(runtimeConfig))
	return cmd
}

func newSecretsSaveCommand(runtimeConfig *rc.RuntimeConfiguration) *cobra.Command {
	var appName string
	var namespace string
	cmd := &cobra.Command{
		Use:   "save (name) (plaintextsecret) (targetEnvironment)",
		Short: "Creates a new secret or updates an existing one",
		Long:  "Creates a new secret or updates an existing one. Secrets are stored seperately per app and environment.",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			currentContext := safeCurrentContext(runtimeConfig)
			secretName := args[0]
			plainTextSecret := args[1]
			environmentName := args[2]

			riserClient := getRiserClient(currentContext)

			err := riserClient.Secrets.Save(appName, namespace, environmentName, secretName, plainTextSecret)
			ui.ExitIfErrorMsg(err, "Error saving secret")

			fmt.Printf("Secret %q saved in environment %q. Changes will take affect for new deployments.\n", secretName, environmentName)
		},
	}
	addAppFlag(cmd.Flags(), &appName)
	addNamespaceFlag(cmd.Flags(), &namespace)

	return cmd
}

func newSecretsListCommand(runtimeConfig *rc.RuntimeConfiguration) *cobra.Command {
	var appName string
	var namespace string
	cmd := &cobra.Command{
		Use:   "list (environment)",
		Short: "Lists secrets configured for a given environment",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			currentContext := safeCurrentContext(runtimeConfig)
			environmentName := args[0]
			riserClient := getRiserClient(currentContext)

			secretMetas, err := riserClient.Secrets.List(appName, namespace, environmentName)
			ui.ExitIfError(err)

			view := &ui.BasicTableView{}
			view.Header("Name", "Rev")

			for _, secretMeta := range secretMetas {
				view.AddRow(
					secretMeta.Name,
					fmt.Sprintf("%d", secretMeta.Revision))
			}

			ui.RenderView(view)
		},
	}

	addAppFlag(cmd.Flags(), &appName)
	addNamespaceFlag(cmd.Flags(), &namespace)
	addOutputFlag(cmd.Flags())

	return cmd
}
