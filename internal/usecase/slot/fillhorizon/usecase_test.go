package fillhorizon_test

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
	commonmocks "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/common/mocks"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/slot/fillhorizon"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/slot/fillhorizon/mocks"
)

func Test_FillTail_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	slotGen := slot.NewGenerator()
	slotRepo := mocks.NewMockslotRepo(ctrl)
	schedRepo := mocks.NewMockscheduleRepo(ctrl)
	txManager := commonmocks.NewTxMngrStubJustCall()
	slotHorizon := 14 * 24 * time.Hour
	uc := fillhorizon.NewUsecase(slotGen, slotRepo, schedRepo, txManager, slotHorizon)

	start, err := schedule.ParseTimeOfDay("0:30")
	require.NoError(t, err)
	end, err := schedule.ParseTimeOfDay("23:30")
	require.NoError(t, err)
	now := time.Now().UTC()

	roomID1 := uuid.New()
	roomID2 := uuid.New()
	allWeekDays := []int{1, 2, 3, 4, 5, 6, 7}

	schedule1, err := schedule.New(roomID1, allWeekDays, start, end, now.Add(-time.Hour*24*20))
	require.NoError(t, err)

	schedule2, err := schedule.New(roomID2, allWeekDays, start, end, now.Add(-time.Hour*24*25))
	require.NoError(t, err)

	schedRepo.EXPECT().
		List(gomock.Any()).
		Return([]schedule.Schedule{*schedule1, *schedule2}, nil)

	slotRepo.EXPECT().
		GetLastEndAtByRooms(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, roomIDs []uuid.UUID) (map[uuid.UUID]time.Time, error) {
			require.Contains(t, roomIDs, roomID1)
			require.Contains(t, roomIDs, roomID2)
			return map[uuid.UUID]time.Time{
				roomID1: now.Add(slotHorizon - 5*time.Hour),
				roomID2: now.Add(-time.Minute * 30),
			}, nil
		})

	slotRepo.EXPECT().
		Add(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, slots []slot.Slot) error {
			require.NotEmpty(t, slots)

			roomIDs := make(map[uuid.UUID]struct{})
			for _, s := range slots {
				roomIDs[s.RoomID] = struct{}{}
			}

			require.Contains(t, roomIDs, roomID1)
			require.Contains(t, roomIDs, roomID2)

			for _, s := range slots {
				require.True(t, s.StartAt.Before(now.Add(slotHorizon)))
				require.True(t, s.StartAt.After(now))
			}

			return nil
		})

	err = uc.FillTail(t.Context())
	require.NoError(t, err)
}

func Test_Repair_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	slotGen := slot.NewGenerator()
	slotRepo := mocks.NewMockslotRepo(ctrl)
	schedRepo := mocks.NewMockscheduleRepo(ctrl)
	txManager := commonmocks.NewTxMngrStubJustCall()
	slotHorizon := 14 * 24 * time.Hour
	uc := fillhorizon.NewUsecase(slotGen, slotRepo, schedRepo, txManager, slotHorizon)

	start, err := schedule.ParseTimeOfDay("0:30")
	require.NoError(t, err)
	end, err := schedule.ParseTimeOfDay("23:30")
	require.NoError(t, err)
	now := time.Now().UTC()

	roomID1 := uuid.New()
	roomID2 := uuid.New()
	allWeekDays := []int{1, 2, 3, 4, 5, 6, 7}

	schedule1, err := schedule.New(roomID1, allWeekDays, start, end, now.Add(-time.Hour*24*20))
	require.NoError(t, err)

	schedule2, err := schedule.New(roomID2, allWeekDays, start, end, now.Add(-time.Hour*24*25))
	require.NoError(t, err)

	schedRepo.EXPECT().
		List(gomock.Any()).
		Return([]schedule.Schedule{*schedule1, *schedule2}, nil)

	slotRepo.EXPECT().
		Add(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, slots []slot.Slot) error {
			require.NotEmpty(t, slots)

			roomIDs := make(map[uuid.UUID]struct{})
			for _, s := range slots {
				roomIDs[s.RoomID] = struct{}{}
			}

			require.Contains(t, roomIDs, roomID1)
			require.Contains(t, roomIDs, roomID2)

			for _, s := range slots {
				require.True(t, s.StartAt.After(now))
				require.True(t, s.StartAt.Before(now.Add(slotHorizon)))
			}

			return nil
		})

	err = uc.FullRepair(t.Context())
	require.NoError(t, err)
}

func Test_Repair_RepoErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	slotGen := slot.NewGenerator()
	slotRepo := mocks.NewMockslotRepo(ctrl)
	schedRepo := mocks.NewMockscheduleRepo(ctrl)
	txManager := commonmocks.NewTxMngrStubJustCall()
	slotHorizon := 14 * 24 * time.Hour
	uc := fillhorizon.NewUsecase(slotGen, slotRepo, schedRepo, txManager, slotHorizon)

	schedRepo.EXPECT().
		List(gomock.Any()).
		Return(nil, errors.New("db error"))

	err := uc.FullRepair(t.Context())
	require.Error(t, err)
}
