package cmd

import (
	"github.com/spf13/pflag"
)

func addAppFlag(flags *pflag.FlagSet, p *string) {
	flags.StringVarP(p, "app", "", "", "The name of the application. Defaults to the app name in ./app.yml or ./app.yaml if specified")
}
