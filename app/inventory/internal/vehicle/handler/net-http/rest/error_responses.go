package rest

import (
	"fmt"
	"github.com/google/uuid"
)

type errorResponse struct {
	Error string `json:"error"`
}

func newInternalError() *errorResponse {
	return &errorResponse{Error: "internal error"}
}

func newBadRequest(err error) *errorResponse {
	return &errorResponse{Error: fmt.Sprintf("bad request: %v", err)}
}

func newNotFound(id uuid.UUID) *errorResponse {
	return &errorResponse{Error: fmt.Sprintf("not found: %v", id)}
}

func newDomainConflict(err error) *errorResponse {
	return &errorResponse{Error: err.Error()}
}