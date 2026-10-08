package listrange

import (
	"time"

	"github.com/google/uuid"
)

type Request struct {
	RoomID uuid.UUID
	From   time.Time // first date, inclusive (UTC)
	To     time.Time // last date, inclusive (UTC)
}
