// Package cmd contains commands accessible via the CLI
package cmd

import (
	"fmt"
	"os"
	"riser/pkg/logger"

	"github.com/spf13/cobra"
)

var verbose bool

// Execute creates the root command and executes it
func Execute(runtime *Runtime) {
	currentContext, _ := runtime.Configuration.CurrentContext()
	// TODO: Lazy load current context os that we don't have an error when using context commands which don't require the current context to be set
	// ui.ExitIfErrorMsg(err, "Error loading current context")

	cmd := &cobra.Command{
		Use:   os.Args[0],
		Short: "Riser platform",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logger.SetLogger(logger.NewScreenLogger(verbose))
		},
	}

	cmd.AddCommand(newAppsCommand(currentContext))
	cmd.AddCommand(newContextCommand(runtime.Configuration))
	cmd.AddCommand(newDemoCommand(runtime.Configuration, runtime.Assets))
	cmd.AddCommand(newDeployCommand(currentContext))
	cmd.AddCommand(newDeploymentsCommand(currentContext))
	cmd.AddCommand(newNamespacesCommand(currentContext))
	cmd.AddCommand(newOpsCommand())
	cmd.AddCommand(newRolloutCommand(currentContext))
	cmd.AddCommand(newStagesCommand(currentContext))
	cmd.AddCommand(newSecretsCommand(currentContext))
	cmd.AddCommand(newStatusCommand(currentContext))
	cmd.AddCommand(newValidateCommand(currentContext))
	cmd.AddCommand(newVersionCmd(runtime.Version))
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	err := cmd.Execute()
	if err != nil {
		if verbose {
			fmt.Printf("%#v\n", err)
		}
		os.Exit(1)
	}
}
