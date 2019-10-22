package steps

import (
	"bytes"
	"fmt"
	"os/exec"
)

type ExecStep struct {
	StepMeta
	cmd   *exec.Cmd
	state map[string]interface{}
}

func NewExecStep(name string, cmd *exec.Cmd) *ExecStep {
	return &ExecStep{StepMeta: StepMeta{name}, cmd: cmd, state: map[string]interface{}{}}
}

// NewShellExecStep is a convenience wrapper to execute a command using sh.
// WARNING: This may not be cross platform (need to validate with linux for Windows subsystem)
func NewShellExecStep(name, shellCmd string) *ExecStep {
	return NewExecStep(name, exec.Command("sh", "-c", shellCmd))
}

func (step *ExecStep) Exec() error {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	step.cmd.Stdout = stdout
	step.cmd.Stderr = stderr

	// TODO: Stream stdout if in verbose mode

	err := step.cmd.Run()
	step.state["stdout"] = stdout.String()
	step.state["stderr"] = stderr.String()

	if err != nil {
		stderrMsg := fmt.Sprintf("\n%s", stderr)
		if stderr.Len() == 0 {
			stderrMsg = "<nil>"
		}
		return fmt.Errorf(fmt.Sprintf("%s. Stderr: %s", err.Error(), stderrMsg))
	}

	return nil
}

func (step *ExecStep) Meta() StepMeta {
	return step.StepMeta
}

func (step *ExecStep) State(key string) interface{} {
	return step.state[key]
}
