package ensureroomdate

import (
	"context"

	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/schedule"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/slot"
)

//go:generate mockgen -source=repo.go -destination=mocks/mocks.go -package=mocks
type scheduleRepo interface {
	GetByRoomID(ctx context.Context, roomID uuid.UUID) (*schedule.Schedule, error)
}

type slotRepo interface {
	Add(ctx context.Context, slots []slot.Slot) error
}
