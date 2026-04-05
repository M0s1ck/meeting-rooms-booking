package create_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/room"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
	"github.com/internships-backend/test-backend-M0s1ck/internal/service/authjwt"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/room/create"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/room/create/mocks"
)

func TestUsecase_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockroomRepo(ctrl)
	uc := create.NewUsecase(repo)

	name := "Room A"
	desc := "Nice room"
	capac := 10

	req := &create.Request{
		Name:        name,
		Description: &desc,
		Capacity:    &capac,
	}

	identity := &authjwt.Identity{
		Role: user.RoleAdmin,
	}

	repo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, r *room.Room) error {
			require.Equal(t, "Room A", r.Name)
			require.NotZero(t, r.ID)
			return nil
		})

	resp, err := uc.Execute(context.Background(), req, identity)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "Room A", resp.Name)
}

func TestUsecase_Execute_NotAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockroomRepo(ctrl)
	uc := create.NewUsecase(repo)

	req := &create.Request{Name: "Room"}

	identity := &authjwt.Identity{
		Role: user.RoleUser,
	}

	resp, err := uc.Execute(context.Background(), req, identity)

	require.Error(t, err)
	require.ErrorIs(t, err, user.ErrAdminRoleRequired)
	require.Nil(t, resp)
}

func TestUsecase_Execute_InvalidName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockroomRepo(ctrl)
	uc := create.NewUsecase(repo)

	invalidReq := &create.Request{
		Name: "   ",
	}

	identity := &authjwt.Identity{
		Role: user.RoleAdmin,
	}

	resp, err := uc.Execute(context.Background(), invalidReq, identity)

	require.Error(t, err)
	require.Nil(t, resp)
}

func TestUsecase_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockroomRepo(ctrl)
	uc := create.NewUsecase(repo)

	req := &create.Request{
		Name: "Room",
	}

	identity := &authjwt.Identity{
		Role: user.RoleAdmin,
	}

	repo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(errors.New("db error"))

	resp, err := uc.Execute(context.Background(), req, identity)

	require.Error(t, err)
	require.Nil(t, resp)
}
