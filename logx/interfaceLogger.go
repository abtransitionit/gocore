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
}
