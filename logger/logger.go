package logger

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

// Event defines the basic structure of a type of log.
// It includes a id, and a message, that will be the default message
// to send.
type Event struct {
	id      int
	message string
}

// StandardLogger is a *log.Entry abstraction layer.
type StandardLogger struct {
	*log.Entry
}

// NewLogger creates a logger based on the app that is using it.
// the parameter tool will be appended to the log, so it can be
// identified correctly.
// It returns a *StandardLogger object
func NewLogger(tool string) *StandardLogger {
	log.SetHandler(cli.Default)
	log.SetLevel(log.DebugLevel)
	var baseLogger = log.WithFields(log.Fields{
		"tool": tool,
	})
	var standardLogger = &StandardLogger{baseLogger}

	return standardLogger
}

// Declare variables to store log messages as new Events
// Every log message should be declared here, so we follow
// a standard.
var (
	customMessage       = Event{0, "INFO: %s"}
	invalidArgMessage   = Event{1, "Invalid arg: %s"}
	invalidValueMessage = Event{2, "Invalid value for argument: %s: %v"}
	missingArgMessage   = Event{3, "Missing arg: %s"}
)

// InvalidArg is a standard info message
func (l *StandardLogger) ICustom(s string) {
	l.Infof(customMessage.message, s)
}

// InvalidArg is a standard error message
func (l *StandardLogger) ECustom(argumentName string) {
	l.Errorf(customMessage.message, argumentName)
}

// InvalidArg is a standard error message
func (l *StandardLogger) EInvalidArg(argumentName string) {
	l.Errorf(invalidArgMessage.message, argumentName)
}

// InvalidValue is a standard error message
func (l *StandardLogger) EInvalidValue(argumentName string, argumentValue string) {
	l.Errorf(invalidValueMessage.message, argumentName, argumentValue)
}

// MissingArg is a standard error message
func (l *StandardLogger) EMissingArg(argumentName string) {
	l.Errorf(missingArgMessage.message, argumentName)
}
