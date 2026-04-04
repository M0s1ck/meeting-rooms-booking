package fillhorizon

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/schedule"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/slot"
)

type scheduleRepo interface {
	List(ctx context.Context) ([]schedule.Schedule, error)
}

type slotRepo interface {
	GetLastEndAtByRooms(ctx context.Context, roomIDS []uuid.UUID) (map[uuid.UUID]time.Time, error)
	Add(ctx context.Context, slots []slot.Slot) error
}
