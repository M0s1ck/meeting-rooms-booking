package create

import (
	"context"

	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/booking"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/slot"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
)

type bookingRepo interface {
	Create(ctx context.Context, booking *booking.Booking) error
}

type slotRepo interface {
	GetByID(ctx context.Context, slotID uuid.UUID) (*slot.Slot, error)
}

type userRepo interface {
	Ensure(ctx context.Context, userID uuid.UUID, role user.Role) error
}
