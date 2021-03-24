package main

import (
	"riser/assets"
	"riser/pkg/cmd"
	"riser/pkg/rc"
	"riser/pkg/ui"

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
		Assets:        assets.Assets,
		Version:       currentVersion,
		Configuration: config,
	})
}
