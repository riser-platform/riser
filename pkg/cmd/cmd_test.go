package cmd

import (
	"riser/pkg/rc"
	"testing"
)

func Test_CreateCmd_DoesNotExitOnEmptyRC(t *testing.T) {
	// This guards against someone attempting to access the currentContext while building up the commands
	Execute(&Runtime{
		Configuration: &rc.RuntimeConfiguration{},
	})
}
