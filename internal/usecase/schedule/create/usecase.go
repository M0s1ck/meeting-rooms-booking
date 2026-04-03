package create

import (
	"context"
	"time"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/schedule"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
	"github.com/internships-backend/test-backend-M0s1ck/internal/service/authjwt"
)

type Usecase struct {
	repo scheduleRepo
}

func NewUsecase(repo scheduleRepo) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) Execute(ctx context.Context, req *Request, identity *authjwt.Identity) (*Response, error) {
	if err := u.authorize(identity); err != nil {
		return nil, err
	}

	entity, err := schedule.New(
		req.RoomID,
		req.DaysOfWeek,
		req.StartTime,
		req.EndTime,
		time.Now(),
	)
	if err != nil {
		return nil, err
	}

	if err = u.repo.Create(ctx, entity); err != nil {
		return nil, err
	}

	return &Response{
		ID:         entity.ID,
		RoomID:     entity.RoomID,
		DaysOfWeek: entity.DaysOfWeek,
		StartTime:  entity.StartTime,
		EndTime:    entity.EndTime,
	}, nil
}

func (u *Usecase) authorize(identity *authjwt.Identity) error {
	if identity == nil || identity.Role != user.RoleAdmin {
		return user.ErrAdminRoleRequired
	}

	return nil
}
