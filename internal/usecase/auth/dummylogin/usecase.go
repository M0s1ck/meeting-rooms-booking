package dummylogin

import (
	"context"

	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
)

type Usecase struct {
	tokenGenerator TokenGenerator
}

func NewUsecase(tokenGenerator TokenGenerator) *Usecase {
	return &Usecase{
		tokenGenerator: tokenGenerator,
	}
}

var (
	adminID = uuid.MustParse("48323e70-a002-447c-86c5-d640142aaef0")
	userID  = uuid.MustParse("230e467a-31b8-4a4a-88aa-d4757876b950")
)

func (u *Usecase) Execute(_ context.Context, req *Request) (*Response, error) {
	var token string
	var err error

	switch req.Role {
	case user.RoleAdmin:
		token, err = u.tokenGenerator.Generate(adminID, user.RoleAdmin)
	case user.RoleUser:
		token, err = u.tokenGenerator.Generate(userID, user.RoleUser)
	default:
		return nil, user.ErrInvalidRole
	}

	if err != nil {
		return nil, err
	}

	return &Response{
		Token: token,
	}, nil
}
