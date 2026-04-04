package list

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/slot"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/slot/ensureroomdate"
)

type roomRepo interface {
	Exists(ctx context.Context, roomID uuid.UUID) (bool, error)
}

type slotRepo interface {
	ListAvailableByRoomAndDate(ctx context.Context, roomID uuid.UUID, date time.Time) ([]slot.Slot, error)
}

type ensureDateFilledUsecase interface {
	Execute(ctx context.Context, req *ensureroomdate.Request) error
}
