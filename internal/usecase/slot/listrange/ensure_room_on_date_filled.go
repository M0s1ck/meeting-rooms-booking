package listrange

import (
	"context"

	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/slot/ensureroomdate"
)

//go:generate mockgen -source=ensure_room_on_date_filled.go -destination=mocks/ensure_room_on_date_filled_mock.go -package=mocks

type ensureRoomOnDateFilledUsecase interface {
	Execute(ctx context.Context, req *ensureroomdate.Request) error
}
