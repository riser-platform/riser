// Package cmd contains commands accessible via the CLI
package cmd

import (
	"fmt"
	"os"
	"riser/logger"
	"riser/ui"

	"github.com/spf13/cobra"
)

var verbose bool

// Execute creates the root command and executes it
func Execute(runtime *Runtime) {
	currentContext, err := runtime.Configuration.CurrentContext()
	ui.ExitIfErrorMsg(err, "Error loading current context")

	cmd := &cobra.Command{
		Use:   os.Args[0],
		Short: "Riser platform",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logger.SetLogger(logger.NewScreenLogger(verbose))
		},
	}

	cmd.AddCommand(newAppsCommand(currentContext))
	cmd.AddCommand(newDeployCommand(currentContext))
	cmd.AddCommand(newOpsCommand())
	cmd.AddCommand(newSecretsCommand(currentContext))
	cmd.AddCommand(newStatusCommand(currentContext))
	cmd.AddCommand(newValidateCommand())
	cmd.AddCommand(newVersionCmd(runtime.Version))
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	err = cmd.Execute()
	// Only show err in verbose mode as must errors are cobra
	// errors and already printed
	if err != nil && verbose {
		fmt.Printf("%#v\n", err)
	}
}
