// File gocore/logx/interfaceLogger.go
/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

defines the interface any logger driver must implements.

*/
package logx

// Name: Logger
//
// Description: defines the logging methods used in the application that all loggers must implements.

// Notes:
// - This allows for easy swapping of the underlying logging implementation.
type Logger interface {
	// Info logs a plain informational message
	Info(msg string)

	// Infof logs a formatted informational message
	Infof(format string, v ...any)

	// Structured informational logging with key-value pairs
	Infow(msg string, keysAndValues ...any)

	// Error logs a plain error message
	Error(msg string)

	// Structured error logging with key-value pairs
	Errorw(msg string, keysAndValues ...any)

	// Errorf logs a formatted error message
	Errorf(format string, v ...any)

	// ErrorWithStack logs an error with its stack trace
	ErrorWithStack(err error, format string, v ...any)

	// ErrorWithNoStack logs an error without the stack trace
	ErrorWithNoStack(err error, format string, v ...any)

	// Warn logs a plain warning message
	Warn(msg string)

	// Warnf logs a formatted warning message
	Warnf(format string, v ...any)

	// Warnw logs a structured warning message
	Warnw(msg string, keysAndValues ...any)

	// Debug logs a plain debug message
	Debug(msg string)

	// Debugf logs a formatted debug message
	Debugf(format string, v ...any)

	// Debugw logs a structured debug message
	Debugw(msg string, keysAndValues ...any)
}
