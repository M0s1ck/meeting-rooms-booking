package list

import (
	"time"

	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/booking"
)

type Booking struct {
	ID             uuid.UUID
	SlotID         uuid.UUID
	UserID         uuid.UUID
	Status         booking.Status
	ConferenceLink *string
	CreatedAt      time.Time
}

type Pagination struct {
	Page     int
	PageSize int
	Total    int
}

type Response struct {
	Bookings   []Booking
	Pagination Pagination
}
