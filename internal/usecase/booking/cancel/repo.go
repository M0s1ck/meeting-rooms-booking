package cancel

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/booking"
)

type bookingRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (*booking.Booking, error)
	Cancel(ctx context.Context, id uuid.UUID, now time.Time) error
}
