package domain

import (
	"fmt"
)

type CustomError string

var (
	ErrorBadRequest                 CustomError = "bad_request"
	ErrorPurchaseNotInPendingStatus CustomError = "purchase_is_not_pending_status"
	ErrorPurchaseNotExist           CustomError = "purchase_not_exist"
	ErrorPurchaseAlreadyExists      CustomError = "purchase_already_exist"
	ErrorRouteNotFound              CustomError = "route_not_found"
	ErrorUnexpected                 CustomError = "unexpected_error"
)

type AppError struct {
	CustomError CustomError
	Message     string
	Metadata    map[string]any
	Err         error
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.CustomError, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) WithMetadata(key string, value any) *AppError {
	if e.Metadata == nil {
		e.Metadata = make(map[string]any)
	}
	e.Metadata[key] = value
	return e
}

func NewError(customErr CustomError, msg string, err error) *AppError {
	return &AppError{
		CustomError: customErr,
		Message:     msg,
		Err:         err,
	}
}
