// File gocore/logx/loggerStd.go
/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

defines the default implementation using the standard Go log package.

*/

package logx

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/abtransitionit/gocore/errorx"
)

// ANSI color codes
const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorWhite   = "\033[37m"

	// Bright versions
	colorBrightRed     = "\033[91m"
	colorBrightGreen   = "\033[92m"
	colorBrightYellow  = "\033[93m"
	colorBrightBlue    = "\033[94m" // <-- light blue
	colorBrightMagenta = "\033[95m"
	colorBrightCyan    = "\033[96m"
	colorBrightWhite   = "\033[97m"
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

func pathFormatter(fullPath string) string {
	parts := strings.Split(filepath.ToSlash(fullPath), "/")
	if len(parts) <= 3 {
		return fullPath
	}
	return strings.Join(parts[len(parts)-3:], "/")
}

// Info logs a simple info message
func (l *stdLogger) Info(msg string) {
	// l.logger.Println("INFO:", msg)
	l.logger.Output(2, colorCyan+"INFO:   "+colorReset+msg)

}

// Infof logs a formatted info message
func (l *stdLogger) Infof(format string, v ...any) {
	l.logger.Output(2, colorCyan+"INFO:   "+fmt.Sprintf(format, v...)+colorReset)
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
	l.logger.Output(2, "INFO:   "+msg)
}

// Error logs a simple error message
func (l *stdLogger) Error(msg string) {
	l.logger.Output(2, colorRed+"ERROR:   "+colorReset+msg)

}

// Errorf logs a formatted error message
func (l *stdLogger) Errorf(format string, v ...any) {
	l.logger.Output(2, colorRed+"ERROR:   "+colorReset+fmt.Sprintf(format, v...))

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
	l.logger.Output(2, colorRed+"ERROR:   "+colorReset+colorReset+msg)
}

// ErrorWithStack logs an error with a stack trace if available
func (l *stdLogger) ErrorWithStack(err error, format string, v ...any) {
	var sb strings.Builder
	// sb.WriteString("ERROR: ")
	sb.WriteString(fmt.Sprintf(format, v...))
	sb.WriteString(": ")
	sb.WriteString(err.Error())
	sb.WriteString("\n")

	if stack := errorx.GetStack(err); stack != nil {
		sb.WriteString(errorx.FormatStack(stack))
	}

	l.logger.Output(2, colorRed+"ERROR:   "+colorReset+sb.String())
}

// ErrorWithNoStack logs an error without the stack trace
func (l *stdLogger) ErrorWithNoStack(err error, format string, v ...any) {
	msg := fmt.Sprintf(format, v...) + ": " + err.Error()
	l.logger.Output(2, "ERROR:   "+colorReset+msg)

}

// File gocore/logx/loggerStd.go

// Warn logs a simple warning message
func (l *stdLogger) Warn(msg string) {
	l.logger.Output(2, colorYellow+"WARN:   "+msg+colorReset)
}

// Warnf logs a formatted warning message
func (l *stdLogger) Warnf(format string, v ...any) {
	l.logger.Output(2, colorYellow+"WARN:   "+fmt.Sprintf(format, v...)+colorReset)
}

func (l *stdLogger) Warnw(msg string, keysAndValues ...any) {
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
	l.logger.Output(2, colorYellow+"WARN:   "+msg+colorReset)
}

// Debug logs a simple debug message
func (l *stdLogger) Debug(msg string) {
	l.logger.Output(2, colorBrightBlue+"DEBUG:   "+colorReset+msg)
}

// Debugf logs a formatted debug message
func (l *stdLogger) Debugf(format string, v ...any) {
	l.logger.Output(2, colorBrightBlue+"DEBUG:   "+colorReset+fmt.Sprintf(format, v...))
}

func (l *stdLogger) Debugw(msg string, keysAndValues ...any) {
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
	l.logger.Output(2, colorBrightBlue+"DEBUG:   "+colorReset+msg)
}
