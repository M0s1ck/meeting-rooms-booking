package register

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
)

type Usecase struct {
	repo       userRepo
	passHasher passHasher
}

func NewUsecase(repo userRepo, passHasher passHasher) *Usecase {
	return &Usecase{
		repo:       repo,
		passHasher: passHasher,
	}
}

func (u *Usecase) Execute(ctx context.Context, req *Request) (*Response, error) {
	usr, err := u.buildUser(req)
	if err != nil {
		return nil, err
	}

	if err := u.repo.Create(ctx, usr); err != nil {
		return nil, err
	}

	return &Response{
		ID:        usr.ID,
		Email:     usr.Email,
		Role:      usr.Role,
		CreatedAt: usr.CreatedAt,
	}, nil
}

func (u *Usecase) buildUser(req *Request) (*user.User, error) {
	now := time.Now().UTC()

	email, err := user.NewEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if err := user.ValidatePassStrength(req.Password); err != nil {
		return nil, err
	}

	passHash, err := u.passHasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	role, err := user.ParseRole(req.Role)
	if err != nil {
		return nil, err
	}

	id := uuid.New()

	usr, err := user.New(id, email, passHash, role, now)
	if err != nil {
		return nil, err
	}

	return usr, nil
}
