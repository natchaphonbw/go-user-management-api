package errors

import (
	"net/http"
)

type AppError struct {
	Code    int    // HTTP Status Code
	Message string // Human-readable message
	Err     error  // Raw error (optional)
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Helper functions
func BadRequest(message string, err error) *AppError {
	return New(http.StatusBadRequest, message, err)
}

func NotFound(message string, err error) *AppError {
	return New(http.StatusNotFound, message, err)
}

func InternalServer(message string, err error) *AppError {
	return New(http.StatusInternalServerError, message, err)
}

func Unautherized(message string, err error) *AppError {
	return New(http.StatusUnauthorized, message, err)
}

func Conflict(message string, err error) *AppError {
	return New(http.StatusConflict, message, err)
}
