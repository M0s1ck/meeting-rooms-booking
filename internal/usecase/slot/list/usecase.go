package list

import (
	"context"
	"time"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/room"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/slot"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/slot/ensureroomdate"
)

type Usecase struct {
	roomRepo         roomRepo
	slotRepo         slotRepo
	ensureDateFilled ensureRoomOnDateFilledUsecase
	slotHorizon      time.Duration
}

func NewUsecase(
	roomRepo roomRepo,
	slotRepo slotRepo,
	ensureDateFilled ensureRoomOnDateFilledUsecase,
	slotHorizon time.Duration,
) *Usecase {

	return &Usecase{
		roomRepo:         roomRepo,
		slotRepo:         slotRepo,
		ensureDateFilled: ensureDateFilled,
		slotHorizon:      slotHorizon,
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

	now := time.Now().UTC()

	if isPastDate(req.Date, now) {
		return &Response{Slots: []Slot{}}, nil
	}

	if isAfterSlotHorizon(req.Date, now, u.slotHorizon) {
		fillReq := &ensureroomdate.Request{
			RoomID: req.RoomID,
			Date:   req.Date,
		}

		if err := u.ensureDateFilled.Execute(ctx, fillReq); err != nil {
			return nil, err
		}
	}

	slots, err := u.slotRepo.ListAvailableByRoomAndDate(ctx, req.RoomID, req.Date, now)
	if err != nil {
		return nil, err
	}

	resp := buildResp(slots)
	return resp, nil
}

func isAfterSlotHorizon(date, now time.Time, slotHorizon time.Duration) bool {
	dayStart := time.Date(date.UTC().Year(), date.UTC().Month(), date.UTC().Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.AddDate(0, 0, 1)

	return dayEnd.After(now.Add(slotHorizon))
}

func isPastDate(date, now time.Time) bool {
	d := date.UTC()
	n := now.UTC()

	dateStart := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
	nowStart := time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.UTC)

	return dateStart.Before(nowStart)
}

func buildResp(slots []slot.Slot) *Response {
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

	return resp
}
