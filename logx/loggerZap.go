/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

This file defines a specific implementation using the Zap package.

*/

package logx

import (
	"fmt"

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
// - implementation for the method in the Logger interface.
func (z *zapLogger) Info(format string, v ...any) {
	z.logger.Info(fmt.Sprintf(format, v...))
}

// Name: Error
// Description: logs a message with the zap logger at Error level.
// Notes:
// - implementation for the method in the Logger interface.
func (z *zapLogger) Error(format string, v ...any) {
	z.logger.Error(fmt.Sprintf(format, v...))
}
