// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package month

import (
	"time"

	"git.sr.ht/~kota/calendar/date"
	tea "github.com/charmbracelet/bubbletea"
)

// move the selection based on a keypress.
func (m *Month) move(msg tea.KeyMsg) {
	switch {
	case m.config.KeyLastSunday.Contains(msg.String()):
		m.selected = date.LastSunday(m.selected)
	case m.config.KeyNextSaturday.Contains(msg.String()):
		m.selected = date.NextSaturday(m.selected)
	case m.config.KeyNextSunday.Contains(msg.String()):
		m.selected = date.NextSunday(m.selected)
	case m.config.KeyMonthDown.Contains(msg.String()):
		m.selected = date.NextMonth(m.selected)
	case m.config.KeyMonthUp.Contains(msg.String()):
		m.selected = date.LastMonth(m.selected)
	}

	if m.layout == LayoutGrid {
		m.gridMove(msg)
	} else {
		m.columnMove(msg)
	}
}

func (m *Month) columnMove(msg tea.KeyMsg) {
	switch {
	case m.config.KeySelectLeft.Contains(msg.String()):
		m.selected = m.selected.AddDate(0, 0, -1)
	case m.config.KeySelectRight.Contains(msg.String()):
		m.selected = m.selected.AddDate(0, 0, 1)
	case m.config.KeySelectDown.Contains(msg.String()):
		m.selected = m.selected.AddDate(0, 0, 7)
	case m.config.KeySelectUp.Contains(msg.String()):
		m.selected = m.selected.AddDate(0, 0, -7)
	}
}

func (m *Month) gridMove(msg tea.KeyMsg) {
	switch {
	case m.config.KeySelectLeft.Contains(msg.String()):
		m.selected = m.selected.AddDate(0, 0, -1)
	case m.config.KeySelectRight.Contains(msg.String()):
		m.selected = m.selected.AddDate(0, 0, 1)
	case m.config.KeySelectDown.Contains(msg.String()):
		m.selected = gridDown(m.selected)
	case m.config.KeySelectUp.Contains(msg.String()):
		m.selected = gridUp(m.selected)
	}
}

func gridDown(t time.Time) time.Time {
	if date.LastWeek(t) {
		down := time.Date(t.Year(), t.Month()+4, 1, 0, 0, 0, 0, t.Location())
		offset := int(t.Weekday() - down.Weekday())
		if offset < 0 {
			offset = offset + 7
		}
		return time.Date(
			t.Year(),
			t.Month()+4,
			1+offset,
			0,
			0,
			0,
			0,
			t.Location(),
		)
	}
	return t.AddDate(0, 0, 7)
}

func gridUp(t time.Time) time.Time {
	if date.FirstWeek(t) {
		up := time.Date(t.Year(), t.Month()-3, 0, 0, 0, 0, 0, t.Location())
		offset := int(t.Weekday() - up.Weekday())
		if offset > 0 {
			offset = offset - 7
		}
		return time.Date(
			t.Year(),
			t.Month()-3,
			0+offset,
			0,
			0,
			0,
			0,
			t.Location(),
		)
	}
	return t.AddDate(0, 0, -7)
}
