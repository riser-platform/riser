package cmd

import (
	"io/fs"
	"riser/pkg/rc"

	version "github.com/hashicorp/go-version"
)

// Runtime provides runtime information to commands
type Runtime struct {
	Version       *version.Version
	Configuration *rc.RuntimeConfiguration
	Assets        fs.FS
}
