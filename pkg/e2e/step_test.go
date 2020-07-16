package e2e

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_stepColor(t *testing.T) {
	testNames := len(allStepColors)
	for i := 0; i < testNames; i++ {
		// Run twice to ensure that it's consistent
		assert.Equal(t, allStepColors[i], stepColor(fmt.Sprintf("testName-%d", i)))
		assert.Equal(t, allStepColors[i], stepColor(fmt.Sprintf("testName-%d", i)))
	}
	// Reuse the colors
	result := stepColor("testReuseColor")
	assert.Equal(t, allStepColors[0], result)
}
