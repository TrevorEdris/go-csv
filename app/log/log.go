package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(logLevelString string) *zap.Logger {
	logLevel, err := parseLogLevel(logLevelString)
	if err != nil {
		fmt.Printf("Failed to parse log level: %v\n", err)
		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		return logger
	}

	// Configure Zap logger options
	config := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(logLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:          "time",
			LevelKey:         "level",
			NameKey:          "logger",
			MessageKey:       "msg",
			EncodeTime:       zapcore.ISO8601TimeEncoder,
			EncodeLevel:      zapcore.CapitalLevelEncoder,
			EncodeCaller:     zapcore.ShortCallerEncoder,
			ConsoleSeparator: " ",
		},
	}

	// Create the logger instance
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	return logger
}

func parseLogLevel(logLevelString string) (zapcore.Level, error) {
	var logLevel zapcore.Level
	err := logLevel.UnmarshalText([]byte(logLevelString))
	if err != nil {
		return zapcore.InfoLevel, err
	}
	return logLevel, nil
}
