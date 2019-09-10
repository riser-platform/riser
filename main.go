package main

import (
	"riser/cmd"
	"riser/rc"
	"riser/ui"

	"github.com/hashicorp/go-version"
)

// versionString is a var because it's overwritten by the compiler using ldflags
var versionString = "0.0.0-local"

func main() {
	currentVersion, err := version.NewVersion(versionString)
	ui.ExitIfErrorMsg(err, "Invalid version")

	config, err := rc.LoadRc()
	ui.ExitIfErrorMsg(err, "Unable to load runtime configuration")

	// Main execution path
	cmd.Execute(&cmd.Runtime{
		Version:       currentVersion,
		Configuration: config,
	})
}
