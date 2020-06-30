package cmd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validateNewDeployCommand(t *testing.T) {
	tests := []struct {
		wait          bool
		manualRollout bool
		expected      error
	}{
		{true, false, nil},
		{false, true, nil},
		{true, true, errors.New(`You cannot specify both "--wait" and "--manual-rollout"`)},
	}

	for _, tt := range tests {
		err := validateNewDeployCommand(tt.manualRollout, tt.wait)
		assert.Equal(t, tt.expected, err)
	}
}
