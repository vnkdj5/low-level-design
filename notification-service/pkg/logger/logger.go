package logger

import (
	"lld/pkg/config"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

/**

Usage:

zap.L().Info("Some log message", zap.Any("key", "ValueOrStruct"))

**/

func getLoggerConfig() zap.Config {

	switch strings.ToLower(config.LogLevel) {
	case "debug":
		return zap.NewDevelopmentConfig()
	case "error":
		{
			config := zap.NewProductionConfig()
			config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
			return config
		}
	case "info":
		return zap.NewProductionConfig()
	default:
		return zap.NewProductionConfig()
	}
}

// Initializes the zap logger and sets it to the global logger instance
func InitLogger() *zap.Logger {

	config := getLoggerConfig()

	//config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	logger, _ = config.Build()
	zap.ReplaceGlobals(logger)

	return logger
}

// Return SingleTon Logger Instance
func Logger() *zap.Logger {
	if logger != nil {
		return logger
	} else {
		return InitLogger()
	}

}
