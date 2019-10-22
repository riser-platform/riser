package cmd

import (
	"fmt"
	"riser/pkg/logger"
	"riser/pkg/rc"
	"riser/pkg/ui"

	"github.com/spf13/cobra"
)

func newContextCommand(config *rc.RuntimeConfiguration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Commands related to managing contexts.",
		Long:  "Commands related to managing contexts. Contexts are used to manage switching between multiple riser instances. Since Riser manages apps across multiple stages (clusters), this is typically only used for demo or development purposes.",
	}

	cmd.AddCommand(newContextSaveCommand(config))
	cmd.AddCommand(newContextRemoveCommand(config))
	cmd.AddCommand(newContextCurrentCommand(config))
	cmd.AddCommand(newContextListCommand(config))
	return cmd
}

func newContextSaveCommand(config *rc.RuntimeConfiguration) *cobra.Command {
	secure := true
	cmd := &cobra.Command{
		Use:   "save <contextName> <serverUrl> <apikey>",
		Short: "Adds or updates a context",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			contextName := args[0]
			ctx := &rc.Context{Name: contextName, ServerURL: args[1], Apikey: args[2], Secure: &secure}
			config.SaveContext(ctx)
			err := rc.SaveRc(config)
			ui.ExitIfErrorMsg(err, "Error saving rc file")

			logger.Log().Info(fmt.Sprintf("Context %q saved. Current context is now set to %q.", contextName, contextName))
		},
	}

	cmd.Flags().BoolVar(&secure, "secure", true, "Set to false to skip TLS verification")

	return cmd
}

func newContextRemoveCommand(config *rc.RuntimeConfiguration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <contextName>",
		Args:  cobra.ExactArgs(1),
		Short: "Removes a context",
		Run: func(cmd *cobra.Command, args []string) {
			err := config.RemoveContext(args[0])
			ui.ExitIfErrorMsg(err, "Error removing context")
			err = rc.SaveRc(config)
			ui.ExitIfErrorMsg(err, "Error saving to rc file")
		},
	}

	return cmd
}

func newContextCurrentCommand(config *rc.RuntimeConfiguration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current [setCurrentContextName]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Gets or sets the current context",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 && len(args[0]) > 0 {
				err := config.SetCurrentContext(args[0])
				ui.ExitIfErrorMsg(err, "unable to set context")
				err = rc.SaveRc(config)
				ui.ExitIfErrorMsg(err, "Error saving to rc file")

				logger.Log().Info(fmt.Sprintf("Successfully loaded context \"%s\"\n", config.CurrentContextName))
			} else {
				logger.Log().Info(fmt.Sprintf("Current Context: \"%s\"\n", config.CurrentContextName))
			}
		},
	}

	return cmd
}

func newContextListCommand(config *rc.RuntimeConfiguration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists contexts",
		Run: func(cmd *cobra.Command, args []string) {
			contexts := config.GetContexts()
			if len(contexts) == 0 {
				logger.Log().Info("No contexts configured. Use \"riser context add\" to add a new context")
			} else {
				for _, context := range contexts {
					logger.Log().Info(context.Name)
				}
			}
		},
	}

	return cmd
}
