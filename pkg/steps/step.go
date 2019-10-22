package steps

import (
	"fmt"
	"riser/pkg/logger"

	"github.com/pkg/errors"
	"github.com/wzshiming/ctc"
)

type Step interface {
	Meta() StepMeta
	Exec() error
	State(key string) interface{}
}

type StepMeta struct {
	Name string
}

func Run(steps ...Step) error {
	for _, step := range steps {
		logger.Log().Info(styleStepExecute(step.Meta()))
		err := step.Exec()
		if err != nil {
			return errors.Wrap(err, styleStepError(step.Meta()))
		}
		logger.Log().Info(styleStepComplete(step.Meta()))
	}

	return nil
}

func styleStepExecute(stepMeta StepMeta) string {
	return fmt.Sprintf(
		"Executing %s...",
		fmt.Sprint(ctc.ForegroundBrightCyan, stepMeta.Name, ctc.Reset),
	)
}

func styleStepComplete(stepMeta StepMeta) string {
	return fmt.Sprint(ctc.ForegroundBrightGreen, "✔", ctc.Reset, " Complete")
}

func styleStepError(stepMeta StepMeta) string {
	return fmt.Sprint(ctc.ForegroundBrightRed, "✘", ctc.Reset, " Error executing step ", ctc.ForegroundBrightRed, stepMeta.Name, ctc.Reset)
}
