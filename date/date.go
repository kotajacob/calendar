// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package date

import "time"

// LastMonth returns a time representing the previous month from time t.
func LastMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month()-1, 1, 0, 0, 0, 0, t.Location())
}

// NextMonth returns a time representing the next month after time t.
func NextMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
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
