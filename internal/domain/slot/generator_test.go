package slot_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/schedule"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/slot"
)

const timeLayout = "2006-01-02 15:04"

func Test_Generate(t *testing.T) {
	type testCase struct {
		name  string
		sched *schedule.Schedule
		from  time.Time
		to    time.Time
		res   []slot.Slot
	}

	tests := []func() testCase{
		func() testCase {
			name := "create_two_week_horizon"

			shCreateAt, _ := time.Parse(timeLayout, "2026-03-10 09:44") // Tuesday
			sched := newSched([]time.Weekday{time.Wednesday, time.Friday}, "9:00", "10:30", shCreateAt)
			from := shCreateAt
			to := shCreateAt.AddDate(0, 0, 14) // 2 weeks after schedule creation

			slots := []slot.Slot{
				newSlot(sched.RoomID, "2026-03-11 09:00"),
				newSlot(sched.RoomID, "2026-03-11 09:30"),
				newSlot(sched.RoomID, "2026-03-11 10:00"),

				newSlot(sched.RoomID, "2026-03-13 09:00"),
				newSlot(sched.RoomID, "2026-03-13 09:30"),
				newSlot(sched.RoomID, "2026-03-13 10:00"),

				newSlot(sched.RoomID, "2026-03-18 09:00"),
				newSlot(sched.RoomID, "2026-03-18 09:30"),
				newSlot(sched.RoomID, "2026-03-18 10:00"),

				newSlot(sched.RoomID, "2026-03-20 09:00"),
				newSlot(sched.RoomID, "2026-03-20 09:30"),
				newSlot(sched.RoomID, "2026-03-20 10:00"),
			}

			return testCase{
				name:  name,
				sched: sched,
				from:  from,
				to:    to,
				res:   slots,
			}
		},
		func() testCase {
			name := "created_during_the_day"

			shCreateAt, _ := time.Parse(timeLayout, "2026-03-11 09:44") // Wednesday during the worktime
			sched := newSched([]time.Weekday{time.Wednesday, time.Friday}, "9:00", "10:30", shCreateAt)
			from := shCreateAt
			to := shCreateAt.AddDate(0, 0, 7) // week later, Wednesday, 9:44

			slots := []slot.Slot{
				newSlot(sched.RoomID, "2026-03-11 10:00"),

				newSlot(sched.RoomID, "2026-03-13 09:00"),
				newSlot(sched.RoomID, "2026-03-13 09:30"),
				newSlot(sched.RoomID, "2026-03-13 10:00"),

				newSlot(sched.RoomID, "2026-03-18 09:00"),
			}

			return testCase{
				name:  name,
				sched: sched,
				from:  from,
				to:    to,
				res:   slots,
			}
		},
		func() testCase {
			name := "add_missing_horizon"

			shCreateAt, _ := time.Parse(timeLayout, "2026-03-10 09:44")
			sched := newSched([]time.Weekday{time.Wednesday, time.Friday}, "9:00", "10:30", shCreateAt)

			from := time.Date(2026, time.March, 11, 10, 0, 0, 0, time.UTC) // Wednesday 10:00
			to := time.Date(2026, time.March, 13, 10, 12, 0, 0, time.UTC)  // Friday 10:12

			slots := []slot.Slot{
				newSlot(sched.RoomID, "2026-03-11 10:00"),

				newSlot(sched.RoomID, "2026-03-13 09:00"),
				newSlot(sched.RoomID, "2026-03-13 09:30"),
			}

			return testCase{
				name:  name,
				sched: sched,
				from:  from,
				to:    to,
				res:   slots,
			}
		},
		func() testCase {
			name := "no_space_for_a_slot"

			shCreateAt, _ := time.Parse(timeLayout, "2026-03-10 09:44")
			sched := newSched([]time.Weekday{time.Wednesday, time.Friday}, "9:00", "10:30", shCreateAt)

			from := time.Date(2026, time.March, 11, 10, 30, 0, 0, time.UTC) // Wednesday 10:30
			to := time.Date(2026, time.March, 13, 9, 12, 0, 0, time.UTC)    // Friday 9:12

			slots := make([]slot.Slot, 0)

			return testCase{
				name:  name,
				sched: sched,
				from:  from,
				to:    to,
				res:   slots,
			}
		},
		func() testCase {
			name := "schedule_time_not_even"

			shCreateAt, _ := time.Parse(timeLayout, "2026-03-10 09:44") // Tuesday
			sched := newSched([]time.Weekday{time.Wednesday, time.Friday}, "8:55", "10:47", shCreateAt)

			from := shCreateAt
			to := shCreateAt.AddDate(0, 0, 7) // week after schedule creation

			slots := []slot.Slot{
				newSlot(sched.RoomID, "2026-03-11 8:55"),
				newSlot(sched.RoomID, "2026-03-11 9:25"),
				newSlot(sched.RoomID, "2026-03-11 9:55"),

				newSlot(sched.RoomID, "2026-03-13 8:55"),
				newSlot(sched.RoomID, "2026-03-13 9:25"),
				newSlot(sched.RoomID, "2026-03-13 9:55"),
			}

			return testCase{
				name:  name,
				sched: sched,
				from:  from,
				to:    to,
				res:   slots,
			}
		},
		func() testCase {
			name := "schedule_time_is_from_and_to"

			shCreateAt, _ := time.Parse(timeLayout, "2026-03-10 09:44") // Tuesday
			sched := newSched([]time.Weekday{time.Wednesday, time.Friday}, "9:00", "10:30", shCreateAt)

			from := time.Date(2026, time.March, 11, 9, 0, 0, 0, time.UTC) // Wednesday 9:00
			to := time.Date(2026, time.March, 11, 10, 30, 0, 0, time.UTC) // Wednesday 10:30

			slots := []slot.Slot{
				newSlot(sched.RoomID, "2026-03-11 09:00"),
				newSlot(sched.RoomID, "2026-03-11 09:30"),
				newSlot(sched.RoomID, "2026-03-11 10:00"),
			}

			return testCase{
				name:  name,
				sched: sched,
				from:  from,
				to:    to,
				res:   slots,
			}
		},
		func() testCase {
			name := "week_days_right_order"

			shCreateAt, _ := time.Parse(timeLayout, "2026-03-07 09:44") // Saturday
			sched := newSched([]time.Weekday{time.Wednesday, time.Friday, time.Sunday}, "10:00", "10:30", shCreateAt)

			from := shCreateAt          // Saturday
			to := from.AddDate(0, 0, 5) // Thursday

			slots := []slot.Slot{
				newSlot(sched.RoomID, "2026-03-08 10:00"), // Sunday
				newSlot(sched.RoomID, "2026-03-11 10:00"), // Wednesday
			}

			return testCase{
				name:  name,
				sched: sched,
				from:  from,
				to:    to,
				res:   slots,
			}
		},
	}

	generator := new(slot.Generator)

	for _, setup := range tests {
		tt := setup()
		t.Run(tt.name, func(t *testing.T) {

			slots := generator.Generate(*tt.sched, tt.from, tt.to)
			require.Equal(t, len(tt.res), len(slots))

			for i := range slots {
				require.Equal(t, tt.sched.RoomID, slots[i].RoomID)
				require.Equal(t, tt.res[i].StartAt, slots[i].StartAt)
				require.Equal(t, tt.res[i].EndAt, slots[i].EndAt)
			}
		})
	}
}

func newSched(days []time.Weekday, start, end string, createdAt time.Time) *schedule.Schedule {
	startAt, _ := schedule.ParseTimeOfDay(start)
	endAt, _ := schedule.ParseTimeOfDay(end)

	normDays := make([]int, len(days))
	for i, d := range days {
		if d == time.Sunday {
			normDays[i] = 7
		} else {
			normDays[i] = int(d)
		}
	}

	sched, _ := schedule.New(uuid.New(), normDays, startAt, endAt, createdAt)
	return sched
}

func newSlot(roomID uuid.UUID, start string) slot.Slot {
	startAt, _ := time.Parse(timeLayout, start)
	endAt := startAt.Add(slot.Duration)

	return slot.Slot{
		RoomID:  roomID,
		StartAt: startAt,
		EndAt:   endAt,
	}
}
