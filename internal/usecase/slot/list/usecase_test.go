package list_test

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
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/slot/list"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/slot/list/mocks"
)

func TestUsecase_Execute_DateInHorizon_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roomRepo := mocks.NewMockroomRepo(ctrl)
	slotRepo := mocks.NewMockslotRepo(ctrl)
	ensureDateFilled := mocks.NewMockensureRoomOnDateFilledUsecase(ctrl)
	slotHorizon := time.Hour * 24 * 14
	uc := list.NewUsecase(roomRepo, slotRepo, ensureDateFilled, slotHorizon)

	now := time.Now().UTC()
	tomorrowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).
		AddDate(0, 0, 1)

	req := &list.Request{
		RoomID: uuid.New(),
		Date:   tomorrowDate,
	}

	slt1, err := slot.New(req.RoomID, tomorrowDate.Add(10*time.Hour), now)
	require.NoError(t, err)

	slt2, err := slot.New(req.RoomID, tomorrowDate.Add(11*time.Hour), now)
	require.NoError(t, err)

	roomRepo.EXPECT().
		Exists(gomock.Any(), req.RoomID).
		Return(true, nil)

	slotRepo.EXPECT().
		ListAvailableByRoomAndDate(gomock.Any(), req.RoomID, req.Date).
		Return([]slot.Slot{*slt1, *slt2}, nil)

	resp, err := uc.Execute(t.Context(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Slots, 2)
	require.Equal(t, resp.Slots[0].ID, slt1.ID)
	require.Equal(t, resp.Slots[1].ID, slt2.ID)
	require.Equal(t, resp.Slots[0].RoomID, req.RoomID)
	require.Equal(t, resp.Slots[0].StartAt, slt1.StartAt)
}

func TestUsecase_Execute_DateAfterHorizon_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roomRepo := mocks.NewMockroomRepo(ctrl)
	slotRepo := mocks.NewMockslotRepo(ctrl)
	ensureDateFilled := mocks.NewMockensureRoomOnDateFilledUsecase(ctrl)
	slotHorizon := time.Hour * 24 * 14
	uc := list.NewUsecase(roomRepo, slotRepo, ensureDateFilled, slotHorizon)

	now := time.Now().UTC()
	dateAfterHorizon := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).
		AddDate(0, 0, 20)

	req := &list.Request{
		RoomID: uuid.New(),
		Date:   dateAfterHorizon,
	}

	roomRepo.EXPECT().
		Exists(gomock.Any(), req.RoomID).
		Return(true, nil)

	ensureDateFilled.EXPECT().
		Execute(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, endureReq *ensureroomdate.Request) error {
			require.Equal(t, req.RoomID, endureReq.RoomID)
			require.Equal(t, req.Date, endureReq.Date)
			return nil
		})

	slt1, err := slot.New(req.RoomID, dateAfterHorizon.Add(10*time.Hour), now)
	require.NoError(t, err)

	slt2, err := slot.New(req.RoomID, dateAfterHorizon.Add(11*time.Hour), now)
	require.NoError(t, err)

	slotRepo.EXPECT().
		ListAvailableByRoomAndDate(gomock.Any(), req.RoomID, req.Date).
		Return([]slot.Slot{*slt1, *slt2}, nil)

	resp, err := uc.Execute(t.Context(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Slots, 2)
	require.Equal(t, resp.Slots[0].ID, slt1.ID)
	require.Equal(t, resp.Slots[1].ID, slt2.ID)
	require.Equal(t, resp.Slots[0].RoomID, req.RoomID)
	require.Equal(t, resp.Slots[0].StartAt, slt1.StartAt)
}

func TestUsecase_Execute_RoomNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roomRepo := mocks.NewMockroomRepo(ctrl)
	slotRepo := mocks.NewMockslotRepo(ctrl)
	ensureDateFilled := mocks.NewMockensureRoomOnDateFilledUsecase(ctrl)
	slotHorizon := time.Hour * 24 * 14
	uc := list.NewUsecase(roomRepo, slotRepo, ensureDateFilled, slotHorizon)

	now := time.Now().UTC()
	dateAfterHorizon := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).
		AddDate(0, 0, 20)

	req := &list.Request{
		RoomID: uuid.New(),
		Date:   dateAfterHorizon,
	}

	roomRepo.EXPECT().
		Exists(gomock.Any(), req.RoomID).
		Return(false, nil)

	resp, err := uc.Execute(t.Context(), req)
	require.Equal(t, room.ErrNotFound, err)
	require.Nil(t, resp)
}

func TestUsecase_Execute_DateAfterHorizon_ScheduleNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roomRepo := mocks.NewMockroomRepo(ctrl)
	slotRepo := mocks.NewMockslotRepo(ctrl)
	ensureDateFilled := mocks.NewMockensureRoomOnDateFilledUsecase(ctrl)
	slotHorizon := time.Hour * 24 * 14
	uc := list.NewUsecase(roomRepo, slotRepo, ensureDateFilled, slotHorizon)

	now := time.Now().UTC()
	dateAfterHorizon := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).
		AddDate(0, 0, 20)

	req := &list.Request{
		RoomID: uuid.New(),
		Date:   dateAfterHorizon,
	}

	roomRepo.EXPECT().
		Exists(gomock.Any(), req.RoomID).
		Return(true, nil)

	ensureDateFilled.EXPECT().
		Execute(gomock.Any(), gomock.Any()).
		Return(errors.New("ensure date filled err"))

	resp, err := uc.Execute(t.Context(), req)
	require.Error(t, err)
	require.Nil(t, resp)
}
