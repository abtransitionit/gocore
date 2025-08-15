/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

This file defines the default implementation using the standard Go log package.

*/

package logx

import (
	"log"
	"os"
)

// Name: stdLogger
// Type: structure
// Description: a concrete implementation of the Logger interface
// Notes:
// - uses the standard Go log package.
type stdLogger struct {
	logger *log.Logger
}

// Name: NewStdLogger
// Description: creates a new stdLogger instance.
// Returns:
// - Logger: an instance to the stdLogger that satisfies the Logger interfac
func NewStdLogger() Logger {
	return &stdLogger{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Name: Info
// Description: logs a message with an INFO prefix.
// Notes:
// - implements the method Info of the Logger interface.
func (s *stdLogger) Info(format string, v ...any) {
	s.logger.Printf("INFO: "+format, v...)
}

// Name: Error
// Description: logs a message with an ERROR prefix.
// Notes:
// - implements the method Error of the Logger interface.
func (s *stdLogger) Error(format string, v ...any) {
	s.logger.Printf("ERROR: "+format, v...)
}
