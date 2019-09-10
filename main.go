package main

import (
	"riser/cmd"
	"riser/rc"

	"github.com/hashicorp/go-version"
)

// versionString is a var because it's overwritten by the compiler using ldflags
var versionString = "0.0.0-local"

func main() {
	currentVersion, err := version.NewVersion(versionString)
	if err != nil {
		panic(err)
	}

	// TODO: Load rc from file or prompt user to create new
	config := &rc.RuntimeConfiguration{}

	// Main execution path
	cmd.Execute(&cmd.Runtime{
		Version:       currentVersion,
		Configuration: config,
	})
}
