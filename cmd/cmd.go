// Package cmd contains commands accessible via the CLI
package cmd

import (
	"os"
	"riser/logger"

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

	cmd.AddCommand(newAppsCommand())
	cmd.AddCommand(newDeployCommand())
	cmd.AddCommand(newStatusCommand())
	cmd.AddCommand(newValidateCommand())
	cmd.AddCommand(newVersionCmd(runtime.Version))
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
