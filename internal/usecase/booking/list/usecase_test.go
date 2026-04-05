package list_test

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/booking"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
	"github.com/internships-backend/test-backend-M0s1ck/internal/service/authjwt"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/booking/list"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/booking/list/mocks"
)

func TestUsecase_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockbookingRepo(ctrl)
	uc := list.NewUsecase(repo)

	now := time.Now().UTC()

	b1, err := booking.NewActive(uuid.New(), uuid.New(), now)
	require.NoError(t, err)

	b2, err := booking.NewActive(uuid.New(), uuid.New(), now)
	require.NoError(t, err)

	req := &list.Request{
		Page:     2,
		PageSize: 10,
	}

	identity := &authjwt.Identity{
		Role: user.RoleAdmin,
	}

	repo.EXPECT().
		List(gomock.Any(), 2, 10).
		Return([]booking.Booking{*b1, *b2}, 2, nil)

	resp, err := uc.Execute(t.Context(), req, identity)

	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Len(t, resp.Bookings, 2)
	require.Equal(t, 2, resp.Pagination.Page)
	require.Equal(t, 10, resp.Pagination.PageSize)
	require.Equal(t, 2, resp.Pagination.Total)

	require.Equal(t, b1.ID, resp.Bookings[0].ID)
	require.Equal(t, b2.Status, resp.Bookings[1].Status)
}

func TestUsecase_Execute_NotAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockbookingRepo(ctrl)
	uc := list.NewUsecase(repo)

	req := &list.Request{}
	identity := &authjwt.Identity{
		Role: user.RoleUser,
	}

	resp, err := uc.Execute(t.Context(), req, identity)

	require.Error(t, err)
	require.ErrorIs(t, err, user.ErrAdminRoleRequired)
	require.Nil(t, resp)
}

func TestUsecase_Execute_InvalidPage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockbookingRepo(ctrl)
	uc := list.NewUsecase(repo)

	req := &list.Request{
		Page: -1,
	}

	identity := &authjwt.Identity{
		Role: user.RoleAdmin,
	}

	resp, err := uc.Execute(t.Context(), req, identity)

	require.Error(t, err)
	require.ErrorIs(t, err, booking.ErrInvalidPage)
	require.Nil(t, resp)
}

func TestUsecase_Execute_InvalidPageSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockbookingRepo(ctrl)
	uc := list.NewUsecase(repo)

	req := &list.Request{
		PageSize: 1000,
	}

	identity := &authjwt.Identity{
		Role: user.RoleAdmin,
	}

	resp, err := uc.Execute(t.Context(), req, identity)

	require.Error(t, err)
	require.ErrorIs(t, err, booking.ErrInvalidPageSize)
	require.Nil(t, resp)
}

func TestUsecase_Execute_DefaultPagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockbookingRepo(ctrl)
	uc := list.NewUsecase(repo)

	identity := &authjwt.Identity{
		Role: user.RoleAdmin,
	}

	repo.EXPECT().
		List(gomock.Any(), 1, 20).
		Return([]booking.Booking{}, 0, nil)

	resp, err := uc.Execute(t.Context(), nil, identity)

	require.NoError(t, err)
	require.Equal(t, 1, resp.Pagination.Page)
	require.Equal(t, 20, resp.Pagination.PageSize)
}

func TestUsecase_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockbookingRepo(ctrl)
	uc := list.NewUsecase(repo)

	req := &list.Request{}
	identity := &authjwt.Identity{
		Role: user.RoleAdmin,
	}

	repo.EXPECT().
		List(gomock.Any(), 1, 20).
		Return(nil, 0, errors.New("db error"))

	resp, err := uc.Execute(t.Context(), req, identity)

	require.Error(t, err)
	require.Nil(t, resp)
}
