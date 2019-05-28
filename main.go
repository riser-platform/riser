package main

import (
	"riser/cmd"

	"github.com/hashicorp/go-version"
)

// versionString is a var because overwritten by the compiler using ldflags
var versionString = "0.0.0-local"

// TODO: Implement server config and context support
const baseUri = "http://localhost:8000"

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
