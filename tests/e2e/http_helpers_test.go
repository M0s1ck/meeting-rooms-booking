package e2e_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func dummyLogin(t *testing.T, baseURL, role string) string {
	t.Helper()

	body := map[string]any{"role": role}
	var resp struct {
		Token string `json:"token"`
	}

	doJSON(t, "POST", baseURL+"/dummyLogin", "", body, http.StatusOK, &resp)
	require.NotEmpty(t, resp.Token)
	return resp.Token
}

type roomDTO struct {
	ID string `json:"id"`
}

func createRoom(t *testing.T, baseURL, token string, body map[string]any) roomDTO {
	t.Helper()

	var resp struct {
		Room roomDTO `json:"room"`
	}

	doJSON(t, "POST", baseURL+"/rooms/create", token, body, http.StatusCreated, &resp)
	require.NotEmpty(t, resp.Room.ID)
	return resp.Room
}

func createSchedule(t *testing.T, baseURL, token, roomID string, body map[string]any) {
	t.Helper()

	doJSON(t, "POST", baseURL+"/rooms/"+roomID+"/schedule/create", token, body, http.StatusCreated, nil)
}

type slotDTO struct {
	ID    string    `json:"id"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func listSlots(t *testing.T, baseURL, token, roomID, date string) []slotDTO {
	t.Helper()

	var resp struct {
		Slots []slotDTO `json:"slots"`
	}

	doJSON(t, "GET", baseURL+"/rooms/"+roomID+"/slots/list?date="+date,
		token, nil, http.StatusOK, &resp)

	return resp.Slots
}

func listSlotsRange(t *testing.T, baseURL, token, roomID, from, to string, wantStatus int) []slotDTO {
	t.Helper()

	var resp struct {
		Slots []slotDTO `json:"slots"`
	}

	var out any
	if wantStatus == http.StatusOK {
		out = &resp
	}

	doJSON(t, "GET", baseURL+"/rooms/"+roomID+"/slots/range?from="+from+"&to="+to,
		token, nil, wantStatus, out)

	return resp.Slots
}

type bookingDTO struct {
	ID     string `json:"id"`
	SlotID string `json:"slotId"`
	Status string `json:"status"`
}

func createBooking(t *testing.T, baseURL, token string, body map[string]any) bookingDTO {
	t.Helper()

	var resp struct {
		Booking bookingDTO `json:"booking"`
	}

	doJSON(t, "POST", baseURL+"/bookings/create", token, body, http.StatusCreated, &resp)
	return resp.Booking
}

func cancelBooking(t *testing.T, baseURL, token, bookingID string) bookingDTO {
	t.Helper()

	var resp struct {
		Booking bookingDTO `json:"booking"`
	}

	doJSON(
		t,
		"POST",
		baseURL+"/bookings/"+bookingID+"/cancel",
		token,
		nil,
		http.StatusOK,
		&resp,
	)

	return resp.Booking
}

func requireSlotMissing(t *testing.T, slots []slotDTO, slotID string) {
	t.Helper()

	for _, s := range slots {
		require.NotEqual(t, slotID, s.ID, "slot should not be available")
	}
}

func requireSlotPresent(t *testing.T, slots []slotDTO, slotID string) {
	t.Helper()

	for _, s := range slots {
		if s.ID == slotID {
			return
		}
	}

	require.Fail(t, "expected slot to be available again", "slotId=%s", slotID)
}

func doJSON(t *testing.T, method, url, token string, reqBody any, wantStatus int, out any) {
	t.Helper()

	var bodyBytes []byte
	var err error
	if reqBody != nil {
		bodyBytes, err = json.Marshal(reqBody)
		require.NoError(t, err)
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(bodyBytes))
	require.NoError(t, err)

	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() {
		_ = resp.Body.Close()
	}()

	require.Equal(t, wantStatus, resp.StatusCode)

	if out != nil {
		err = json.NewDecoder(resp.Body).Decode(out)
		require.NoError(t, err)
	}
}
