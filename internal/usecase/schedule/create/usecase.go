package create

import (
	"context"
	"time"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/schedule"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/slot"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
	"github.com/internships-backend/test-backend-M0s1ck/internal/service/authjwt"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/common"
)

type Usecase struct {
	slotGen     *slot.Generator
	schedRepo   scheduleRepo
	slotRepo    slotRepo
	txManager   common.TxManager
	slotHorizon time.Duration
}

func NewUsecase(
	slotGen *slot.Generator,
	schedRepo scheduleRepo,
	slotRepo slotRepo,
	txManager common.TxManager,
	slotHorizon time.Duration,
) *Usecase {

	return &Usecase{
		slotGen:     slotGen,
		schedRepo:   schedRepo,
		slotRepo:    slotRepo,
		txManager:   txManager,
		slotHorizon: slotHorizon,
	}
}

func (u *Usecase) Execute(ctx context.Context, req *Request, identity *authjwt.Identity) (*Response, error) {
	now := time.Now().UTC()

	if err := u.authorize(identity); err != nil {
		return nil, err
	}

	sched, err := schedule.New(
		req.RoomID,
		req.DaysOfWeek,
		req.StartTime,
		req.EndTime,
		now,
	)
	if err != nil {
		return nil, err
	}

	err = u.txManager.Do(ctx, func(ctx context.Context) error {
		if err = u.schedRepo.Create(ctx, sched); err != nil {
			return err
		}

		slots := u.slotGen.Generate(*sched, now, now.Add(u.slotHorizon))
		return u.slotRepo.Add(ctx, slots)
	})

	if err != nil {
		return nil, err
	}

	return &Response{
		ID:         sched.ID,
		RoomID:     sched.RoomID,
		DaysOfWeek: sched.DaysOfWeek,
		StartTime:  sched.StartTime,
		EndTime:    sched.EndTime,
	}, nil
}

func (u *Usecase) authorize(identity *authjwt.Identity) error {
	if identity == nil || identity.Role != user.RoleAdmin {
		return user.ErrAdminRoleRequired
	}

	return nil
}
