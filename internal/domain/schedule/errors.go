package schedule

import "errors"

var (
	ErrRoomIDRequired   = errors.New("room id is required")
	ErrEmptyDaysOfWeek  = errors.New("days of week are required")
	ErrInvalidDayOfWeek = errors.New("days of week must contain unique values from 1 to 7")
	ErrInvalidTimeRange = errors.New("start time must be before end time")
	ErrInvalidTimeFmt   = errors.New("time must be in HH:MM format")
	ErrAlreadyExists    = errors.New("schedule for this room already exists")
	ErrNotFound         = errors.New("schedule not found")
)
