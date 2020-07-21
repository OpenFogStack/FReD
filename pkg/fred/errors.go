package fred

import (
	"fmt"
	"net/http"
)

// Error is a custom error type for FReD that also has an HTTP status code.
type Error struct {
	Code  int
	error string
}

// Status Codes we use in FReD that can be used when creating new custom errors.
const (
	StatusConflict      int = http.StatusConflict
	StatusBadRequest    int = http.StatusBadRequest
	StatusNotFound      int = http.StatusNotFound
	StatusInternalError int = http.StatusInternalServerError
)

// Error is needed to satisfy the interface of error for Error.
func (e *Error) Error() string {
	return fmt.Sprintf("%s (Status Code %d)", e.error, e.Code)
}

// newError creates a new custom Error.
func newError(code int, error string) *Error {
	return &Error{
		Code:  code,
		error: error,
	}
}
