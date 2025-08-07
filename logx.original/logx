/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com
*/

// provides logging with Uber zap logger
//
// Features:
//   - Dev and prod mode
//
// Usage:
//
//  1. Initialize logger (in main.go):
//     - logx.Init(true) // true for dev mode, false for prod
//  2. use logger in code:
//     - logx.L.Info("application starting")
//     - logx.L.Infof("installing %s", packageName)
//     - logx.L.Errorf("failed to start: %v", err)
//     - ...
package logx

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// L is the global SugaredLogger instance, initialized by Init() and used in the code.
var L *zap.SugaredLogger

// initializes the global logger with either development or production configuration.
//
// Parameters:
//   - dev : 	When true, uses  dev config (colored output, simplified format)
//     When false, uses prod config (structured, optimized for logging systems)
//
// Behavior:
//   - Panics if logger creation fails (logger initialization is considered critical)
//   - Configures different settings for dev/prod environments
//   - Sets global L variable for package-wide access
func Init(dev bool) {
	var (
		base *zap.Logger
		cfg  zap.Config
		err  error
	)

	// define the type of logger according to the bool
	if dev {
		cfg = devConfig()
	} else {
		cfg = prodConfig()
	}

	// build the logger
	base, err = cfg.Build()

	// error: panic : logger creation MUST succeed
	if err != nil {
		panic(err)
	}
	// Converts the base zap.Logger to a SugaredLogger and assigns it to L
	L = base.Sugar()
}

// devConfig returns a development-optimized zap configuration.
//
// Features:
//   - Colored log levels
//   - No timestamps (for cleaner console output)
//   - Fixed-width caller formatting
//   - Optimized for human readability
func devConfig() zap.Config {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // colored levels for dev
	cfg.EncoderConfig.EncodeTime = func(time.Time, zapcore.PrimitiveArrayEncoder) {
		// This empty function will skip timestamp encoding
	}
	// cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        // readable time
	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder // short caller info
	cfg.EncoderConfig.EncodeCaller = fixedWidthCallerEncoder

	return cfg
}

// prodConfig returns a production-optimized zap configuration.
//
// Features:
//   - ISO8601 timestamps
//   - Structured JSON output
//   - Short caller information
//   - Optimized for log aggregation systems
func prodConfig() zap.Config {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // no color, uppercase
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return cfg
}

// fixedWidthCallerEncoder provides consistent-width caller information formatting.
//
// Format:
//
//	"[filename:line][functionName]"
//	- filename: 10 characters wide
//	- line: 3 digits
//	- functionName: 32 characters wide
//
// Example:
//
//	"main.go:023       handleRequest"
func fixedWidthCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	file := caller.TrimmedPath()
	funcParts := strings.Split(caller.Function, ".")
	funcName := funcParts[len(funcParts)-1]
	fileWidth := 30
	funcWidth := 20
	enc.AppendString(fmt.Sprintf("%-*s %-*s", fileWidth, file, funcWidth, funcName))
}

// file := filepath.Base(caller.File)
// Extract just the function name
// funcName := caller.Function
// enc.AppendString(fmt.Sprintf("%10s:%3d %12s", file, caller.Line, funcName))
