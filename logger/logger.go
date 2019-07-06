// Package logger provides a common logging interface.
package logger

// Logger interface for logging
type Logger interface {
	Info(string)
	Error(string)
	Verbose(string)
}

// logger is the global shared instance of Logger
var logger = Logger(NewScreenLogger(false))

// Log returns the default logger
func Log() Logger {
	return logger
}

// SetLogger sets the shared logger. This should only be done on process startup. This is not thread safe.
func SetLogger(l Logger) {
	logger = l
}
