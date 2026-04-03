package create

import (
	"context"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/schedule"
)

type scheduleRepo interface {
	Create(ctx context.Context, schedule *schedule.Schedule) error
}
