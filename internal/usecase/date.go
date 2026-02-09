package usecase

import (
	"fmt"
	"time"
)

const MonthLayout = "01-2006"

// ParseMonthDate parses month-year (MM-YYYY) into first day of month in UTC.
func ParseMonthDate(s string) (time.Time, error) {
	t, err := time.ParseInLocation(MonthLayout, s, time.UTC)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid month format, expected MM-YYYY: %w", err)
	}
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC), nil
}

func FormatMonthDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(MonthLayout)
}
