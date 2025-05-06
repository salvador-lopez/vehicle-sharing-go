package rest

import (
	"fmt"
	"github.com/google/uuid"
)

type ErrorResponse struct {
	Error string `json:"error" example:"internal server error"`
}

func NewInternalError() *ErrorResponse {
	return &ErrorResponse{Error: "internal error"}
}

func NewBadRequest(err error) *ErrorResponse {
	return &ErrorResponse{Error: fmt.Sprintf("bad request: %v", err)}
}

func NewNotFound(id uuid.UUID) *ErrorResponse {
	return &ErrorResponse{Error: fmt.Sprintf("not found: %v", id)}
}

func NewDomainConflict(err error) *ErrorResponse {
	return &ErrorResponse{Error: err.Error()}
}
