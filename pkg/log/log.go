package log

import (
	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerKey struct{}

var (
	defaultLogger *zap.Logger
)

type (
	noopCallerEncoder struct{}
	noopTimeEncoder   struct{}
)

func FromCtx(ctx context.Context) *zap.Logger {
	if ctx != nil {
		if logger, ok := ctx.Value(loggerKey{}).(*zap.Logger); ok {
			return logger
		}
	}
	return defaultLogger
}

func ToCtx(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func init() {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.EncodeTime = noTime
	cfg.EncoderConfig.EncodeCaller = omitCaller

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	defaultLogger = logger
}

func omitCaller(ec zapcore.EntryCaller, pae zapcore.PrimitiveArrayEncoder) {}

func noTime(t time.Time, pae zapcore.PrimitiveArrayEncoder) {}
