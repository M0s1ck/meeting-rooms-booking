package booking

import "errors"

var (
	ErrSlotIDRequired = errors.New("slot id is required")
	ErrUserIDRequired = errors.New("user id is required")
	ErrSlotInPast     = errors.New("cannot create booking for a slot in the past")
	ErrAlreadyBooked  = errors.New("slot is already booked")
)
