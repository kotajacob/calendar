// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package calendar

import "git.sr.ht/~kota/calendar/month"

// resize the number of months being displayed to fill the window size.
func (c Calendar) resize() Calendar {
	var want int
	switch {
	case c.height > 3*month.MonthHeight:
		want = 3
	default:
		want = 1
	}

	switch want {
	case 3:
		c.months = c.resizeThree()
	default:
		c.months = c.resizeOne()
	}

	// Restore focus. It gets lots when resizing.
	c.SetFocus(c.focus)
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
				lastMonth(c.selected),
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
				nextMonth(c.selected),
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
				nextMonth(c.selected),
				c.today,
				c.selected,
				true,
				c.config,
			),
			month.New(
				nextMonth(nextMonth(c.selected)),
				c.today,
				c.selected,
				true,
				c.config,
			),
		}
	default: // 2
		return []month.Month{
			month.New(
				lastMonth(lastMonth(c.selected)),
				c.today,
				c.selected,
				true,
				c.config,
			),
			month.New(
				lastMonth(c.selected),
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
