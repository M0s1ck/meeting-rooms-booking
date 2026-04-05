package ensureroomdate

import (
	"context"
	"errors"
	"time"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/schedule"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/slot"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/common"
)

type Usecase struct {
	slotGen      *slot.Generator
	scheduleRepo scheduleRepo
	slotRepo     slotRepo
	txManager    common.TxManager
}

func NewUsecase(
	slotGen *slot.Generator,
	scheduleRepo scheduleRepo,
	slotRepo slotRepo,
	txManager common.TxManager,
) *Usecase {

	return &Usecase{
		slotGen:      slotGen,
		scheduleRepo: scheduleRepo,
		slotRepo:     slotRepo,
		txManager:    txManager,
	}
}

func (u *Usecase) Execute(ctx context.Context, req *Request) error {
	return u.txManager.Do(ctx, func(ctx context.Context) error {
		sched, err := u.scheduleRepo.GetByRoomID(ctx, req.RoomID)
		if errors.Is(err, schedule.ErrNotFound) {
			return nil
		}

		if err != nil {
			return err
		}

		dayStart := time.Date(req.Date.UTC().Year(), req.Date.UTC().Month(), req.Date.UTC().Day(), 0, 0, 0, 0, time.UTC)
		dayEnd := dayStart.AddDate(0, 0, 1)

		slots := u.slotGen.Generate(*sched, dayStart, dayEnd)
		return u.slotRepo.Add(ctx, slots)
	})
}
