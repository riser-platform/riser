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
	cmd := &cobra.Command{
		Use:   os.Args[0],
		Short: "Riser platform",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logger.SetLogger(logger.NewScreenLogger(verbose))
		},
	}

	cmd.AddCommand(newAppsCommand(runtime.Configuration))
	cmd.AddCommand(newContextCommand(runtime.Configuration))
	cmd.AddCommand(newDemoCommand(runtime.Configuration, runtime.Assets))
	cmd.AddCommand(newDeployCommand(runtime.Configuration))
	cmd.AddCommand(newDeploymentsCommand(runtime.Configuration))
	cmd.AddCommand(newNamespacesCommand(runtime.Configuration))
	cmd.AddCommand(newOpsCommand())
	cmd.AddCommand(newRolloutCommand(runtime.Configuration))
	cmd.AddCommand(newStagesCommand(runtime.Configuration))
	cmd.AddCommand(newSecretsCommand(runtime.Configuration))
	cmd.AddCommand(newStatusCommand(runtime.Configuration))
	cmd.AddCommand(newValidateCommand(runtime.Configuration))
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
