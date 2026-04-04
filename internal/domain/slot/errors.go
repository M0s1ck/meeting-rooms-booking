package slot

import "errors"

var (
	ErrSlotStartTooLate = errors.New("slot start time is too late")
	ErrNotFound         = errors.New("slot not found")
)
