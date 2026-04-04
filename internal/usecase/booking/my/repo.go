package my

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/booking"
)

type bookingRepo interface {
	ListFutureByUser(ctx context.Context, userID uuid.UUID, now time.Time) ([]booking.Booking, error)
}
