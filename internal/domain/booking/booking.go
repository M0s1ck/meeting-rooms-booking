package booking

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusActive    Status = "active"
	StatusCancelled Status = "cancelled"
)

type Booking struct {
	ID             uuid.UUID
	SlotID         uuid.UUID
	UserID         uuid.UUID
	Status         Status
	ConferenceLink *string
	CreatedAt      time.Time
}

func NewActive(slotID, userID uuid.UUID, now time.Time) (*Booking, error) {
	if slotID == uuid.Nil {
		return nil, ErrSlotIDRequired
	}

	if userID == uuid.Nil {
		return nil, ErrUserIDRequired
	}

	return &Booking{
		ID:        uuid.New(),
		SlotID:    slotID,
		UserID:    userID,
		Status:    StatusActive,
		CreatedAt: now.UTC(),
	}, nil
}
