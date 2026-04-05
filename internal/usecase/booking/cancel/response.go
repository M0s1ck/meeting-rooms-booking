package cancel

import (
	"time"

	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/booking"
)

type Response struct {
	ID             uuid.UUID
	SlotID         uuid.UUID
	UserID         uuid.UUID
	Status         booking.Status
	ConferenceLink *string
	CreatedAt      time.Time
}
