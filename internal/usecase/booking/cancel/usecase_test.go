package cancel_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/booking"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
	"github.com/internships-backend/test-backend-M0s1ck/internal/service/authjwt"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/booking/cancel"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/booking/cancel/mocks"
	"github.com/stretchr/testify/require"
)

func TestUsecase_Execute_Success_CancelActive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockbookingRepo(ctrl)
	uc := cancel.NewUsecase(repo)

	userID := uuid.New()
	bookingID := uuid.New()

	book := &booking.Booking{
		ID:     bookingID,
		UserID: userID,
		Status: booking.StatusActive,
	}

	cancelled := &booking.Booking{
		ID:     bookingID,
		UserID: userID,
		Status: booking.StatusCancelled,
	}

	identity := &authjwt.Identity{
		Role:   user.RoleUser,
		UserID: userID,
	}

	gomock.InOrder(
		repo.EXPECT().
			GetByID(gomock.Any(), bookingID).
			Return(book, nil),

		repo.EXPECT().
			Cancel(gomock.Any(), bookingID, gomock.Any()).
			Return(nil),

		repo.EXPECT().
			GetByID(gomock.Any(), bookingID).
			Return(cancelled, nil),
	)

	resp, err := uc.Execute(t.Context(), bookingID, identity)

	require.NoError(t, err)
	require.Equal(t, booking.StatusCancelled, resp.Status)
}

func TestUsecase_Execute_NotUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockbookingRepo(ctrl)
	uc := cancel.NewUsecase(repo)

	resp, err := uc.Execute(t.Context(), uuid.New(), &authjwt.Identity{
		Role: user.RoleAdmin,
	})

	require.Error(t, err)
	require.ErrorIs(t, err, user.ErrUserRoleRequired)
	require.Nil(t, resp)
}

func TestUsecase_Execute_BookingNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockbookingRepo(ctrl)
	uc := cancel.NewUsecase(repo)

	bookingID := uuid.New()

	identity := &authjwt.Identity{
		Role:   user.RoleUser,
		UserID: uuid.New(),
	}

	repo.EXPECT().
		GetByID(gomock.Any(), bookingID).
		Return(&booking.Booking{}, booking.ErrNotFound)

	resp, err := uc.Execute(t.Context(), bookingID, identity)

	require.Equal(t, booking.ErrNotFound, err)
	require.Nil(t, resp)
}

func TestUsecase_Execute_Forbidden(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockbookingRepo(ctrl)
	uc := cancel.NewUsecase(repo)

	bookingID := uuid.New()

	book := &booking.Booking{
		ID:     bookingID,
		UserID: uuid.New(),
		Status: booking.StatusActive,
	}

	identity := &authjwt.Identity{
		Role:   user.RoleUser,
		UserID: uuid.New(),
	}

	repo.EXPECT().
		GetByID(gomock.Any(), bookingID).
		Return(book, nil)

	resp, err := uc.Execute(t.Context(), bookingID, identity)

	require.Error(t, err)
	require.ErrorIs(t, err, booking.ErrCancelForbidden)
	require.Nil(t, resp)
}

func TestUsecase_Execute_AlreadyCancelled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockbookingRepo(ctrl)
	uc := cancel.NewUsecase(repo)

	userID := uuid.New()
	bookingID := uuid.New()

	book := &booking.Booking{
		ID:     bookingID,
		UserID: userID,
		Status: booking.StatusCancelled,
	}

	identity := &authjwt.Identity{
		Role:   user.RoleUser,
		UserID: userID,
	}

	repo.EXPECT().
		GetByID(gomock.Any(), bookingID).
		Return(book, nil)

	resp, err := uc.Execute(t.Context(), bookingID, identity)

	require.NoError(t, err)
	require.Equal(t, booking.StatusCancelled, resp.Status)
}
