package cmd

import (
	"fmt"
	"riser/config"
	"riser/logger"

	"github.com/sanity-io/litter"
	"github.com/spf13/cobra"
)

func newValidateCommand() *cobra.Command {
	var appFilePath string
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validates an app config",
		Run: func(cmd *cobra.Command, args []string) {
			app, err := config.LoadApp(appFilePath)
			if err == nil {
				logger.Log().Info(fmt.Sprintf("Loaded config %s", appFilePath))
				logger.Log().Verbose(litter.Sdump(app))
			} else {
				logger.Log().Error(fmt.Sprintf("Failed to load config %s: %s", appFilePath, err))
			}
		},
	}

	addAppFilePathFlag(cmd.Flags(), &appFilePath)

	return cmd
}
