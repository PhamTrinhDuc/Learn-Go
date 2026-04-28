package utils

import "errors"

var (
	ErrNilInput          = errors.New("input cannot be nil")
	ErrInvalidID         = errors.New("invalid id")
	ErrNotFound          = errors.New("resource not found")
	ErrDuplicateResource = errors.New("resource already exists")
)
