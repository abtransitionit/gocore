/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

This file defines a specific implementation using the Zap package.

*/

package logx

import (
	"fmt"

	"github.com/abtransitionit/gocore/errorx"
	"go.uber.org/zap"
)

// Name: stdLogger
// Description: a concrete implementation of the Logger interface
// Notes:
// - uses the Zap library logging
type zapLogger struct {
	logger *zap.Logger
}

// Name: NewZapLogger
// Description: creates a new zapLogger instance for any environment (production or development).
// Return:
// - Logger: an instance to the zapLogger configured for an env (dev, prod, ..) and that satisfies the Logger interface.
func NewZapLogger(config zap.Config) Logger {
	l, _ := config.Build()
	return &zapLogger{
		logger: l,
	}
}

// Name: Info
// Description: logs a message with the zap logger at Info level.
// Notes:
// - implements the method of the same name in the Logger interface.
func (z *zapLogger) Info(format string, v ...any) {
	z.logger.Info(fmt.Sprintf(format, v...))
}

// Name: Error
// Description: logs a message with the zap logger at Error level.
// Notes:
// - implements the method of the same name in the Logger interface.
func (l *zapLogger) Error(format string, v ...any) {
	l.logger.Error(fmt.Sprintf(format, v...))
}

// Name: ErrorWithStack
// Description: logs an error, including a stack trace if one is available.
// Notes:
// - implements the method of the same name in the Logger interface.
func (l *zapLogger) ErrorWithStack(err error, format string, v ...any) {
	// create a slice of zap.Field and append the error and stack trace to the original error (zap.Error(err)).
	fields := []zap.Field{
		zap.Error(err),
	}

	// Use GetStack from your errorx package to check for a stack trace.
	if stack := errorx.GetStack(err); stack != nil {
		fields = append(fields, zap.String("stack_trace", errorx.FormatStack(stack)))
	}

	l.logger.Error(fmt.Sprintf(format, v...), fields...)
}

// Name: ErrorWithNoStack
// Description: logs an error without including a stack trace.
// Notes:
// - implements the method of the same name in the Logger interface.
func (l *zapLogger) ErrorWithNoStack(err error, format string, v ...any) {
	l.logger.Error(fmt.Sprintf(format, v...), zap.Error(err))
}
