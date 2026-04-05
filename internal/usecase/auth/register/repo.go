package register

import (
	"context"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
)

//go:generate mockgen -source=repo.go -destination=mocks/repo_mocks.go -package=mocks
type userRepo interface {
	Create(ctx context.Context, user *user.User) error
}
