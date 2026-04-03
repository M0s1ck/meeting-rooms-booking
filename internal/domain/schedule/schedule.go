package schedule

import (
	"slices"
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ID         uuid.UUID
	RoomID     uuid.UUID
	DaysOfWeek []int
	StartTime  TimeOfDay
	EndTime    TimeOfDay
	CreatedAt  time.Time

	WeekDaysSet map[time.Weekday]struct{}
}

func (s *Schedule) HasWeekday(weekday time.Weekday) bool {
	if s.WeekDaysSet == nil {
		s.WeekDaysSet = formWeekdaySet(s.DaysOfWeek)
	}

	_, ok := s.WeekDaysSet[weekday]
	return ok
}

func New(roomID uuid.UUID, daysOfWeek []int, startTime TimeOfDay, endTime TimeOfDay, now time.Time) (*Schedule, error) {
	if roomID == uuid.Nil {
		return nil, ErrRoomIDRequired
	}

	normalizedDays, err := normalizeDaysOfWeek(daysOfWeek)
	if err != nil {
		return nil, err
	}

	if !startTime.Before(endTime) {
		return nil, ErrInvalidTimeRange
	}

	return &Schedule{
		ID:         uuid.New(),
		RoomID:     roomID,
		DaysOfWeek: normalizedDays,
		StartTime:  startTime,
		EndTime:    endTime,
		CreatedAt:  now.UTC(),
	}, nil
}

func normalizeDaysOfWeek(daysOfWeek []int) ([]int, error) {
	if len(daysOfWeek) == 0 {
		return nil, ErrEmptyDaysOfWeek
	}

	seen := make(map[int]struct{}, len(daysOfWeek))
	normalized := make([]int, 0, len(daysOfWeek))

	for _, day := range daysOfWeek {
		if day < 1 || day > 7 {
			return nil, ErrInvalidDayOfWeek
		}

		if _, ok := seen[day]; ok {
			return nil, ErrInvalidDayOfWeek
		}

		seen[day] = struct{}{}
		normalized = append(normalized, day)
	}

	slices.Sort(normalized)
	return normalized, nil
}

func formWeekdaySet(days []int) map[time.Weekday]struct{} {
	wds := make(map[time.Weekday]struct{})

	for _, weekDay := range days {
		wd := time.Weekday(weekDay)
		if wd == 7 {
			wd = time.Sunday
		}
		wds[wd] = struct{}{}
	}

	return wds
}
