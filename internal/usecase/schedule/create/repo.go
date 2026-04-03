package create

import (
	"context"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/schedule"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/slot"
)

type scheduleRepo interface {
	Create(ctx context.Context, schedule *schedule.Schedule) error
}

type slotRepo interface {
	Add(ctx context.Context, slots []slot.Slot) error
}
