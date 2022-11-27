// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package date

import "time"

// LastMonth returns a time representing the previous month from time t.
// The day of the month will be the same, or truncated to the last day.
func LastMonth(t time.Time) time.Time {
	day := t.Day()
	lastDay := DaysIn(t.Month()-1, t.Year())
	if day > lastDay {
		day = lastDay
	}
	return time.Date(t.Year(), t.Month()-1, day, 0, 0, 0, 0, t.Location())
}

// NextMonth returns a time representing the next month after time t.
// The day of the month will be the same, or truncated to the last day.
func NextMonth(t time.Time) time.Time {
	day := t.Day()
	lastDay := DaysIn(t.Month()+1, t.Year())
	if day > lastDay {
		day = lastDay
	}
	return time.Date(t.Year(), t.Month()+1, day, 0, 0, 0, 0, t.Location())
}

// DaysIn reports the number of days in the month of time t.
func DaysIn(m time.Month, year int) int {
	// The reason it works is that we generate a date one month from the target,
	// but set the day of month to 0. Days are 1-indexed, so this has the effect
	// of rolling back one day to the last day of the previous month.
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// SameMonth returns true if both times are in the same month and year.
func SameMonth(x, y time.Time) bool {
	if x.Year() == y.Year() && int(x.Month()) == int(y.Month()) {
		return true
	}
	return false
}

// FirstDay returns a time representing the first day of the month for time t.
func FirstDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// LastDay returns a time representing the last day of the month for time t.
func LastDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location())
}

// LastSunday returns a time representing the last Sunday for time t.
func LastSunday(t time.Time) time.Time {
	offset := int(t.Weekday())
	if offset == 0 {
		// If it's already sunday, go to the previous one.
		offset = 7
	}
	return time.Date(
		t.Year(),
		t.Month(),
		t.Day()-offset,
		0, 0, 0, 0,
		t.Location(),
	)
}

// NextSunday returns a time representing the next Sunday for time t.
func NextSunday(t time.Time) time.Time {
	offset := int(7 - t.Weekday())
	if offset == 0 {
		// If it's already sunday, go to the next one.
		offset = 7
	}
	return time.Date(
		t.Year(),
		t.Month(),
		t.Day()+offset,
		0, 0, 0, 0,
		t.Location(),
	)
}

// NextSaturday returns a time representing the next Saturday for time t.
func NextSaturday(t time.Time) time.Time {
	offset := int(6 - t.Weekday())
	if offset == 0 {
		// If it's already saturday, go to the next one.
		offset = 7
	}
	return time.Date(
		t.Year(),
		t.Month(),
		t.Day()+offset,
		0, 0, 0, 0,
		t.Location(),
	)
}
