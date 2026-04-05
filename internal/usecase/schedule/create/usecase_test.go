package create_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/schedule"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/slot"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
	"github.com/internships-backend/test-backend-M0s1ck/internal/service/authjwt"
	commonmocks "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/common/mocks"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/schedule/create"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/schedule/create/mocks"
)

func Test_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	slotGen := slot.NewGenerator()
	slotRepo := mocks.NewMockslotRepo(ctrl)
	schedRepo := mocks.NewMockscheduleRepo(ctrl)
	txManager := commonmocks.NewTxMngrStubJustCall()
	uc := create.NewUsecase(slotGen, schedRepo, slotRepo, txManager, 14*24*time.Hour)

	start, err := schedule.ParseTimeOfDay("10:00")
	require.NoError(t, err)
	end, err := schedule.ParseTimeOfDay("12:30")
	require.NoError(t, err)

	req := &create.Request{
		RoomID:     uuid.New(),
		DaysOfWeek: []int{1, 2, 5},
		StartTime:  start,
		EndTime:    end,
	}

	identity := &authjwt.Identity{
		Role: user.RoleAdmin,
	}

	schedRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, s *schedule.Schedule) error {
			require.Equal(t, s.RoomID, req.RoomID)
			require.Equal(t, s.DaysOfWeek, req.DaysOfWeek)
			return nil
		})

	slotRepo.EXPECT().
		Add(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, slots []slot.Slot) error {
			require.Equal(t, slots[0].RoomID, req.RoomID)
			return nil
		})

	resp, err := uc.Execute(t.Context(), req, identity)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, req.RoomID, resp.RoomID)
}

func Test_Execute_NotAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	slotGen := slot.NewGenerator()
	slotRepo := mocks.NewMockslotRepo(ctrl)
	schedRepo := mocks.NewMockscheduleRepo(ctrl)
	txManager := commonmocks.NewTxMngrStubJustCall()
	uc := create.NewUsecase(slotGen, schedRepo, slotRepo, txManager, 14*24*time.Hour)

	start, err := schedule.ParseTimeOfDay("10:00")
	require.NoError(t, err)
	end, err := schedule.ParseTimeOfDay("12:30")
	require.NoError(t, err)

	req := &create.Request{
		RoomID:     uuid.New(),
		DaysOfWeek: []int{1, 2, 5},
		StartTime:  start,
		EndTime:    end,
	}

	identity := &authjwt.Identity{
		Role: user.RoleUser,
	}

	resp, err := uc.Execute(t.Context(), req, identity)
	require.Equal(t, user.ErrAdminRoleRequired, err)
	require.Nil(t, resp)
}

func Test_Execute_InvalidSchedule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	slotGen := slot.NewGenerator()
	slotRepo := mocks.NewMockslotRepo(ctrl)
	schedRepo := mocks.NewMockscheduleRepo(ctrl)
	txManager := commonmocks.NewTxMngrStubJustCall()
	uc := create.NewUsecase(slotGen, schedRepo, slotRepo, txManager, 14*24*time.Hour)

	start, err := schedule.ParseTimeOfDay("12:00")
	require.NoError(t, err)
	end, err := schedule.ParseTimeOfDay("9:30")
	require.NoError(t, err)

	req := &create.Request{
		RoomID:     uuid.New(),
		DaysOfWeek: []int{1, 2, 5},
		StartTime:  start,
		EndTime:    end,
	}

	identity := &authjwt.Identity{
		Role: user.RoleAdmin,
	}

	resp, err := uc.Execute(t.Context(), req, identity)
	require.Equal(t, schedule.ErrInvalidTimeRange, err)
	require.Nil(t, resp)
}

func Test_Execute_ScheduleRepoAlreadyExistsErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	slotGen := slot.NewGenerator()
	slotRepo := mocks.NewMockslotRepo(ctrl)
	schedRepo := mocks.NewMockscheduleRepo(ctrl)
	txManager := commonmocks.NewTxMngrStubJustCall()
	uc := create.NewUsecase(slotGen, schedRepo, slotRepo, txManager, 14*24*time.Hour)

	start, err := schedule.ParseTimeOfDay("9:00")
	require.NoError(t, err)
	end, err := schedule.ParseTimeOfDay("12:30")
	require.NoError(t, err)

	req := &create.Request{
		RoomID:     uuid.New(),
		DaysOfWeek: []int{1, 2, 5},
		StartTime:  start,
		EndTime:    end,
	}

	identity := &authjwt.Identity{
		Role: user.RoleAdmin,
	}

	schedRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(schedule.ErrAlreadyExists)

	resp, err := uc.Execute(t.Context(), req, identity)
	require.Equal(t, schedule.ErrAlreadyExists, err)
	require.Nil(t, resp)
}

func Test_Execute_SlotRepoErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	slotGen := slot.NewGenerator()
	slotRepo := mocks.NewMockslotRepo(ctrl)
	schedRepo := mocks.NewMockscheduleRepo(ctrl)
	txManager := commonmocks.NewTxMngrStubJustCall()
	uc := create.NewUsecase(slotGen, schedRepo, slotRepo, txManager, 14*24*time.Hour)

	start, err := schedule.ParseTimeOfDay("9:00")
	require.NoError(t, err)
	end, err := schedule.ParseTimeOfDay("12:30")
	require.NoError(t, err)

	req := &create.Request{
		RoomID:     uuid.New(),
		DaysOfWeek: []int{1, 2, 5},
		StartTime:  start,
		EndTime:    end,
	}

	identity := &authjwt.Identity{
		Role: user.RoleAdmin,
	}

	schedRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	slotRepo.EXPECT().
		Add(gomock.Any(), gomock.Any()).
		Return(errors.New("db error"))

	resp, err := uc.Execute(t.Context(), req, identity)
	require.Error(t, err)
	require.Nil(t, resp)
}
