package cmd

import (
	"riser/config"

	"github.com/spf13/cobra"

	"github.com/spf13/pflag"
)

// addAppFlag adds --app flag for the app name
func addAppFlag(flags *pflag.FlagSet, appName *string) {
	defaultAppName := config.SafeLoadDefaultAppName()
	flags.StringVarP(appName, "app", "a", defaultAppName, "The name of the application. Required if no app config is present in the current directory.")
	if len(defaultAppName) == 0 {
		_ = cobra.MarkFlagRequired(flags, "app")
	}
}

// addAppFilePathFlag adds --file flag for the app file name
func addAppFilePathFlag(flags *pflag.FlagSet, appFilePath *string) {
	defaultPath := config.GetAppConfigPathFromDefaults()
	flags.StringVarP(appFilePath, "file", "f", defaultPath, "Path to the application config file")
	if len(defaultPath) == 0 {
		_ = cobra.MarkFlagRequired(flags, "file")
	}
}
