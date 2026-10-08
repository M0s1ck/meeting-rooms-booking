package listrange

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/slot"
)

//go:generate mockgen -source=repo.go -destination=mocks/repo_mocks.go -package=mocks

type roomRepo interface {
	Exists(ctx context.Context, roomID uuid.UUID) (bool, error)
}

type slotRepo interface {
	// ListAvailableByRoomAndRange returns free slots of the room with start in [from, to), not earlier than now.
	ListAvailableByRoomAndRange(ctx context.Context, roomID uuid.UUID, from, to, now time.Time) ([]slot.Slot, error)
}
