package ui

import (
	"fmt"
	"os"
	"riser/pkg/logger"
)

// ExitIfError exits if the error is not null and prints the error message
func ExitIfError(err error) {
	if err != nil {
		ExitErrorMsg(err.Error())
	}
}

// ExitIfErrorMsg exits if the error is not null and prints a custom message with the error message
func ExitIfErrorMsg(err error, exitMessage string) {
	if err != nil {
		ExitErrorMsg(fmt.Sprintf("%s: %s", exitMessage, err))
	}
}

// ExitErrorMsg exits with error code 1 and prints a custom message
func ExitErrorMsg(exitMessage string) {
	logger.Log().Error(exitMessage)
	os.Exit(1)
}
