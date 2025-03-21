package domain

import (
	"errors"
	"fmt"
)

var ErrConflict = errors.New("domain conflict")

func WrapConflict(err error) error {
	return fmt.Errorf("%w: %w", ErrConflict, err)
}
