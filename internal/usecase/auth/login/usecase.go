package login

import (
	"context"
	"errors"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
)

type Usecase struct {
	repo         userRepo
	tokenGen     tokenGenerator
	passComparor passComparor
}

func NewUsecase(repo userRepo, tokenGen tokenGenerator, passComparor passComparor) *Usecase {
	return &Usecase{
		repo:         repo,
		tokenGen:     tokenGen,
		passComparor: passComparor,
	}
}

func (u *Usecase) Execute(ctx context.Context, req *Request) (*Response, error) {
	email, err := user.NewEmail(req.Email)
	if err != nil {
		return nil, user.ErrInvalidCredentials
	}

	usr, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return nil, user.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := u.passComparor.Compare(usr.PassHash, req.Password); err != nil {
		if errors.Is(err, user.ErrWrongPassword) {
			return nil, user.ErrInvalidCredentials
		}
		return nil, err
	}

	token, err := u.tokenGen.Generate(usr.ID, usr.Role)
	if err != nil {
		return nil, err
	}

	return &Response{
		Token: token,
	}, nil
}
