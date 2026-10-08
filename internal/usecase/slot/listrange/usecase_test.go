package listrange_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/room"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/slot"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/slot/ensureroomdate"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/slot/listrange"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/slot/listrange/mocks"
)

const horizon = 14 * 24 * time.Hour

type fixture struct {
	roomRepo *mocks.MockroomRepo
	slotRepo *mocks.MockslotRepo
	ensure   *mocks.MockensureRoomOnDateFilledUsecase
	uc       *listrange.Usecase
}

func newFixture(t *testing.T) *fixture {
	ctrl := gomock.NewController(t)
	f := &fixture{
		roomRepo: mocks.NewMockroomRepo(ctrl),
		slotRepo: mocks.NewMockslotRepo(ctrl),
		ensure:   mocks.NewMockensureRoomOnDateFilledUsecase(ctrl),
	}
	f.uc = listrange.NewUsecase(f.roomRepo, f.slotRepo, f.ensure, horizon)
	return f
}

func today() time.Time {
	n := time.Now().UTC()
	return time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.UTC)
}

func TestExecute_WithinHorizon_ReadsWholeRange(t *testing.T) {
	f := newFixture(t)
	from := today().AddDate(0, 0, 1)
	to := from.AddDate(0, 0, 6)
	req := &listrange.Request{RoomID: uuid.New(), From: from, To: to}

	s1, err := slot.New(req.RoomID, from.Add(9*time.Hour), time.Now())
	require.NoError(t, err)
	s2, err := slot.New(req.RoomID, to.Add(17*time.Hour), time.Now())
	require.NoError(t, err)

	f.roomRepo.EXPECT().Exists(gomock.Any(), req.RoomID).Return(true, nil)
	f.ensure.EXPECT().Execute(gomock.Any(), gomock.Any()).Times(0)
	// 'to' is inclusive, so the repo gets the exclusive upper bound to+1 day
	f.slotRepo.EXPECT().
		ListAvailableByRoomAndRange(gomock.Any(), req.RoomID, from, to.AddDate(0, 0, 1), gomock.Any()).
		Return([]slot.Slot{*s1, *s2}, nil)

	resp, err := f.uc.Execute(t.Context(), req)
	require.NoError(t, err)
	require.Len(t, resp.Slots, 2)
	require.Equal(t, s1.ID, resp.Slots[0].ID)
	require.Equal(t, s2.StartAt, resp.Slots[1].StartAt)
}

func TestExecute_PartlyBeyondHorizon_FillsOnlyThoseDays(t *testing.T) {
	f := newFixture(t)
	from := today().AddDate(0, 0, 10)
	to := today().AddDate(0, 0, 20)
	req := &listrange.Request{RoomID: uuid.New(), From: from, To: to}

	f.roomRepo.EXPECT().Exists(gomock.Any(), req.RoomID).Return(true, nil)

	var filled []time.Time
	f.ensure.EXPECT().Execute(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, r *ensureroomdate.Request) error {
			require.Equal(t, req.RoomID, r.RoomID)
			filled = append(filled, r.Date)
			return nil
		}).AnyTimes()
	f.slotRepo.EXPECT().ListAvailableByRoomAndRange(gomock.Any(), req.RoomID, from, to.AddDate(0, 0, 1), gomock.Any()).
		Return([]slot.Slot{}, nil)

	_, err := f.uc.Execute(t.Context(), req)
	require.NoError(t, err)

	// days fully inside [now, now+horizon) are pre-filled; the rest up to 'to' is filled on demand
	require.NotEmpty(t, filled)
	require.Equal(t, to, filled[len(filled)-1])
	for _, d := range filled {
		require.True(t, d.AddDate(0, 0, 1).After(time.Now().Add(horizon)), "day %s is inside horizon", d)
	}
}

