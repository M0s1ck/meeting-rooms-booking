package list

import (
	"context"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/room"
)

type roomRepo interface {
	List(ctx context.Context) ([]*room.Room, error)
}
