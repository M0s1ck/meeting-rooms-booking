package slot_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/slot"
)

func Test_New_Error(t *testing.T) {
	roomID := uuid.New()
	startAt := time.Date(2026, time.April, 4, 23, 50, 0, 0, time.UTC)

	slt, err := slot.New(roomID, startAt, startAt.Add(-time.Hour*24))
	require.Equal(t, err, slot.ErrSlotStartTooLate)
	require.Nil(t, slt)
}

func Test_New_UsesDeterministicID(t *testing.T) {
	roomID := uuid.New()
	startAt := time.Date(2026, time.April, 4, 9, 0, 0, 0, time.FixedZone("UTC+3", 3*60*60))

	first, err := slot.New(roomID, startAt, time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)

	second, err := slot.New(roomID, startAt.UTC(), time.Date(2026, time.April, 2, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)

	require.Equal(t, first.ID, second.ID)
	require.Equal(t, first.StartAt, second.StartAt)
	require.Equal(t, first.EndAt, second.EndAt)
}
