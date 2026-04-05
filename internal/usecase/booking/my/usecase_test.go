package my_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/booking"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
	"github.com/internships-backend/test-backend-M0s1ck/internal/service/authjwt"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/booking/my"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/booking/my/mocks"
	"github.com/stretchr/testify/require"
)

func TestUsecase_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockbookingRepo(ctrl)
	uc := my.NewUsecase(repo)

	now := time.Now().UTC()
	userID := uuid.New()

	b1, err := booking.NewActive(uuid.New(), userID, now.Add(-24*time.Hour))
	require.NoError(t, err)

	b2, err := booking.NewActive(uuid.New(), userID, now.Add(-20*time.Hour))
	require.NoError(t, err)

	identity := &authjwt.Identity{
		Role:   user.RoleUser,
		UserID: userID,
	}

	repo.EXPECT().
		ListFutureActiveByUser(
			gomock.Any(),
			userID,
			gomock.Any(),
		).
		Return([]booking.Booking{*b1, *b2}, nil)

	resp, err := uc.Execute(t.Context(), identity)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Bookings, 2)

	require.Equal(t, b1.ID, resp.Bookings[0].ID)
	require.Equal(t, b2.SlotID, resp.Bookings[1].SlotID)
}

func TestUsecase_Execute_NotUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockbookingRepo(ctrl)
	uc := my.NewUsecase(repo)

	identity := &authjwt.Identity{
		Role: user.RoleAdmin,
	}

	resp, err := uc.Execute(t.Context(), identity)

	require.Error(t, err)
	require.ErrorIs(t, err, user.ErrUserRoleRequired)
	require.Nil(t, resp)
}

func TestUsecase_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockbookingRepo(ctrl)
	uc := my.NewUsecase(repo)

	userID := uuid.New()

	identity := &authjwt.Identity{
		Role:   user.RoleUser,
		UserID: userID,
	}

	repo.EXPECT().
		ListFutureActiveByUser(gomock.Any(), userID, gomock.Any()).
		Return(nil, errors.New("db error"))

	resp, err := uc.Execute(context.Background(), identity)

	require.Error(t, err)
	require.Nil(t, resp)
}
