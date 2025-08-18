/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

gocore/logx/interfaceLogger.go defines the interface any logger driver must implements.
*/
package logx

// Name: Logger
// Description: defines the logging methods used in the application that all loggers must implements.
// Notes:
// - This allows for easy swapping of the underlying logging implementation.
type Logger interface {
	Info(format string, v ...any)
	Error(format string, v ...any)
	//logs an error, including a stack trace if one is available.
	ErrorWithStack(err error, format string, v ...any)
	// logs an error without including a stack trace.
	ErrorWithNoStack(err error, format string, v ...any)
}
