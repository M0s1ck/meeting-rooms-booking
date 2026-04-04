package dummylogin

import (
	"context"

	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
)

type userRepo interface {
	Ensure(ctx context.Context, userID uuid.UUID, role user.Role) error
}
