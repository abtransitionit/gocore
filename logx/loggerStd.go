// File gocore/logx/loggerStd.go
/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

defines the default implementation using the standard Go log package.

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

// Info logs a simple info message
func (l *stdLogger) Info(msg string) {
	l.logger.Println("INFO:", msg)
}

// Infof logs a formatted info message
func (l *stdLogger) Infof(format string, v ...any) {
	l.logger.Printf("INFO: "+format, v...)
}

func (l *stdLogger) Infow(msg string, keysAndValues ...any) {
	if len(keysAndValues) > 0 {
		msg += " | "
		for i := 0; i < len(keysAndValues); i += 2 {
			k := keysAndValues[i]
			v := "<nil>"
			if i+1 < len(keysAndValues) {
				v = fmt.Sprint(keysAndValues[i+1])
			}
			msg += fmt.Sprintf("%v=%v ", k, v)
		}
	}
	l.logger.Println("INFO:", msg)
}

// Error logs a simple error message
func (l *stdLogger) Error(msg string) {
	l.logger.Println("ERROR:", msg)
}

// Errorf logs a formatted error message
func (l *stdLogger) Errorf(format string, v ...any) {
	l.logger.Printf("ERROR: "+format, v...)
}

func (l *stdLogger) Errorw(msg string, keysAndValues ...any) {
	if len(keysAndValues) > 0 {
		msg += " | "
		for i := 0; i < len(keysAndValues); i += 2 {
			k := keysAndValues[i]
			v := "<nil>"
			if i+1 < len(keysAndValues) {
				v = fmt.Sprint(keysAndValues[i+1])
			}
			msg += fmt.Sprintf("%v=%v ", k, v)
		}
	}
	l.logger.Println("ERROR:", msg)
}

// ErrorWithStack logs an error with a stack trace if available
func (l *stdLogger) ErrorWithStack(err error, format string, v ...any) {
	var sb strings.Builder
	sb.WriteString("ERROR: ")
	sb.WriteString(fmt.Sprintf(format, v...))
	sb.WriteString(": ")
	sb.WriteString(err.Error())
	sb.WriteString("\n")

	if stack := errorx.GetStack(err); stack != nil {
		sb.WriteString(errorx.FormatStack(stack))
	}

	l.logger.Println(sb.String())
}

// ErrorWithNoStack logs an error without the stack trace
func (l *stdLogger) ErrorWithNoStack(err error, format string, v ...any) {
	l.logger.Printf("ERROR: "+format, v...)
	l.logger.Println("Original Error:", err.Error())
}
