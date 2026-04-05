package ensureroomdate_test

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/schedule"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/slot"
	commonmocks "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/common/mocks"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/slot/ensureroomdate"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/slot/ensureroomdate/mocks"
)

func Test_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	slotGen := slot.NewGenerator()
	scheduleRepo := mocks.NewMockscheduleRepo(ctrl)
	slotRepo := mocks.NewMockslotRepo(ctrl)
	txManager := commonmocks.NewTxMngrStubJustCall()

	uc := ensureroomdate.NewUsecase(slotGen, scheduleRepo, slotRepo, txManager)

	now := time.Now().UTC()

	req := &ensureroomdate.Request{
		RoomID: uuid.New(),
		Date:   time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
	}

	start, err := schedule.ParseTimeOfDay("10:00")
	require.NoError(t, err)
	end, err := schedule.ParseTimeOfDay("12:30")
	require.NoError(t, err)
	schedCreatedAt := now.Add(-1 * time.Hour * 24 * 21)

	sched, err := schedule.New(req.RoomID, []int{1, 2, 5}, start, end, schedCreatedAt)
	require.NoError(t, err)

	scheduleRepo.EXPECT().
		GetByRoomID(gomock.Any(), req.RoomID).
		Return(sched, nil)

	slotRepo.EXPECT().
		Add(gomock.Any(), gomock.Any()).
		Return(nil)

	err = uc.Execute(t.Context(), req)
	require.NoError(t, err)
}

func Test_Execute_ScheduleNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	slotGen := slot.NewGenerator()
	scheduleRepo := mocks.NewMockscheduleRepo(ctrl)
	slotRepo := mocks.NewMockslotRepo(ctrl)
	txManager := commonmocks.NewTxMngrStubJustCall()

	uc := ensureroomdate.NewUsecase(slotGen, scheduleRepo, slotRepo, txManager)

	now := time.Now().UTC()

	req := &ensureroomdate.Request{
		RoomID: uuid.New(),
		Date:   time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
	}

	scheduleRepo.EXPECT().
		GetByRoomID(gomock.Any(), req.RoomID).
		Return(nil, schedule.ErrNotFound)

	err := uc.Execute(t.Context(), req)
	require.NoError(t, err)
}

func Test_Execute_RepoErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	slotGen := slot.NewGenerator()
	scheduleRepo := mocks.NewMockscheduleRepo(ctrl)
	slotRepo := mocks.NewMockslotRepo(ctrl)
	txManager := commonmocks.NewTxMngrStubJustCall()

	uc := ensureroomdate.NewUsecase(slotGen, scheduleRepo, slotRepo, txManager)

	now := time.Now().UTC()

	req := &ensureroomdate.Request{
		RoomID: uuid.New(),
		Date:   time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
	}

	scheduleRepo.EXPECT().
		GetByRoomID(gomock.Any(), req.RoomID).
		Return(nil, errors.New("db error"))

	err := uc.Execute(t.Context(), req)
	require.Error(t, err)
}
