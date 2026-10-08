package listrange

import (
	"time"

	"github.com/google/uuid"
)

type Slot struct {
	ID      uuid.UUID
	RoomID  uuid.UUID
	StartAt time.Time
	EndAt   time.Time
}

type Response struct {
	Slots []Slot
}
