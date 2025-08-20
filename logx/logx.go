// File to create in gocore/logx/logx.go
/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

define The entry point for any application to interact with the logger.

*/

package logx

import (
	"os"
	"sync"
)

// Name: GlobalLogger
//
// Description: the main logging instance for the application.
//
// Notes:
// - reference the GALI (Global Application Logger Instance) of type Logger (the interface)
// - private variable
var globalLogger Logger

// Name: once
//
// Description: ensure that the global logger is initialized only once, even with concurrent calls.
var once sync.Once

// Name: GetLogger
//
// Returns:
// - Logger: the global logger instance.
//
// Notes:
// - This is the canonical way to access the application's logger.
// - If not initialized yet, it will trigger Init().
func GetLogger() Logger {
	if globalLogger == nil {
		Init()
	}
	return globalLogger
}

// Name: NewLogger
//
// Return:
// - Logger: A Logger instance
//
// Notes:
// - The logger is chosen based on the following environment variables
//   - APP_LOG_DRIVER : to choose the logger driver
//   - APP_ENV : to choose the logger driver format
func NewLogger() Logger {
	// get environment variables
	appLogDriver := os.Getenv("APP_LOG_DRIVER")
	appEnv := os.Getenv("APP_ENV")

	// instanciate the logger
	switch appLogDriver {
	case "zap":
		if appEnv == "prod" {
			return NewZapLogger(NewProdConfig())
		}
		return NewZapLogger(NewDevConfig())
	default:
		if appEnv == "prod" {
			return NewStdLogger(NewStdProdConfig())
		}
		return NewStdLogger(NewStdDevConfig())
	}
}

// Name: InitLogger
//
// Description: initializes the GALI
func InitLogger() {
	globalLogger = NewLogger()
}

// Name: Init
//
// Description: ensures logger is initialized only once, even with concurrent calls.
func Init() {
	once.Do(InitLogger)
}

// Name: Info
//
// Description: a convenience function
//
// Notes:
// - ensures lazy initialization if globalLogger is nil.
func Info(format string, v ...any) {
	if globalLogger == nil {
		Init()
	}
	globalLogger.Info(format, v...)
}

// Name: Error
//
// Description: a convenience function
//
// Notes:
// - ensures lazy initialization if globalLogger is nil.
func Error(format string, v ...any) {
	if globalLogger == nil {
		Init()
	}
	globalLogger.Error(format, v...)
}

// Name: ErrorWithStack
//
// Description: a convenience function
//
// Notes:
// - ensures lazy initialization if globalLogger is nil.
func ErrorWithStack(err error, format string, v ...any) {
	if globalLogger == nil {
		Init()
	}
	globalLogger.ErrorWithStack(err, format, v...)
}

// Name: ErrorWithNoStack
// Description: a convenience function
// Notes:
// - ensures lazy initialization if globalLogger is nil.
func ErrorWithNoStack(err error, format string, v ...any) {
	if globalLogger == nil {
		Init()
	}
	globalLogger.ErrorWithNoStack(err, format, v...)
}
