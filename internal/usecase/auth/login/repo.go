package login

import (
	"context"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
)

type userRepo interface {
	GetByEmail(ctx context.Context, email user.Email) (*user.User, error)
}
