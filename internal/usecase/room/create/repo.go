package create

import (
	"context"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/room"
)

//go:generate mockgen -source=repo.go -destination=mocks/mock_room_repo.go -package=mocks
type roomRepo interface {
	Create(ctx context.Context, room *room.Room) error
}
