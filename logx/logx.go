/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

The main entry point for your application to interact with the logger.

gocore/logx/logx.go

*/

package logx

import "os"

// Name: GlobalLogger
// Description: the main logging instance for the application,
// Notes:
// -  reference the GALI (Global Application Logger Instance) of type Logger (the interface)
var GlobalLogger Logger

// Name: Init
// Description: initializes the GALI
// Notes: The logger is choosen based on the following environment variables
// - APP_LOG_DRIVER : to choose the logger driver
// - APP_ENV : to choose the logger driver format
func Init() {

	// retrive the value of some environment variables
	appLogDriver := os.Getenv("APP_LOG_DRIVER")
	appEnv := os.Getenv("APP_ENV")

	// instanciate the logger
	switch appLogDriver {
	case "zap":
		// Choose the zap logger configuration based on the environment
		if appEnv == "prod" {
			GlobalLogger = NewZapLogger(NewProdConfig())
		} else {
			GlobalLogger = NewZapLogger(NewDevConfig())
		}
	default:
		if appEnv == "prod" {
			// Default logger driver if only this variable is set (Std logger driver and prod config).
			GlobalLogger = NewStdLogger(NewStdProdConfig())
		} else {
			// Default logger driver if both variables are not set or unknown (Std logger driver and dev config).
			GlobalLogger = NewStdLogger(NewStdDevConfig())
		}
	}
}

// Name: Info
// Description:  a convenience function that delegates to the global logger's Info method.
// Notes:
// - a convenience for the function of the same name in the interface
func Info(format string, v ...any) {
	if GlobalLogger == nil {
		Init() // Ensure logger is initialized if not already
	}
	GlobalLogger.Info(format, v...)
}

// Name: Error
// Description: a convenience function that delegates to the global logger's Error method.
// Notes:
// - a convenience for the function of the same name in the interface
func Error(format string, v ...any) {
	if GlobalLogger == nil {
		Init() // Ensure logger is initialized if not already
	}
	GlobalLogger.Error(format, v...)
}

// Name: ErrorWithStack
// Description: is a convenience function that delegates to the global logger's ErrorWithStack method.
// Notes:
// - a convenience for the function of the same name in the interface
func ErrorWithStack(err error, format string, v ...any) {
	if GlobalLogger == nil {
		Init()
	}
	GlobalLogger.ErrorWithStack(err, format, v...)
}

// Name: ErrorWithNoStack
// Description: is a convenience function that delegates to the global logger's ErrorWithNoStack method.
// Notes:
// - a convenience for the function of the same name in the interface
func ErrorWithNoStack(err error, format string, v ...any) {
	if GlobalLogger == nil {
		Init()
	}
	GlobalLogger.ErrorWithNoStack(err, format, v...)
}
