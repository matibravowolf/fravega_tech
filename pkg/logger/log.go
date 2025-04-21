package logger

import (
	"context"

	"go.uber.org/zap"
)

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	From(ctx).Info(msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	From(ctx).Warn(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	From(ctx).Error(msg, fields...)
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	From(ctx).Debug(msg, fields...)
}
