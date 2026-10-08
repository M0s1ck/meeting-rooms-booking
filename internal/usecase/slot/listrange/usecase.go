package listrange

import (
	"context"
	"time"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/room"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/slot"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/slot/ensureroomdate"
)

// Usecase lists free slots of a room for a range of dates [From, To] (inclusive, UTC).
// Same approach as the single-date listing: slots within the horizon are already in the DB,
// dates beyond the horizon are filled on demand (idempotently) before reading.
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
	from := startOfDay(req.From)
	to := startOfDay(req.To)

	if from.After(to) {
		return nil, ErrFromAfterTo
	}
	if daysInclusive(from, to) > MaxRangeDays {
		return nil, ErrRangeTooLarge
	}

	exists, err := u.roomRepo.Exists(ctx, req.RoomID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, room.ErrNotFound
	}

	now := time.Now().UTC()
	today := startOfDay(now)

	// past days have no bookable slots
	if from.Before(today) {
		from = today
	}
	if from.After(to) {
		return &Response{Slots: []Slot{}}, nil
	}

	// generate slots for days that are beyond the pre-filled horizon
	horizonEnd := now.Add(u.slotHorizon)
	for day := from; !day.After(to); day = day.AddDate(0, 0, 1) {
		if !day.AddDate(0, 0, 1).After(horizonEnd) {
			continue
		}

		fillReq := &ensureroomdate.Request{RoomID: req.RoomID, Date: day}
		if err := u.ensureDateFilled.Execute(ctx, fillReq); err != nil {
			return nil, err
		}
	}

	slots, err := u.slotRepo.ListAvailableByRoomAndRange(ctx, req.RoomID, from, to.AddDate(0, 0, 1), now)
	if err != nil {
		return nil, err
	}

	return buildResp(slots), nil
}

func startOfDay(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func daysInclusive(from, to time.Time) int {
	return int(to.Sub(from).Hours()/24) + 1
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
