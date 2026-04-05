package create_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/booking"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/slot"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
	"github.com/internships-backend/test-backend-M0s1ck/internal/service/authjwt"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/booking/create"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/booking/create/mocks"
)

func TestUsecase_Execute_Success_NoLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bookingRepo := mocks.NewMockbookingRepo(ctrl)
	slotRepo := mocks.NewMockslotRepo(ctrl)
	conf := mocks.NewMockconferenceLinkProvider(ctrl)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	uc := create.NewUsecase(bookingRepo, slotRepo, conf, logger)

	userID := uuid.New()
	slotID := uuid.New()

	now := time.Now().UTC()

	sl := &slot.Slot{
		ID:      slotID,
		StartAt: now.Add(time.Hour),
	}

	req := &create.Request{
		SlotID: slotID,
	}

	identity := &authjwt.Identity{
		Role:   user.RoleUser,
		UserID: userID,
	}

	slotRepo.EXPECT().
		GetByID(gomock.Any(), slotID).
		Return(sl, nil)

	bookingRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, b *booking.Booking) error {
			require.Equal(t, slotID, b.SlotID)
			require.Equal(t, userID, b.UserID)
			require.Equal(t, booking.StatusActive, b.Status)
			require.Nil(t, b.ConferenceLink)
			return nil
		})

	resp, err := uc.Execute(t.Context(), req, identity)

	require.NoError(t, err)
	require.Equal(t, slotID, resp.SlotID)
	require.Equal(t, userID, resp.UserID)
}

func TestUsecase_Execute_Success_WithLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bookingRepo := mocks.NewMockbookingRepo(ctrl)
	slotRepo := mocks.NewMockslotRepo(ctrl)
	conf := mocks.NewMockconferenceLinkProvider(ctrl)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	uc := create.NewUsecase(bookingRepo, slotRepo, conf, logger)

	userID := uuid.New()
	slotID := uuid.New()
	now := time.Now().UTC()

	link := "https://meet"

	sl := &slot.Slot{
		ID:      slotID,
		StartAt: now.Add(time.Hour),
	}

	req := &create.Request{
		SlotID:               slotID,
		CreateConferenceLink: true,
	}

	identity := &authjwt.Identity{
		Role:   user.RoleUser,
		UserID: userID,
	}

	slotRepo.EXPECT().
		GetByID(gomock.Any(), slotID).
		Return(sl, nil)

	conf.EXPECT().
		CreateLink(gomock.Any(), slotID, userID).
		Return(link, nil)

	bookingRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, b *booking.Booking) error {
			require.NotNil(t, b.ConferenceLink)
			require.Equal(t, link, *b.ConferenceLink)
			return nil
		})

	resp, err := uc.Execute(t.Context(), req, identity)

	require.NoError(t, err)
	require.NotNil(t, resp.ConferenceLink)
}

func TestUsecase_Execute_LinkError_NonBlocking(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bookingRepo := mocks.NewMockbookingRepo(ctrl)
	slotRepo := mocks.NewMockslotRepo(ctrl)
	conf := mocks.NewMockconferenceLinkProvider(ctrl)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	uc := create.NewUsecase(bookingRepo, slotRepo, conf, logger)

	userID := uuid.New()
	slotID := uuid.New()
	now := time.Now().UTC()

	sl := &slot.Slot{
		ID:      slotID,
		StartAt: now.Add(time.Hour),
	}

	req := &create.Request{
		SlotID:               slotID,
		CreateConferenceLink: true,
	}

	identity := &authjwt.Identity{
		Role:   user.RoleUser,
		UserID: userID,
	}

	slotRepo.EXPECT().
		GetByID(gomock.Any(), slotID).
		Return(sl, nil)

	conf.EXPECT().
		CreateLink(gomock.Any(), slotID, userID).
		Return("", errors.New("fail"))

	bookingRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	resp, err := uc.Execute(t.Context(), req, identity)

	require.NoError(t, err)
	require.Nil(t, resp.ConferenceLink)
}

func TestUsecase_Execute_SlotInPast(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bookingRepo := mocks.NewMockbookingRepo(ctrl)
	slotRepo := mocks.NewMockslotRepo(ctrl)
	conf := mocks.NewMockconferenceLinkProvider(ctrl)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	uc := create.NewUsecase(bookingRepo, slotRepo, conf, logger)

	slotID := uuid.New()

	sl := &slot.Slot{
		ID:      slotID,
		StartAt: time.Now().Add(-time.Hour),
	}

	slotRepo.EXPECT().
		GetByID(gomock.Any(), slotID).
		Return(sl, nil)

	resp, err := uc.Execute(t.Context(), &create.Request{
		SlotID: slotID,
	}, &authjwt.Identity{
		Role:   user.RoleUser,
		UserID: uuid.New(),
	})

	require.Error(t, err)
	require.ErrorIs(t, err, booking.ErrSlotInPast)
	require.Nil(t, resp)
}

func TestUsecase_Execute_NotUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := create.NewUsecase(
		mocks.NewMockbookingRepo(ctrl),
		mocks.NewMockslotRepo(ctrl),
		mocks.NewMockconferenceLinkProvider(ctrl),
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)

	resp, err := uc.Execute(t.Context(), &create.Request{}, &authjwt.Identity{
		Role: user.RoleAdmin,
	})

	require.Error(t, err)
	require.ErrorIs(t, err, user.ErrUserRoleRequired)
	require.Nil(t, resp)
}

func TestUsecase_Execute_SlotRepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	slotRepo := mocks.NewMockslotRepo(ctrl)

	uc := create.NewUsecase(
		mocks.NewMockbookingRepo(ctrl),
		slotRepo,
		mocks.NewMockconferenceLinkProvider(ctrl),
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)

	slotID := uuid.New()

	slotRepo.EXPECT().
		GetByID(gomock.Any(), slotID).
		Return(nil, errors.New("db error"))

	resp, err := uc.Execute(t.Context(), &create.Request{
		SlotID: slotID,
	}, &authjwt.Identity{
		Role:   user.RoleUser,
		UserID: uuid.New(),
	})

	require.Error(t, err)
	require.Nil(t, resp)
}

func TestUsecase_Execute_CreateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bookingRepo := mocks.NewMockbookingRepo(ctrl)
	slotRepo := mocks.NewMockslotRepo(ctrl)

	uc := create.NewUsecase(
		bookingRepo,
		slotRepo,
		mocks.NewMockconferenceLinkProvider(ctrl),
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)

	slotID := uuid.New()

	sl := &slot.Slot{
		ID:      slotID,
		StartAt: time.Now().Add(time.Hour),
	}

	slotRepo.EXPECT().
		GetByID(gomock.Any(), slotID).
		Return(sl, nil)

	bookingRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(errors.New("fail"))

	resp, err := uc.Execute(t.Context(), &create.Request{
		SlotID: slotID,
	}, &authjwt.Identity{
		Role:   user.RoleUser,
		UserID: uuid.New(),
	})

	require.Error(t, err)
	require.Nil(t, resp)
}
