package slot

import (
	"sort"
	"time"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/schedule"
)

type Generator struct {
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(sched schedule.Schedule, from, to time.Time) []Slot {
	slots := make([]Slot, 0)
	if from.After(to) || from.Equal(to) {
		return slots
	}

	dates := getDates(sched, from, to)

	if len(dates) == 0 {
		return slots
	}

	now := time.Now()

	if sched.HasWeekday(from.Weekday()) {
		slots = addHeadingDate(slots, sched, from, to, now)
		dates = dates[1:]
	}

	if len(dates) == 0 {
		return slots
	}

	if !sched.HasWeekday(to.Weekday()) {
		return addFullDates(slots, sched, dates, now)
	}

	slots = addFullDates(slots, sched, dates[:len(dates)-1], now)
	slots = addTailingDate(slots, sched, to, now)
	return slots
}

func addHeadingDate(slots []Slot, sched schedule.Schedule, from, to time.Time, now time.Time) []Slot {
	beginning := getStartTimeShiftedAfterFrom(sched, from)
	end := getDateTimeWithTimeOfDay(sched.EndTime, from)
	if to.Before(end) {
		end = getEndTimeShiftedBeforeTo(sched, to)
	}

	if beginning.After(end) || beginning.Equal(end) {
		return slots
	}

	slots = addContinuousSlots(slots, sched, beginning, end, now)
	return slots
}

func addFullDates(slots []Slot, sched schedule.Schedule, dates []time.Time, now time.Time) []Slot {
	for _, date := range dates {
		begin := getDateTimeWithTimeOfDay(sched.StartTime, date)
		end := getDateTimeWithTimeOfDay(sched.EndTime, date)
		slots = addContinuousSlots(slots, sched, begin, end, now)
	}

	return slots
}

func addTailingDate(slots []Slot, sched schedule.Schedule, to time.Time, now time.Time) []Slot {
	begin := getDateTimeWithTimeOfDay(sched.StartTime, to)
	end := getEndTimeShiftedBeforeTo(sched, to)
	slots = addContinuousSlots(slots, sched, begin, end, now)
	return slots
}

func addContinuousSlots(slots []Slot, sched schedule.Schedule, begin, end time.Time, now time.Time) []Slot {
	for start := begin; start.Add(Duration).Before(end) || start.Add(Duration).Equal(end); start = start.Add(Duration) {
		slot, err := New(sched.RoomID, start, now)
		if err != nil {
			continue
		}

		slots = append(slots, *slot)
	}

	return slots
}

func getDates(sched schedule.Schedule, from, to time.Time) []time.Time {
	dates := make([]time.Time, 0)

	start := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
	to = time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1)

	end := to
	weekAfterStart := start.AddDate(0, 0, 7)
	if weekAfterStart.Before(to) {
		end = weekAfterStart
	}

	for date := start; date.Before(end); date = date.AddDate(0, 0, 1) {
		if sched.HasWeekday(date.Weekday()) {
			dates = addSuchWeekDay(date, to, dates)
		}
	}

	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	return dates
}

func getStartTimeShiftedAfterFrom(sched schedule.Schedule, from time.Time) time.Time {
	startTime := getDateTimeWithTimeOfDay(sched.StartTime, from)

	if startTime.After(from) || startTime.Equal(from) {
		return startTime
	}

	shiftSinceStartTime := getShiftToSchedTime(startTime, from)
	return startTime.Add(shiftSinceStartTime)
}

func getEndTimeShiftedBeforeTo(sched schedule.Schedule, to time.Time) time.Time {
	endTime := getDateTimeWithTimeOfDay(sched.EndTime, to)

	if endTime.Before(to) || endTime.Equal(to) {
		return endTime
	}

	shiftBeforeEndTime := getShiftToSchedTime(to, endTime)
	return endTime.Add(-shiftBeforeEndTime)
}

func getShiftToSchedTime(earlier, later time.Time) time.Duration {
	minutesDiff := int(later.Sub(earlier).Minutes())
	shiftToSchedTime := time.Duration(minutesDiff/int(Duration.Minutes())) * Duration

	if minutesDiff%int(Duration.Minutes()) != 0 {
		shiftToSchedTime += Duration
	}

	return shiftToSchedTime
}

func getDateTimeWithTimeOfDay(td schedule.TimeOfDay, date time.Time) time.Time {
	return time.Date(
		date.Year(), date.Month(), date.Day(),
		td.Hours(), td.MinutesOfHour(),
		0, 0, date.Location(),
	)
}

func addSuchWeekDay(date time.Time, to time.Time, dates []time.Time) []time.Time {
	for d := date; d.Before(to); d = d.AddDate(0, 0, 7) {
		dates = append(dates, d)
	}

	return dates
}
