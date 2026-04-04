package register

import (
	"context"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
)

type userRepo interface {
	Create(ctx context.Context, user *user.User) error
}
