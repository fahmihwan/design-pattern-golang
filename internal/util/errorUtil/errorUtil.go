package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Response represents the general response structure
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data.omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(errors interface{}) Response {
	return Response{
		Status:  http.StatusBadRequest,
		Message: "Error",
		Errors:  errors,
	}
}

// SendErrorResponse sends an error response with the provided status code
func SendErrorResponse(w http.ResponseWriter, statusCode int, errors interface{}) {
	errorResponse := NewErrorResponse(errors)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", e.Message)
}

func NewNotFound(format string, a ...any) *NotFoundError {
	return &NotFoundError{Message: fmt.Sprintf(format, a...)}
}

var (
	ErrNotFound *NotFoundError
)

var ErrPipelineNotApproved = errors.New("pipeline must be in approved state")
