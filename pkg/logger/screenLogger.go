package logger

import (
	"fmt"
	"riser/pkg/ui/style"
)

// ScreenLogger logs to a terminal
type ScreenLogger struct {
	VerboseMode bool
}

// NewScreenLogger creates a logger instance designed for printing messages to the screen for humans
func NewScreenLogger(verbose bool) *ScreenLogger {
	return &ScreenLogger{VerboseMode: verbose}
}

// Verbose logs a verbose message
func (logger *ScreenLogger) Verbose(message string) {
	if logger.VerboseMode {
		fmt.Println(style.Muted(message))
	}
}

// Info logs an information message to the screen
func (logger *ScreenLogger) Info(message string) {
	fmt.Println(message)
}

// Warn logs a warning message to the screen
func (logger *ScreenLogger) Warn(message string) {
	fmt.Println(style.Warn(message))
}

// Error logs an error message to the screen
func (logger *ScreenLogger) Error(message string) {
	fmt.Println(style.Bad(message))
}
