package rest

import (
	"fmt"
	"github.com/google/uuid"
)

type ErrorResponse struct {
	Error string `json:"error" example:"internal server error"`
}

func newInternalError() *ErrorResponse {
	return &ErrorResponse{Error: "internal error"}
}

func newBadRequest(err error) *ErrorResponse {
	return &ErrorResponse{Error: fmt.Sprintf("bad request: %v", err)}
}

func newNotFound(id uuid.UUID) *ErrorResponse {
	return &ErrorResponse{Error: fmt.Sprintf("not found: %v", id)}
}

func newDomainConflict(err error) *ErrorResponse {
	return &ErrorResponse{Error: err.Error()}
}
