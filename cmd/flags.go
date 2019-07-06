package cmd

import (
	"riser/config"

	"github.com/spf13/cobra"

	"github.com/spf13/pflag"
)

func addAppFlag(flags *pflag.FlagSet, p *string) {
	defaultAppName := config.SafeLoadDefaultAppName()
	flags.StringVarP(p, "app", "", defaultAppName, "The name of the application. Required if no app config is present in the current directory.")
	if len(defaultAppName) == 0 {
		_ = cobra.MarkFlagRequired(flags, "app")
	}
}
