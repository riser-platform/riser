package e2e

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/pkg/errors"
)

const defaultCommandTimeout = 20 * time.Second

func shellOrFailTimeout(t *testing.T, timeout time.Duration, format string, args ...interface{}) string {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return shellOrFailContext(t, ctx, format, args...)
}

func shellOrFail(t *testing.T, format string, args ...interface{}) string {
	return shellOrFailTimeout(t, defaultCommandTimeout, format, args...)
}

func shellOrFailContext(t *testing.T, ctx context.Context, format string, args ...interface{}) string {
	output, err := shellContext(ctx, format, args...)
	if err != nil {
		t.Fatalf("Shell command failed: %v", err)
	}

	return output
}

func shell(format string, args ...interface{}) (string, error) {
	return shellContext(context.Background(), format, args...)
}

func shellContext(ctx context.Context, format string, args ...interface{}) (string, error) {
	command := fmt.Sprintf(format, args...)
	c := exec.CommandContext(ctx, "sh", "-c", command)
	stdout, err := c.StdoutPipe()
	if err != nil {
		return "", errors.Wrap(err, "error getting stdout pipe")
	}
	c.Stderr = c.Stdout

	var output []byte

	// The exec package is broken when it comes to cancellation. Without this hack a long running process cannot be cancelled.
	// https://github.com/golang/go/issues/23019
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			output = append(output, scanner.Bytes()...)
		}
	}()

	err = c.Run()

	if err != nil {
		if ctx.Err() != nil {
			err = errors.Wrap(ctx.Err(), err.Error())
		}
		return string(output), fmt.Errorf("command %q failed: %q %v", command, string(output), err)
	}
	return string(output), nil
}
