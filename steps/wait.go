package steps

type WaitStep struct {
	StepMeta
	inner       Step
	maxAttempts int
	attempts    int
	retryFunc   WaitStepRetryFunc
}

type WaitStepRetryFunc func(stepError error) bool

func NewWaitStep(step Step, maxAttempts int, retryFunc WaitStepRetryFunc) *WaitStep {
	return &WaitStep{StepMeta: step.Meta(), inner: step, maxAttempts: maxAttempts, retryFunc: retryFunc}
}

func (step *WaitStep) Exec() error {
	var stepErr error
	for step.attempts = 1; step.attempts < step.maxAttempts; step.attempts++ {
		stepErr = step.inner.Exec()
		if stepErr == nil || !step.retryFunc(stepErr) {
			break
		}
	}

	return stepErr
}

func (step *WaitStep) Meta() StepMeta {
	return step.StepMeta
}

func (step *WaitStep) State(key string) interface{} {
	return step.inner.State(key)
}
