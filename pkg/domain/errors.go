package domain

import (
	"errors"
	"fmt"
)

var ErrConflict = errors.New("domain conflict")

func WrapErrConflict(err error) error {
	return fmt.Errorf("%w: %w", ErrConflict, err)
}
