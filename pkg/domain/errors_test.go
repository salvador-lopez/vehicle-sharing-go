//go:build unit

package domain_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"vehicle-sharing-go/pkg/domain"
)

func TestWrapErrConflict(t *testing.T) {
	cause := errors.New("unique constraint violated")

	err := domain.WrapErrConflict(cause)

	require.EqualError(t, err, "domain conflict: unique constraint violated")
	require.ErrorIs(t, err, domain.ErrConflict)
	require.ErrorIs(t, err, cause)
}
