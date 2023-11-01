// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package calendar

import (
	"time"

	"git.sr.ht/~kota/calendar/date"
	"git.sr.ht/~kota/calendar/month"
)

// resize the number of months being displayed to fill the window size.
func (c Calendar) resize() Calendar {
	want := 1
	if c.height > 3*month.MonthHeight {
		want = 3
		if c.previewMode == previewModeHidden {
			if c.width > 4*month.MonthWidth+c.config.LeftPadding*3 {
				want = 12
			}
		}
	}

	switch want {
	case 3:
		c.months = c.resizeThree()
	case 12:
		c.months = c.resizeTwelve()
	default:
		c.months = c.resizeOne()
	}

	// Restore focus. It gets lost when resizing.
	c.SetFocus(c.previewMode)
	return c
}

func (c Calendar) resizeOne() []month.Month {
	return []month.Month{month.New(
		c.selected,
		c.today,
		c.selected,
		true,
		c.config,
	)}
}

func (c Calendar) resizeThree() []month.Month {
	offset := mod(int(c.today.Month()-c.selected.Month()), 3)
	switch offset {
	case 0:
		return []month.Month{
			month.New(
				date.LastMonth(c.selected),
				c.today,
				c.selected,
				true,
				c.config,
			),
			month.New(
				c.selected,
				c.today,
				c.selected,
				true,
				c.config,
			),
			month.New(
				date.NextMonth(c.selected),
				c.today,
				c.selected,
				true,
				c.config,
			),
		}
	case 1:
		return []month.Month{
			month.New(
				c.selected,
				c.today,
				c.selected,
				true,
				c.config,
			),
			month.New(
				date.NextMonth(c.selected),
				c.today,
				c.selected,
				true,
				c.config,
			),
			month.New(
				date.NextMonth(date.NextMonth(c.selected)),
				c.today,
				c.selected,
				true,
				c.config,
			),
		}
	default: // 2
		return []month.Month{
			month.New(
				date.LastMonth(date.LastMonth(c.selected)),
				c.today,
				c.selected,
				true,
				c.config,
			),
			month.New(
				date.LastMonth(c.selected),
				c.today,
				c.selected,
				true,
				c.config,
			),
			month.New(
				c.selected,
				c.today,
				c.selected,
				true,
				c.config,
			),
		}
	}
}

func (c Calendar) resizeTwelve() []month.Month {
	var months []month.Month
	for i := 1; i <= 12; i++ {
		m := date.Month(time.Month(i), c.selected.Year())
		if date.SameMonth(m, c.selected) {
			months = append(months, month.New(
				m,
				c.today,
				c.selected,
				false,
				c.config,
			))
		} else {
			months = append(months, month.New(
				m,
				c.today,
				c.selected,
				false,
				c.config,
			))
		}
	}
	return months
}

// mod is a fast and simple integer modulo implementation.
func mod(x, y int) int {
	m := x % y
	if m > 0 && y < 0 {
		m -= y
	}
	if m < 0 && y > 0 {
		m += y
	}
	return m
}
