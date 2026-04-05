package create

import (
	"context"

	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/booking"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/slot"
)

//go:generate mockgen -source=repo.go -destination=mocks/repo_mocks.go -package=mocks
type bookingRepo interface {
	Create(ctx context.Context, booking *booking.Booking) error
}

type slotRepo interface {
	GetByID(ctx context.Context, slotID uuid.UUID) (*slot.Slot, error)
}
