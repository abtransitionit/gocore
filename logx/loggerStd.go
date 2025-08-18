/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

gocore/logx/interfaceLogger.go defines the default implementation using the standard Go log package.

*/

package logx

import (
	"fmt"
	"log"
	"strings"

	"github.com/abtransitionit/gocore/errorx"
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
// - Logger: an instance to the stdLogger that satisfies the Logger using the provided configuration.
func NewStdLogger(config StdLoggerConfig) Logger {
	return &stdLogger{
		logger: log.New(config.Out, config.Prefix, config.Flag),
	}
}

// func NewStdLogger() Logger {
// 	return &stdLogger{
// 		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),
// 	}
// }

// Name: Info
// Description: logs a message with an INFO prefix.
// Notes:
// - implements the method of the same name in the Logger interface.
func (l *stdLogger) Info(format string, v ...any) {
	l.logger.Printf("INFO: "+format, v...)
}

// Name: Error
// Description: logs a message with an ERROR prefix.
// Notes:
// - implements the method of the same name in the Logger interface.
func (l *stdLogger) Error(format string, v ...any) {
	l.logger.Printf("ERROR: "+format, v...)
}

// Name: ErrorWithStack
// Description: logs an error, including a stack trace if one is available.
// Notes:
// - implements the method of the same name in the Logger interface.
func (l *stdLogger) ErrorWithStack(err error, format string, v ...any) {
	var sb strings.Builder
	sb.WriteString("ERROR: ")
	sb.WriteString(fmt.Sprintf(format, v...))
	sb.WriteString(": ")
	sb.WriteString(err.Error())
	sb.WriteString("\n")

	// Use GetStack from your errorx package to check for a stack trace.
	if stack := errorx.GetStack(err); stack != nil {
		sb.WriteString(errorx.FormatStack(stack))
	}

	l.logger.Println(sb.String())
}

// Name: ErrorWithNoStack
// Description: logs an error without including a stack trace.
// Notes:
// - implements the method of the same name in the Logger interface.
func (l *stdLogger) ErrorWithNoStack(err error, format string, v ...any) {
	l.logger.Printf("ERROR: "+format, v...)
	l.logger.Println("Original Error:", err.Error())
}
