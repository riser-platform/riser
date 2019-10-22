package steps

import (
	"fmt"
	"riser/pkg/logger"
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
)

func Test_Run(t *testing.T) {
	step1 := newFakeStep("step1")
	step2 := newFakeStep("step2")
	log := logger.NewFakeLogger()
	logger.SetLogger(log)

	err := Run(step1, step2)

	assert.NoError(t, err)
	assert.True(t, step1.State("executed").(bool))
	assert.True(t, step2.State("executed").(bool))
	assert.Len(t, log.InfoLogs, 4)
	assert.Equal(t, styleStepExecute(step1.Meta()), log.InfoLogs[0])
	assert.Equal(t, styleStepComplete(step1.Meta()), log.InfoLogs[1])
	assert.Equal(t, styleStepExecute(step2.Meta()), log.InfoLogs[2])
	assert.Equal(t, styleStepComplete(step2.Meta()), log.InfoLogs[3])
}

func Test_Run_ErrorHalts(t *testing.T) {
	step1 := newFakeStep("step1")
	step1.err = errors.New("step1 error")
	step2 := newFakeStep("step2")
	logger.SetLogger(logger.NewFakeLogger())

	err := Run(step1, step2)

	assert.Equal(t, fmt.Sprintf("%s: step1 error", styleStepError(step1.Meta())), err.Error())
	assert.Nil(t, step2.State("executed"))
}

type fakeStep struct {
	StepMeta
	state map[string]interface{}
	err   error
}

func newFakeStep(name string) *fakeStep {
	return &fakeStep{StepMeta: StepMeta{name}, state: map[string]interface{}{}}
}

func (step *fakeStep) Meta() StepMeta {
	return step.StepMeta
}

func (step *fakeStep) Exec() error {
	if step.err != nil {
		return step.err
	}
	step.state["executed"] = true
	return nil
}

func (step *fakeStep) State(key string) interface{} {
	return step.state[key]
}
