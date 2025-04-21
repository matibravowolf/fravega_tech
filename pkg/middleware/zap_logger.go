package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uMakeMeCrazy/fravega_tech/pkg/logger"
	"go.uber.org/zap"
)

// WithLogger is a Gin middleware that injects a request-scoped logger
func WithLogger(baseLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID
		requestID := uuid.New().String()

		// Create a logger with request-specific fields
		requestLogger := baseLogger.With(
			zap.String("request_id", requestID),
		)

		// Inject the logger into context
		ctx := c.Request.Context()
		ctx = logger.Inject(ctx, requestLogger)
		ctx = logger.InjectRequestID(ctx, requestID)

		// Propagate the new context with the request
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
