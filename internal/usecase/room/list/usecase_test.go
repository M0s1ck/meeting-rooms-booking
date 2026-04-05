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
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/room/list"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/room/list/mocks"
)

func TestUsecase_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockroomRepo(ctrl)
	uc := list.NewUsecase(repo)

	now := time.Now()

	room1 := &room.Room{
		ID:        uuid.New(),
		Name:      "Room A",
		CreatedAt: now,
	}

	room2 := &room.Room{
		ID:        uuid.New(),
		Name:      "Room B",
		CreatedAt: now,
	}

	repo.EXPECT().
		List(gomock.Any()).
		Return([]*room.Room{room1, room2}, nil)

	resp, err := uc.Execute(context.Background())

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Rooms, 2)

	require.Equal(t, "Room A", resp.Rooms[0].Name)
	require.Equal(t, "Room B", resp.Rooms[1].Name)
}

func TestUsecase_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockroomRepo(ctrl)
	uc := list.NewUsecase(repo)

	repo.EXPECT().
		List(gomock.Any()).
		Return(nil, errors.New("db error"))

	resp, err := uc.Execute(context.Background())

	require.Error(t, err)
	require.Nil(t, resp)
}
