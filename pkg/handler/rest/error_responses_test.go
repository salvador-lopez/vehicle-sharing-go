//go:build unit

package rest_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"vehicle-sharing-go/pkg/handler/rest"
)

func TestErrorResponses(t *testing.T) {
	id := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	cause := errors.New("validation failed")

	tests := []struct {
		name     string
		response *rest.ErrorResponse
		expected string
	}{
		{
			name:     "internal error",
			response: rest.NewInternalError(),
			expected: "internal error",
		},
		{
			name:     "bad request",
			response: rest.NewBadRequest(cause),
			expected: "bad request: validation failed",
		},
		{
			name:     "not found",
			response: rest.NewNotFound(id),
			expected: "not found: 123e4567-e89b-12d3-a456-426614174000",
		},
		{
			name:     "domain conflict",
			response: rest.NewDomainConflict(cause),
			expected: "validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.response.Error)

			encoded, err := json.Marshal(tt.response)
			require.NoError(t, err)
			require.JSONEq(t, `{"error":"`+tt.expected+`"}`, string(encoded))
		})
	}
}
