package main

import (
	"riser/cmd"

	"github.com/hashicorp/go-version"
)

// versionString is a var because overwritten by the compiler using ldflags
var versionString = "0.0.0-local"

func main() {
	currentVersion, err := version.NewVersion(versionString)
	if err != nil {
		panic(err)
	}

	// Main execution path
	cmd.Execute(&cmd.Runtime{
		Version: currentVersion,
	})
}
