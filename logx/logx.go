/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com
*/

package logx

import (
	"log"
	"os"
)

// Logger is the main logging instance for the application.
var Logger *log.Logger

// Init initializes the logger.
func Init() {
	Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
}

// Info logs a message with an INFO prefix.
func Info(format string, v ...any) {
	Logger.Printf("INFO: "+format, v...)
}

// Error logs a message with an ERROR prefix.
func Error(format string, v ...any) {
	Logger.Printf("ERROR: "+format, v...)
}
