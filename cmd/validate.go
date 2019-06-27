package cmd

import (
	"fmt"
	"riser/config"
	"riser/logger"

	"github.com/sanity-io/litter"
	"github.com/spf13/cobra"
)

func newValidateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "validate (path/to/app.yml)",
		Short: "Validates an app config",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			appConfigPath := args[0]
			app, err := config.LoadApp(appConfigPath)
			if err == nil {
				logger.Log().Info(fmt.Sprintf("Loaded config %s", appConfigPath))
				logger.Log().Verbose(litter.Sdump(app))
			} else {
				logger.Log().Error(fmt.Sprintf("Failed to load config %s", appConfigPath), err)
			}

		},
	}
}
