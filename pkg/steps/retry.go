package steps

import (
	"fmt"
	"riser/pkg/logger"
	"time"

	"github.com/pkg/errors"
)

type RetryStep struct {
	StepMeta
	stepFunc    CreateStepFunc
	currentStep Step
	maxAttempts int
	attempts    int
	retryFunc   ShouldRetryFunc
	// No need to export for configuration - just speeds up testing
	sleepTime time.Duration
}

type ShouldRetryFunc func(stepError error) bool
type CreateStepFunc func() Step

func AlwaysRetry() ShouldRetryFunc {
	return func(error) bool {
		return true
	}
}

func NewRetryStep(stepFunc CreateStepFunc, maxAttempts int, shouldRetry ShouldRetryFunc) *RetryStep {
	return &RetryStep{StepMeta: stepFunc().Meta(), stepFunc: stepFunc, maxAttempts: maxAttempts, retryFunc: shouldRetry, sleepTime: 1 * time.Second}
}

func (step *RetryStep) Exec() error {
	var stepErr error
	for step.attempts = 1; step.attempts < step.maxAttempts; step.attempts++ {
		step.currentStep = step.stepFunc()
		stepErr = step.currentStep.Exec()
		if stepErr == nil || !step.retryFunc(stepErr) {
			break
		}

		logger.Log().Verbose(fmt.Sprintf("Step %q failed and will be retried. Error: %v", step.Name, stepErr))
		time.Sleep(step.sleepTime)
	}

	return errors.Wrap(stepErr, fmt.Sprintf("failed after %d attempts", step.attempts))
}

func (step *RetryStep) Meta() StepMeta {
	return step.StepMeta
}

func (step *RetryStep) State(key string) interface{} {
	if step.currentStep == nil {
		return nil
	}
	return step.currentStep.State(key)
}
