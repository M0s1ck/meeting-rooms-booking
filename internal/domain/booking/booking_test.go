package booking_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/booking"
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

func Test_SetConfLink(t *testing.T) {
	tests := []struct {
		name string
		link string
		err  error
	}{
		{
			name: "blank_link",
			link: "  ",
			err:  booking.ErrConfLinkBlank,
		},
		{
			name: "valid",
			link: "https://mock-conference.local/1737f916-c6c6-514e-b6e7-eee9cec14d0a/b566a8f5-4fbc-42c7-ba6b-3a85dcb6ed6c/380cae68-65a4-4d85-9831-46f949fc92ad",
			err:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slotId, userID := uuid.New(), uuid.New()
			now := time.Now().UTC()

			book, err := booking.NewActive(slotId, userID, now)
			require.Nil(t, err)

			err = book.SetConferenceLink(tt.link)
			require.Equal(t, tt.err, err)
		})
	}
}
