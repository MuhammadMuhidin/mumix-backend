package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error is a custom error type for the application
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

// Error implements the error interface
func (e *Error) Error() string {
	return e.Message
}

// WriteJSON writes an error as JSON response
func WriteJSON(w http.ResponseWriter, err *Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(err)
}

// New creates a new Error
func New(code, message string, status int) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, code, message string, status int) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf("%s: %v", message, err),
		Status:  status,
	}
}

// Common error constructors
func NotFound(resource string) *Error {
	return New("NOT_FOUND", fmt.Sprintf("%s not found", resource), 404)
}

func BadRequest(message string) *Error {
	return New("BAD_REQUEST", message, 400)
}

func Unauthorized(message string) *Error {
	return New("UNAUTHORIZED", message, 401)
}

func Conflict(message string) *Error {
	return New("CONFLICT", message, 409)
}

func Internal(message string) *Error {
	return New("INTERNAL_ERROR", message, 500)
}
