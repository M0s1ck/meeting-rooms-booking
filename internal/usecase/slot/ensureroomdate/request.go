package ensureroomdate

import (
	"time"

	"github.com/google/uuid"
)

type Request struct {
	RoomID uuid.UUID
	Date   time.Time
}
