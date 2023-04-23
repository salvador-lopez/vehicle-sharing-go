package rest

import "errors"

var (
	ErrInternal = errors.New("internal error")
	ErrNotFound = errors.New("resource not found")
)
