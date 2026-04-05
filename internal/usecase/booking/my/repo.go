package my

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/booking"
)

//go:generate mockgen -source=repo.go -destination=mocks/repo_mocks.go -package=mocks
type bookingRepo interface {
	ListFutureActiveByUser(ctx context.Context, userID uuid.UUID, now time.Time) ([]booking.Booking, error)
}
