package list

import (
	"context"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/booking"
)

type bookingRepo interface {
	List(ctx context.Context, page, pageSize int) ([]booking.Booking, int, error)
}
