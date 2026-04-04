package fillhorizon

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/slot"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/common"
)

const slotHorizon = time.Hour * 24 * 14

type Usecase struct {
	slotGen      *slot.Generator
	slotRepo     slotRepo
	scheduleRepo scheduleRepo
	txManager    common.TxManager
}

func NewUsecase(
	slotGen *slot.Generator,
	slotRepo slotRepo,
	scheduleRepo scheduleRepo,
	txManager common.TxManager,
) *Usecase {

	return &Usecase{
		slotGen:      slotGen,
		slotRepo:     slotRepo,
		scheduleRepo: scheduleRepo,
		txManager:    txManager,
	}
}

func (u *Usecase) FillTail(ctx context.Context) error {
	return u.txManager.Do(ctx, func(ctx context.Context) error {
		now := time.Now().UTC()
		to := now.Add(slotHorizon)

		schedules, err := u.scheduleRepo.List(ctx)
		if err != nil {
			return err
		}

		roomIDs := make([]uuid.UUID, len(schedules))
		for i := range schedules {
			roomIDs[i] = schedules[i].RoomID
		}

		roomEnds, err := u.slotRepo.GetLastEndAtByRooms(ctx, roomIDs)
		if err != nil {
			return err
		}

		slots := make([]slot.Slot, 0, len(roomEnds))

		for _, sch := range schedules {
			from, ok := roomEnds[sch.RoomID]
			if !ok || from.Before(now) {
				from = now
			}

			slots = append(slots, u.slotGen.Generate(sch, from, to)...)
		}

		return u.slotRepo.Add(ctx, slots)
	})
}

func (u *Usecase) FullRepair(ctx context.Context) error {
	return u.txManager.Do(ctx, func(ctx context.Context) error {
		now := time.Now().UTC()
		to := now.Add(slotHorizon)

		schedules, err := u.scheduleRepo.List(ctx)
		if err != nil {
			return err
		}

		slots := make([]slot.Slot, 0)

		for _, sch := range schedules {
			slots = append(slots, u.slotGen.Generate(sch, now, to)...)
		}

		return u.slotRepo.Add(ctx, slots)
	})
}
