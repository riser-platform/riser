package logger

type FakeLogger struct {
	VerboseLogs []string
	InfoLogs    []string
	WarnLogs    []string
	ErrorLogs   []string
}

func NewFakeLogger() *FakeLogger {
	return &FakeLogger{
		VerboseLogs: []string{},
		InfoLogs:    []string{},
		WarnLogs:    []string{},
		ErrorLogs:   []string{},
	}
}

func (logger *FakeLogger) Verbose(message string) {
	logger.VerboseLogs = append(logger.VerboseLogs, message)
}

func (logger *FakeLogger) Info(message string) {
	logger.InfoLogs = append(logger.InfoLogs, message)
}

func (logger *FakeLogger) Warn(message string) {
	logger.WarnLogs = append(logger.WarnLogs, message)
}

func (logger *FakeLogger) Error(message string) {
	logger.ErrorLogs = append(logger.ErrorLogs, message)
}
