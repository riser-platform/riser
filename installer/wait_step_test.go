package installer

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_WaitStep_NoError_DoesNotRetry(t *testing.T) {
	step1 := newFakeStep("step1")
	waitStep := NewWaitStep(step1, 3, func(error) bool {
		assert.Fail(t, "RetryFunc should not be called")
		return false
	})

	err := Run(waitStep)

	assert.NoError(t, err)
	assert.Equal(t, 1, waitStep.attempts)
}

func Test_WaitStep_RetryFuncFalse_DoesNotRetry(t *testing.T) {
	step1 := newFakeStep("step1")
	step1.err = errors.New("step1 error")
	waitStep := NewWaitStep(step1, 3, func(error) bool { return false })

	err := Run(waitStep)

	assert.Contains(t, err.Error(), "step1 error")
	assert.Equal(t, 1, waitStep.attempts)
}

func Test_WaitStep_MaxRetries(t *testing.T) {
	step1 := newFakeStep("step1")
	step1.err = errors.New("step1 error")
	waitStep := NewWaitStep(step1, 3, func(error) bool { return true })

	err := Run(waitStep)

	assert.Contains(t, err.Error(), "step1 error")
	assert.Equal(t, 3, waitStep.attempts)
}
