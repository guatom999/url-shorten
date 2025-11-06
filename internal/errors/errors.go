package errors

import "errors"

// Custom error types
var (
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrInternal     = errors.New("internal server error")
	ErrConflict     = errors.New("resource already exists")
)

// Error types for checking
type ErrorType string

const (
	NotFound     ErrorType = "NOT_FOUND"
	InvalidInput ErrorType = "INVALID_INPUT"
	Internal     ErrorType = "INTERNAL"
	Conflict     ErrorType = "CONFLICT"
)

// AppError represents application error with type
type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// New error functions
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Type:    NotFound,
		Message: message,
	}
}

func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Type:    Internal,
		Message: message,
		Err:     err,
	}
}

func NewInvalidInputError(message string) *AppError {
	return &AppError{
		Type:    InvalidInput,
		Message: message,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Type:    Conflict,
		Message: message,
	}
}
