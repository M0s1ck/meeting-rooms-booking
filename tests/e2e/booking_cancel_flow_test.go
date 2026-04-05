package e2e_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestE2E_UserCanCancelBooking(t *testing.T) {
	adminToken := dummyLogin(t, baseURL, "admin")
	userToken := dummyLogin(t, baseURL, "user")

	room := createRoom(t, baseURL, adminToken, map[string]any{
		"name":        "Sigma",
		"description": "Cancel flow room",
		"capacity":    8,
	})

	targetDate := nextWeekdayUTC(time.Now().UTC(), time.Monday)

	createSchedule(t, baseURL, adminToken, room.ID, map[string]any{
		"roomId":     room.ID,
		"daysOfWeek": []int{1, 2, 3, 4, 5},
		"startTime":  "10:00",
		"endTime":    "12:00",
	})

	slotsBeforeBooking := listSlots(t, baseURL, userToken, room.ID, targetDate.Format("2006-01-02"))
	require.NotEmpty(t, slotsBeforeBooking)

	targetSlot := slotsBeforeBooking[0]

	created := createBooking(t, baseURL, userToken, map[string]any{
		"slotId": targetSlot.ID,
	})

	require.Equal(t, "active", created.Status)
	require.Equal(t, targetSlot.ID, created.SlotID)

	slotsAfterBooking := listSlots(t, baseURL, userToken, room.ID, targetDate.Format("2006-01-02"))
	requireSlotMissing(t, slotsAfterBooking, targetSlot.ID)

	cancelled := cancelBooking(t, baseURL, userToken, created.ID)
	require.Equal(t, created.ID, cancelled.ID)
	require.Equal(t, "cancelled", cancelled.Status)

	cancelledAgain := cancelBooking(t, baseURL, userToken, created.ID)
	require.Equal(t, created.ID, cancelledAgain.ID)
	require.Equal(t, "cancelled", cancelledAgain.Status)

	slotsAfterCancel := listSlots(t, baseURL, userToken, room.ID, targetDate.Format("2006-01-02"))
	requireSlotPresent(t, slotsAfterCancel, targetSlot.ID)
}
