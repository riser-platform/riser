package cmd

import (
	version "github.com/hashicorp/go-version"
)

// Runtime provides runtime information to commands
type Runtime struct {
	Version *version.Version
}
