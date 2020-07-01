package cmd

import (
	"fmt"
	"riser/pkg/config"
	"riser/pkg/logger"
	"riser/pkg/rc"
	"riser/pkg/ui"

	"github.com/spf13/cobra"
)

func newValidateCommand(runtimeConfig *rc.RuntimeConfiguration) *cobra.Command {
	var appFilePath string
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validates an app config",
		Run: func(cmd *cobra.Command, args []string) {
			currentContext := safeCurrentContext(runtimeConfig)
			app, err := config.LoadAppFromConfig(appFilePath)
			if err == nil {
				riserClient := getRiserClient(currentContext)

				err := riserClient.Validate.AppConfig(app)
				ui.ExitIfError(err)

				fmt.Println("App config is valid")
			} else {
				logger.Log().Error(fmt.Sprintf("Failed to load app config %s: %s", appFilePath, err))
			}
		},
	}

	addAppFilePathFlag(cmd.Flags(), &appFilePath)

	return cmd
}
