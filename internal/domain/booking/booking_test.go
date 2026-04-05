package booking_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/booking"
	"github.com/stretchr/testify/require"
)

func Test_New_Success(t *testing.T) {
	slotId, userID := uuid.New(), uuid.New()
	now := time.Now().UTC()

	book, err := booking.NewActive(slotId, userID, now)
	require.Nil(t, err)
	require.NotNil(t, book)
	require.Equal(t, slotId, book.SlotID)
	require.Equal(t, userID, book.UserID)
	require.Equal(t, booking.StatusActive, book.Status)
	require.Equal(t, now, book.CreatedAt)
}

func Test_New_Fail(t *testing.T) {
	tests := []struct {
		name   string
		slotId uuid.UUID
		userID uuid.UUID
		now    time.Time
		err    error
	}{
		{
			name:   "invalid_slot_id",
			slotId: uuid.Nil,
			userID: uuid.New(),
			now:    time.Now().UTC(),
			err:    booking.ErrSlotIDRequired,
		},
		{
			name:   "invalid_user_id",
			slotId: uuid.New(),
			userID: uuid.Nil,
			now:    time.Now().UTC(),
			err:    booking.ErrUserIDRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			book, err := booking.NewActive(tt.slotId, tt.userID, time.Now())
			require.Equal(t, tt.err, err)
			require.Nil(t, book)
		})
	}
}
