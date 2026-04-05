package room_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/room"
)

func TestRoom_New_Success(t *testing.T) {
	desc := " test "
	capac := 5

	now := time.Now()

	r, err := room.New(" Room ", &desc, &capac, now)

	require.NoError(t, err)
	require.Equal(t, "Room", r.Name)
	require.NotNil(t, r.Description)
	require.Equal(t, "test", *r.Description)
	require.Equal(t, 5, *r.Capacity)
	require.False(t, r.CreatedAt.IsZero())
}

func TestRoom_New_EmptyName(t *testing.T) {
	r, err := room.New("   ", nil, nil, time.Now())

	require.Error(t, err)
	require.ErrorIs(t, err, room.ErrEmptyName)
	require.Nil(t, r)
}

func TestRoom_New_InvalidCapacity(t *testing.T) {
	capac := 0

	r, err := room.New("Room", nil, &capac, time.Now())

	require.Error(t, err)
	require.ErrorIs(t, err, room.ErrInvalidCapacity)
	require.Nil(t, r)
}
