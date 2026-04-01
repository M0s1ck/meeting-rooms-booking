package create

import (
	"context"
	"time"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/room"
)

type Usecase struct {
	repo roomRepo
}

func NewUsecase(repo roomRepo) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (u *Usecase) Execute(ctx context.Context, req *Request) (*Response, error) {
	entity, err := room.New(
		req.Name,
		req.Description,
		req.Capacity,
		time.Now(),
	)

	if err != nil {
		return nil, err
	}

	if err = u.repo.Create(ctx, entity); err != nil {
		return nil, err
	}

	return &Response{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		Capacity:    entity.Capacity,
		CreatedAt:   entity.CreatedAt,
	}, nil
}
