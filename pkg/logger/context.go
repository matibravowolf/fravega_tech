package logger

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey string

const (
	ctxKeyLogger    ctxKey = "logger"
	ctxKeyRequestID ctxKey = "request_id"
)

func Inject(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKeyLogger, logger)
}

func InjectRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ctxKeyRequestID, requestID)
}

func From(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(ctxKeyLogger).(*zap.Logger); ok {
		return logger
	}
	return zap.NewNop()
}

func RequestID(ctx context.Context) string {
	if id, ok := ctx.Value(ctxKeyRequestID).(string); ok {
		return id
	}
	return ""
}
