package list

import (
	"context"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/room"
)

//go:generate mockgen -source=repo.go -destination=mocks/mock_room_repo.go -package=mocks
type roomRepo interface {
	List(ctx context.Context) ([]*room.Room, error)
}
