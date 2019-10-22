package steps

type FuncStep struct {
	StepMeta
	stepFunc func() error
}

func NewFuncStep(name string, stepFunc func() error) *FuncStep {
	return &FuncStep{StepMeta: StepMeta{name}, stepFunc: stepFunc}
}

func (step *FuncStep) Exec() error {
	return step.stepFunc()
}

func (step *FuncStep) Meta() StepMeta {
	return step.StepMeta
}

func (step *FuncStep) State(key string) interface{} {
	// For the moment we don't have a way for the underlying func to modify state
	return nil
}
