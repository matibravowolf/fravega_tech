package domain

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type responseError struct {
	CustomError CustomError    `json:"error"`
	Message     string         `json:"message"`
	Metadata    map[string]any `json:"metadata"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err == nil {
			return
		}

		appErr, ok := err.Err.(*AppError)
		if ok {
			status := mapErrorCodeToStatus(appErr.CustomError)
			c.JSON(status, responseError{
				CustomError: appErr.CustomError,
				Message:     appErr.Message,
				Metadata:    appErr.Metadata,
			})
			return
		}

		// fallback for unexpected errors
		c.JSON(http.StatusInternalServerError, responseError{
			CustomError: ErrorUnexpected,
			Message:     "unexpected error",
		})
	}
}

func mapErrorCodeToStatus(customErr CustomError) int {
	switch customErr {
	case ErrorRouteNotFound, ErrorPurchaseNotExist:
		return http.StatusNotFound
	case ErrorBadRequest, ErrorPurchaseAlreadyExists, ErrorPurchaseNotInPendingStatus:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
