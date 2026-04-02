package list

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Capacity    *int
	CreatedAt   time.Time
}

type Response struct {
	Rooms []Room
}
