/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

The main entry point for your application to interact with the logger.

gocore/logx/logx.go

*/

package logx

import (
	"os"
	"sync"
)

// Name: GlobalLogger
//
// Description: the main logging instance for the application,
//
// Notes:
//
// -  reference the GALI (Global Application Logger Instance) of type Logger (the interface)
// -  private variable
var globalLogger Logger

// Name: GetLogger
//
// Returns:
//
// - Logger: the global logger instance.
//
// Notes:
// - This is the canonical way to access the application's logger.
// - It is safe to call before Init, as it will return a nil logger.
func GetLogger() Logger {
	return globalLogger
}

// Name: NewLogger
//
// Return:
//
// - Logger: A Logger instance
//
// Notes:
//
// - The logger is choosen based on the following environment variables
//   - APP_LOG_DRIVER : to choose the logger driver
//   - APP_ENV : to choose the logger driver format
func NewLogger() Logger {
	// retrive the value of some environment variables
	appLogDriver := os.Getenv("APP_LOG_DRIVER")
	appEnv := os.Getenv("APP_ENV")

	// instanciate the logger
	switch appLogDriver {
	case "zap":
		// Choose the zap logger configuration based on the environment
		if appEnv == "prod" {
			return NewZapLogger(NewProdConfig())
		} else {
			return NewZapLogger(NewDevConfig())
		}
	default:
		if appEnv == "prod" {
			// Default logger driver if only this variable is set (Std logger driver and prod config).
			return NewStdLogger(NewStdProdConfig())
		} else {
			// Default logger driver if both variables are not set or unknown (Std logger driver and dev config).
			return NewStdLogger(NewStdDevConfig())
		}
	}
}

// Name: once
//
// Description: sync.Once variable to ensure that the global logger is initialized only once, even with concurrent calls.
var once sync.Once

// Name: InitLogger
//
// Description: initializes the GALI
//
// Notes:
//
// - The logger is choosen based on the following environment variables
//   - APP_LOG_DRIVER : to choose the logger driver
//   - APP_ENV : to choose the logger driver format
func InitLogger() {
	globalLogger = NewLogger()
}

// func InitLogger() {

// 	// retrive the value of some environment variables
// 	appLogDriver := os.Getenv("APP_LOG_DRIVER")
// 	appEnv := os.Getenv("APP_ENV")

// 	// instanciate the logger
// 	switch appLogDriver {
// 	case "zap":
// 		// Choose the zap logger configuration based on the environment
// 		if appEnv == "prod" {
// 			globalLogger = NewZapLogger(NewProdConfig())
// 		} else {
// 			globalLogger = NewZapLogger(NewDevConfig())
// 		}
// 	default:
// 		if appEnv == "prod" {
// 			// Default logger driver if only this variable is set (Std logger driver and prod config).
// 			globalLogger = NewStdLogger(NewStdProdConfig())
// 		} else {
// 			// Default logger driver if both variables are not set or unknown (Std logger driver and dev config).
// 			globalLogger = NewStdLogger(NewStdDevConfig())
// 		}
// 	}
// }

// Name: init
//
// Description: The application's logger is initialization method
func Init() {
	once.Do(InitLogger)
}

// Name: Info
// Description:  a convenience function that delegates to the global logger's Info method.
// Notes:
// - a convenience for the function of the same name in the interface
func Info(format string, v ...any) {
	// if GlobalLogger == nil {
	// 	Init() // Ensure logger is initialized if not already
	// }
	globalLogger.Info(format, v...)
}

// Name: Error
// Description: a convenience function that delegates to the global logger's Error method.
// Notes:
// - a convenience for the function of the same name in the interface
func Error(format string, v ...any) {
	// if GlobalLogger == nil {
	// 	Init() // Ensure logger is initialized if not already
	// }
	globalLogger.Error(format, v...)
}

// Name: ErrorWithStack
// Description: is a convenience function that delegates to the global logger's ErrorWithStack method.
// Notes:
// - a convenience for the function of the same name in the interface
func ErrorWithStack(err error, format string, v ...any) {
	// if GlobalLogger == nil {
	// 	Init()
	// }
	globalLogger.ErrorWithStack(err, format, v...)
}

// Name: ErrorWithNoStack
// Description: is a convenience function that delegates to the global logger's ErrorWithNoStack method.
// Notes:
// - a convenience for the function of the same name in the interface
func ErrorWithNoStack(err error, format string, v ...any) {
	// if GlobalLogger == nil {
	// 	Init()
	// }
	globalLogger.ErrorWithNoStack(err, format, v...)
}
