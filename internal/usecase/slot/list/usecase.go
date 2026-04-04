package list

import (
	"context"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/room"
)

type Usecase struct {
	roomRepo roomRepo
	slotRepo slotRepo
}

func NewUsecase(roomRepo roomRepo, slotRepo slotRepo) *Usecase {
	return &Usecase{
		roomRepo: roomRepo,
		slotRepo: slotRepo,
	}
}

func (u *Usecase) Execute(ctx context.Context, req *Request) (*Response, error) {
	exists, err := u.roomRepo.Exists(ctx, req.RoomID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, room.ErrNotFound
	}

	slots, err := u.slotRepo.ListAvailableByRoomAndDate(ctx, req.RoomID, req.Date)
	if err != nil {
		return nil, err
	}

	resp := &Response{
		Slots: make([]Slot, 0, len(slots)),
	}

	for _, item := range slots {
		resp.Slots = append(resp.Slots, Slot{
			ID:      item.ID,
			RoomID:  item.RoomID,
			StartAt: item.StartAt,
			EndAt:   item.EndAt,
		})
	}

	return resp, nil
}
