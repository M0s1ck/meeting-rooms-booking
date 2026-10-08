package e2e_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestE2E_ListSlotsRange(t *testing.T) {
	adminToken := dummyLogin(t, baseURL, "admin")
	userToken := dummyLogin(t, baseURL, "user")

	room := createRoom(t, baseURL, adminToken, map[string]any{"name": "Range"})

	// every day 09:00-10:00 UTC -> 2 slots per day
	createSchedule(t, baseURL, adminToken, room.ID, map[string]any{
		"roomId":     room.ID,
		"daysOfWeek": []int{1, 2, 3, 4, 5, 6, 7},
		"startTime":  "09:00",
		"endTime":    "10:00",
	})

	from := time.Now().UTC().AddDate(0, 0, 1)
	to := from.AddDate(0, 0, 2)
	fromS, toS := from.Format("2006-01-02"), to.Format("2006-01-02")

	slots := listSlotsRange(t, baseURL, userToken, room.ID, fromS, toS, http.StatusOK)
	require.Len(t, slots, 6, "3 days x 2 slots")

	// consistent with the single-date endpoint
	var perDay []slotDTO
	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		perDay = append(perDay, listSlots(t, baseURL, userToken, room.ID, d.Format("2006-01-02"))...)
	}
	require.Equal(t, perDay, slots)

	// booked slot disappears from the range
	createBooking(t, baseURL, userToken, map[string]any{"slotId": slots[2].ID})
	after := listSlotsRange(t, baseURL, userToken, room.ID, fromS, toS, http.StatusOK)
	require.Len(t, after, 5)
	requireSlotMissing(t, after, slots[2].ID)

	// validation
	listSlotsRange(t, baseURL, userToken, room.ID, toS, fromS, http.StatusBadRequest)
	listSlotsRange(t, baseURL, userToken, room.ID, fromS, from.AddDate(0, 0, 31).Format("2006-01-02"), http.StatusBadRequest)
	listSlotsRange(t, baseURL, userToken, room.ID, fromS, "not-a-date", http.StatusBadRequest)
}
