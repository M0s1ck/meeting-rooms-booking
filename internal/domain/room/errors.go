package room

import "errors"

var (
	ErrEmptyName       = errors.New("room name is required")
	ErrInvalidCapacity = errors.New("room capacity must be greater than zero")
)
