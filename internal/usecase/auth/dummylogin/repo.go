package dummylogin

import (
	"context"

	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
)

//go:generate mockgen -source=repo.go -destination=mocks/repo_mocks.go -package=mocks
type userRepo interface {
	Ensure(ctx context.Context, userID uuid.UUID, role user.Role) error
}
