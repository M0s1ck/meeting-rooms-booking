package create

import (
	"context"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/room"
)

type roomRepo interface {
	Create(ctx context.Context, room *room.Room) error
}
