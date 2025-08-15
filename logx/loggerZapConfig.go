/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

This file defines the different config concerning the Zap driver for the different env we want: dev or prod.

*/

// config.go
package logx

import (
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
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	config.EncoderConfig.EncodeTime = func(time.Time, zapcore.PrimitiveArrayEncoder) {
		// This empty function will skip timestamp encoding
	}
	// config.EncoderConfig.EncodeCaller = fixedWidthCallerEncoder

	return config
}

// Name: NewProdConfig
// Return:
// - zap.Config: a configuration instance for the production environment.
// Notes:
// - This is standard the default zap production config.
func NewProdConfig() zap.Config {
	// config := zap.NewProductionConfig()
	// config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // no color, uppercase
	// config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	return zap.NewProductionConfig()
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
