package booking

import "errors"

var (
	ErrSlotIDRequired  = errors.New("slot id is required")
	ErrUserIDRequired  = errors.New("user id is required")
	ErrSlotInPast      = errors.New("cannot create booking for a slot in the past")
	ErrAlreadyBooked   = errors.New("slot is already booked")
	ErrInvalidPage     = errors.New("page must be greater than or equal to 1")
	ErrInvalidPageSize = errors.New("page size must be between 1 and 100")
)
