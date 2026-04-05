package e2e_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestE2E_CreateRoomScheduleAndBooking(t *testing.T) {
	adminToken := dummyLogin(t, baseURL, "admin")
	userToken := dummyLogin(t, baseURL, "user")

	room := createRoom(t, baseURL, adminToken, map[string]any{
		"name":        "Omega",
		"description": "E2E room",
		"capacity":    6,
	})

	targetDate := nextWeekdayUTC(time.Now().UTC(), time.Monday)

	createSchedule(t, baseURL, adminToken, room.ID, map[string]any{
		"roomId":     room.ID,
		"daysOfWeek": []int{1, 2, 3, 4, 5},
		"startTime":  "09:00",
		"endTime":    "11:00",
	})

	slots := listSlots(t, baseURL, userToken, room.ID, targetDate.Format("2006-01-02"))
	require.NotEmpty(t, slots, "expected generated slots for target date")

	booking := createBooking(t, baseURL, userToken, map[string]any{
		"slotId": slots[0].ID,
	})

	require.Equal(t, slots[0].ID, booking.SlotID)
	require.Equal(t, "active", booking.Status)

	slotsAfter := listSlots(t, baseURL, userToken, room.ID, targetDate.Format("2006-01-02"))
	for _, s := range slotsAfter {
		require.NotEqual(t, slots[0].ID, s.ID, "booked slot must disappear from available slots")
	}
}

func nextWeekdayUTC(from time.Time, wd time.Weekday) time.Time {
	d := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1)
	for d.Weekday() != wd {
		d = d.AddDate(0, 0, 1)
	}
	return d
}
