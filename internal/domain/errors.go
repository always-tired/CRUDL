package domain

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrDuplicate       = errors.New("duplicate")
	ErrInvalidArgument = errors.New("invalid argument")
)
