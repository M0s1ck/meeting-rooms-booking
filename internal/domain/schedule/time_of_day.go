package schedule

import "time"

const timeLayout = "15:04"

type TimeOfDay struct {
	minutes int
}

func ParseTimeOfDay(value string) (TimeOfDay, error) {
	parsed, err := time.Parse(timeLayout, value)
	if err != nil {
		return TimeOfDay{}, ErrInvalidTimeFmt
	}

	return TimeOfDay{
		minutes: parsed.Hour()*60 + parsed.Minute(),
	}, nil
}

func (t TimeOfDay) String() string {
	hours := t.minutes / 60
	minutes := t.minutes % 60

	return time.Date(0, 1, 1, hours, minutes, 0, 0, time.UTC).Format(timeLayout)
}

func (t TimeOfDay) Before(other TimeOfDay) bool {
	return t.minutes < other.minutes
}

func (t TimeOfDay) Minutes() int {
	return t.minutes
}

func (t TimeOfDay) Hours() int {
	return t.minutes / 60
}

func (t TimeOfDay) MinutesOfHour() int {
	return t.minutes % 60
}
