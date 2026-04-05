package schedule_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/schedule"
)

func Test_New_Success(t *testing.T) {
	roomID := uuid.New()
	now := time.Now()

	from, err := schedule.ParseTimeOfDay("10:00")
	require.NoError(t, err)

	to, err := schedule.ParseTimeOfDay("20:00")
	require.NoError(t, err)

	days := []int{3, 4, 2}

	sched, err := schedule.New(roomID, days, from, to, now)

	require.NoError(t, err)
	require.NotNil(t, sched)
	require.Equal(t, sched.DaysOfWeek, []int{2, 3, 4})
}

func Test_New_InvalidDays(t *testing.T) {
	tests := []struct {
		name string
		days []int
		err  error
	}{
		{
			name: "not_one_to_seven",
			days: []int{1, 3, 5, 8},
			err:  schedule.ErrInvalidDayOfWeek,
		},
		{
			name: "repeating_days",
			days: []int{1, 3, 3, 5},
			err:  schedule.ErrInvalidDayOfWeek,
		},
		{
			name: "empty",
			days: []int{},
			err:  schedule.ErrEmptyDaysOfWeek,
		},
	}

	from, err := schedule.ParseTimeOfDay("1:00")
	require.NoError(t, err)

	to, err := schedule.ParseTimeOfDay("10:30")
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sched, err := schedule.New(uuid.New(), tt.days, from, to, time.Now())

			require.Equal(t, tt.err, err)
			require.Nil(t, sched)
		})
	}
}

func Test_New_InvalidTimeRange(t *testing.T) {
	from, err := schedule.ParseTimeOfDay("20:00")
	require.NoError(t, err)

	to, err := schedule.ParseTimeOfDay("10:30")
	require.NoError(t, err)

	days := []int{1, 2, 3}

	sched, err := schedule.New(uuid.New(), days, from, to, time.Now())

	require.Equal(t, schedule.ErrInvalidTimeRange, err)
	require.Nil(t, sched)
}

func Test_HasWeekDay(t *testing.T) {
	from, err := schedule.ParseTimeOfDay("9:00")
	require.NoError(t, err)

	to, err := schedule.ParseTimeOfDay("10:30")
	require.NoError(t, err)

	days := []int{1, 3, 7}

	sched, err := schedule.New(uuid.New(), days, from, to, time.Now())
	require.NoError(t, err)
	require.NotNil(t, sched)

	for _, day := range []time.Weekday{time.Monday, time.Wednesday, time.Sunday} {
		require.True(t, sched.HasWeekday(day))
	}

	for _, day := range []time.Weekday{time.Tuesday, time.Thursday, time.Saturday} {
		require.False(t, sched.HasWeekday(day))
	}
}
