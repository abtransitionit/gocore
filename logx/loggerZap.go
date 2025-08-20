// File: gocore/logx/loggerZap.go
/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

Implementation of Logger interface using Zap SugaredLogger.
*/

package logx

import (
	"fmt"

	"github.com/abtransitionit/gocore/errorx"
	"go.uber.org/zap"
)

// zapLogger implements Logger using Zap SugaredLogger
type zapLogger struct {
	logger *zap.SugaredLogger
}

// NewZapLogger creates a new zapLogger instance for any environment (prod/dev)
func NewZapLogger(config zap.Config) Logger {
	l, err := config.Build(zap.AddCallerSkip(2)) // skip frames for wrappers
	if err != nil {
		panic("failed to build zap logger: " + err.Error())
	}
	return &zapLogger{
		logger: l.Sugar(),
	}
}

// Simple info logging
func (l *zapLogger) Info(msg string) {
	l.logger.Info(msg)
}

// Formatted info logging
func (l *zapLogger) Infof(format string, v ...any) {
	l.logger.Infof(format, v...)
}

// Structured info logging
func (l *zapLogger) Infow(msg string, keysAndValues ...any) {
	l.logger.Infow(msg, keysAndValues...)
}

// Simple error logging
func (l *zapLogger) Error(msg string) {
	l.logger.Error(msg)
}

// Formatted error logging
func (l *zapLogger) Errorf(format string, v ...any) {
	l.logger.Errorf(format, v...)
}

// Structured error logging
func (l *zapLogger) Errorw(msg string, keysAndValues ...any) {
	l.logger.Errorw(msg, keysAndValues...)
}

// Error with stack trace
func (l *zapLogger) ErrorWithStack(err error, format string, v ...any) {
	fields := []any{"error", err}
	if stack := errorx.GetStack(err); stack != nil {
		fields = append(fields, "stack_trace", errorx.FormatStack(stack))
	}
	l.logger.Errorw(fmt.Sprintf(format, v...), fields...)
}

// Error without stack trace
func (l *zapLogger) ErrorWithNoStack(err error, format string, v ...any) {
	l.logger.Errorw(fmt.Sprintf(format, v...), "error", err)
}
