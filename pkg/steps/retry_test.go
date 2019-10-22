package steps

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_RetryStep_NoError_DoesNotRetry(t *testing.T) {
	step1 := newFakeStep("step1")
	waitStep := NewRetryStep(func() Step { return step1 }, 3, func(error) bool {
		assert.Fail(t, "RetryFunc should not be called")
		return false
	})

	err := Run(waitStep)

	assert.NoError(t, err)
	assert.Equal(t, 1, waitStep.attempts)
}

func Test_RetryStep_RetryFuncFalse_DoesNotRetry(t *testing.T) {
	step1 := newFakeStep("step1")
	step1.err = errors.New("step1 error")
	waitStep := NewRetryStep(func() Step { return step1 }, 3, func(error) bool { return false })

	err := Run(waitStep)

	assert.Contains(t, err.Error(), "step1 error")
	assert.Equal(t, 1, waitStep.attempts)
}

func Test_RetryStep_MaxRetries(t *testing.T) {
	stepFuncCalls := 0
	step1 := newFakeStep("step1")
	step1.err = errors.New("step1 error")
	waitStep := NewRetryStep(func() Step { stepFuncCalls++; return step1 }, 3, func(error) bool { return true })
	waitStep.sleepTime = 1 * time.Microsecond

	err := Run(waitStep)

	assert.Contains(t, err.Error(), "step1 error")
	assert.Contains(t, err.Error(), "failed after 3 attempts")
	assert.Equal(t, 3, waitStep.attempts)
	assert.Equal(t, 3, stepFuncCalls)
}
