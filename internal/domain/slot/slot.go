package slot

import (
	"time"

	"github.com/google/uuid"
)

type Slot struct {
	ID        uuid.UUID
	RoomID    uuid.UUID
	StartAt   time.Time
	EndAt     time.Time
	CreatedAt time.Time
}

const Duration = 30 * time.Minute

func New(id uuid.UUID, roomID uuid.UUID, startAt time.Time, now time.Time) (*Slot, error) {
	startAt = startAt.UTC()

	endAt, err := getEndAt(startAt)
	if err != nil {
		return nil, err
	}

	return &Slot{
		ID:        id,
		RoomID:    roomID,
		StartAt:   startAt,
		EndAt:     endAt,
		CreatedAt: now.UTC(),
	}, nil
}

func getEndAt(start time.Time) (time.Time, error) {
	end := start.Add(Duration).UTC()

	if end.Hour()*60+end.Minute() < start.Hour()*60+start.Minute() &&
		!(end.Hour() == 0 && end.Minute() == 0) {

		return time.Time{}, ErrSlotStartTooLate
	}

	return end, nil
}