func TestExecute_PastFromIsClampedToToday(t *testing.T) {
	f := newFixture(t)
	req := &listrange.Request{RoomID: uuid.New(), From: today().AddDate(0, 0, -5), To: today().AddDate(0, 0, 2)}

	f.roomRepo.EXPECT().Exists(gomock.Any(), req.RoomID).Return(true, nil)
	f.slotRepo.EXPECT().
		ListAvailableByRoomAndRange(gomock.Any(), req.RoomID, today(), today().AddDate(0, 0, 3), gomock.Any()).
		Return([]slot.Slot{}, nil)

	resp, err := f.uc.Execute(t.Context(), req)
	require.NoError(t, err)
	require.Empty(t, resp.Slots)
}

func TestExecute_RangeFullyInPast_ReturnsEmpty(t *testing.T) {
	f := newFixture(t)
	req := &listrange.Request{RoomID: uuid.New(), From: today().AddDate(0, 0, -7), To: today().AddDate(0, 0, -1)}

	f.roomRepo.EXPECT().Exists(gomock.Any(), req.RoomID).Return(true, nil)
	f.slotRepo.EXPECT().ListAvailableByRoomAndRange(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	resp, err := f.uc.Execute(t.Context(), req)
	require.NoError(t, err)
	require.NotNil(t, resp.Slots)
	require.Empty(t, resp.Slots)
}

func TestExecute_SingleDay(t *testing.T) {
	f := newFixture(t)
	day := today().AddDate(0, 0, 3)
	req := &listrange.Request{RoomID: uuid.New(), From: day, To: day}

	f.roomRepo.EXPECT().Exists(gomock.Any(), req.RoomID).Return(true, nil)
	f.slotRepo.EXPECT().
		ListAvailableByRoomAndRange(gomock.Any(), req.RoomID, day, day.AddDate(0, 0, 1), gomock.Any()).
		Return([]slot.Slot{}, nil)

	_, err := f.uc.Execute(t.Context(), req)
	require.NoError(t, err)
}

func TestExecute_FromAfterTo(t *testing.T) {
	f := newFixture(t)
	req := &listrange.Request{RoomID: uuid.New(), From: today().AddDate(0, 0, 5), To: today().AddDate(0, 0, 4)}

	resp, err := f.uc.Execute(t.Context(), req)
	require.ErrorIs(t, err, listrange.ErrFromAfterTo)
	require.Nil(t, resp)
}

func TestExecute_RangeLimit(t *testing.T) {
	from := today().AddDate(0, 0, 1)

	t.Run("max allowed", func(t *testing.T) {
		f := newFixture(t)
		req := &listrange.Request{RoomID: uuid.New(), From: from, To: from.AddDate(0, 0, listrange.MaxRangeDays-1)}
		f.roomRepo.EXPECT().Exists(gomock.Any(), req.RoomID).Return(true, nil)
		f.ensure.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		f.slotRepo.EXPECT().ListAvailableByRoomAndRange(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return([]slot.Slot{}, nil)

		_, err := f.uc.Execute(t.Context(), req)
		require.NoError(t, err)
	})

	t.Run("one day more", func(t *testing.T) {
		f := newFixture(t)
		req := &listrange.Request{RoomID: uuid.New(), From: from, To: from.AddDate(0, 0, listrange.MaxRangeDays)}

		_, err := f.uc.Execute(t.Context(), req)
		require.ErrorIs(t, err, listrange.ErrRangeTooLarge)
	})
}

func TestExecute_RoomNotFound(t *testing.T) {
	f := newFixture(t)
	req := &listrange.Request{RoomID: uuid.New(), From: today(), To: today()}

	f.roomRepo.EXPECT().Exists(gomock.Any(), req.RoomID).Return(false, nil)

	_, err := f.uc.Execute(t.Context(), req)
	require.ErrorIs(t, err, room.ErrNotFound)
}

func TestExecute_FillError(t *testing.T) {
	f := newFixture(t)
	day := today().AddDate(0, 0, 25)
	req := &listrange.Request{RoomID: uuid.New(), From: day, To: day}

	f.roomRepo.EXPECT().Exists(gomock.Any(), req.RoomID).Return(true, nil)
	f.ensure.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(errors.New("db down"))

	_, err := f.uc.Execute(t.Context(), req)
	require.Error(t, err)
}
