package logx

// Name: Logger
// Description: defines the logging methods used in the application that all loggers must implements.
// Notes:
// - This allows for easy swapping of the underlying logging implementation.
type Logger interface {
	Info(format string, v ...any)
	Error(format string, v ...any)
}
