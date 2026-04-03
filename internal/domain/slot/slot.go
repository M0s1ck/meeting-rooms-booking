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

var idNamespace = uuid.MustParse("9f90954f-aeb5-4b72-8d4e-f20ff0c6b5f5")

func New(roomID uuid.UUID, startAt time.Time, now time.Time) (*Slot, error) {
	startAt = startAt.UTC()

	endAt, err := getEndAt(startAt)
	if err != nil {
		return nil, err
	}

	return &Slot{
		ID:        IDFor(roomID, startAt),
		RoomID:    roomID,
		StartAt:   startAt,
		EndAt:     endAt,
		CreatedAt: now.UTC(),
	}, nil
}

func IDFor(roomID uuid.UUID, startAt time.Time) uuid.UUID {
	startAt = startAt.UTC()
	key := roomID.String() + "|" + startAt.Format(time.RFC3339)
	return uuid.NewSHA1(idNamespace, []byte(key))
}

func getEndAt(start time.Time) (time.Time, error) {
	end := start.Add(Duration).UTC()

	if end.Hour()*60+end.Minute() < start.Hour()*60+start.Minute() {
		return time.Time{}, ErrSlotStartTooLate
	}

	return end, nil
}
