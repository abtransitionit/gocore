// File in gocore/logx/loggerZapConfig.go
/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

defines the different config concerning the Zap logging driver for the different env we want: dev or prod.

*/
package logx

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Name: NewDevConfig
// Return:
// - zap.Config: a configuration instance for the development environment.
// Notes:
// - It includes customizations like colorized levels and short caller paths.
func NewDevConfig() zap.Config {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	// cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	cfg.EncoderConfig.EncodeCaller = shortCallerWithLineEncoder
	cfg.EncoderConfig.EncodeTime = func(time.Time, zapcore.PrimitiveArrayEncoder) {
		// This empty function will skip timestamp encoding
	}
	// cfg.EncoderConfig.EncodeCaller = fixedWidthCallerEncoder
	// cfg.Development = true // enable stack traces for warnings and errors

	return cfg
}

// Name: NewProdConfig
// Return:
// - zap.Config: a configuration instance for the production environment.
// Notes:
// - This is standard the default zap production config.
func NewProdConfig() zap.Config {
	// cfg := zap.NewProductionConfig()
	// cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder       // Uppercase levels
	// cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder         // Standard timestamps
	// cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder       // Short file paths
	// cfg.Sampling = &zap.SamplingConfig{Initial: 100, Thereafter: 100} // Optional, avoids log spam
	// return cfg
	return zap.NewProductionConfig()
}

// Custom caller encoder to show file:line only (similar to Std logger)
func shortCallerWithLineEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	// enc.AppendString(fmt.Sprintf("%s:%d", caller.TrimmedPath(), caller.Line))
	enc.AppendString(fmt.Sprintf("%s ", caller.TrimmedPath()))
}

// func fixedWidthCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
// 	file := caller.TrimmedPath()
// 	funcParts := strings.Split(caller.Function, ".")
// 	funcName := funcParts[len(funcParts)-1]
// 	fileWidth := 30
// 	funcWidth := 20
// 	enc.AppendString(fmt.Sprintf("%-*s %-*s", fileWidth, file, funcWidth, funcName))
// }

// file := filepath.Base(caller.File)
// Extract just the function name
// funcName := caller.Function
// enc.AppendString(fmt.Sprintf("%10s:%3d %12s", file, caller.Line, funcName))
