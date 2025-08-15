/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

The main entry point for your application to interact with the logger.

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
		// Default to the standard logger driver if the variable is not set or unknown.
		GlobalLogger = NewStdLogger()
	}
}

// Name: Info
// Description:  a convenience function that delegates to the global logger's Info method.
func Info(format string, v ...any) {
	if GlobalLogger == nil {
		Init() // Ensure logger is initialized if not already
	}
	GlobalLogger.Info(format, v...)
}

// Name: Error
// Description: a convenience function that delegates to the global logger's Error method.
func Error(format string, v ...any) {
	if GlobalLogger == nil {
		Init() // Ensure logger is initialized if not already
	}
	GlobalLogger.Error(format, v...)
}
