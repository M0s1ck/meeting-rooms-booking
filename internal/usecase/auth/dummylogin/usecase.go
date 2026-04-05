package dummylogin

import (
	"context"

	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
)

var (
	adminID = uuid.MustParse("48323e70-a002-447c-86c5-d640142aaef0")
	userID  = uuid.MustParse("230e467a-31b8-4a4a-88aa-d4757876b950")
)

type Usecase struct {
	repo           userRepo
	tokenGenerator tokenGenerator
}

func NewUsecase(tokenGenerator tokenGenerator, repo userRepo) *Usecase {
	return &Usecase{
		tokenGenerator: tokenGenerator,
		repo:           repo,
	}
}

func (u *Usecase) Execute(ctx context.Context, req *Request) (*Response, error) {
	role := req.Role
	var id uuid.UUID

	switch role {
	case user.RoleUser:
		id = userID
	case user.RoleAdmin:
		id = adminID
	default:
		return nil, user.ErrInvalidRole
	}

	// for bookings etc.
	if err := u.repo.Ensure(ctx, id, role); err != nil {
		return nil, err
	}

	token, err := u.tokenGenerator.Generate(id, role)
	if err != nil {
		return nil, err
	}

	return &Response{
		Token: token,
	}, nil
}
