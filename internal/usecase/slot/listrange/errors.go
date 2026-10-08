package listrange

import (
	"errors"
	"fmt"
)

// MaxRangeDays limits how many days one request may cover: it bounds the response size
// and the amount of on-demand slot generation for dates beyond the horizon.
const MaxRangeDays = 31

var (
	ErrFromAfterTo   = errors.New("'from' must not be after 'to'")
	ErrRangeTooLarge = fmt.Errorf("date range must not exceed %d days", MaxRangeDays)
)
