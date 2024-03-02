package month

import (
	"errors"
	"io/fs"
	"os"
	"time"

	"git.sr.ht/~kota/calendar/config"
	"git.sr.ht/~kota/calendar/date"
	"git.sr.ht/~kota/calendar/note"
	tea "github.com/charmbracelet/bubbletea"
)

// styledDays represents days of the month which should use non-standard styling.
//
// This includes matched holidays, days with a matched keyword, or just any
// days with a written note if configured to style them differently.
// It does not include "today" or the selected day as these do not need to be
// parsed / loaded concurrently.
type styledDays map[string]config.Style

// Match attempts to match a given time with a styled day.
func (sd styledDays) Match(t time.Time) (config.Style, bool) {
	for date, style := range sd {
		if date == t.Format("2006-01-02") {
			return style, true
		}
	}
	return config.Style{}, false
}

type styledDaysMsg struct {
	month      time.Time
	styledDays styledDays
}

// loadStyledDays reads every note file for the given Month to create a tea.Msg
// with days that should be styled differently (matching a keyword, holiday,
// etc).
func (m Month) loadStyledDays() tea.Msg {
	var msg styledDaysMsg
	msg.month = m.date

	sd := make(styledDays)
	last := date.LastDay(m.date)
	for i := 1; i <= last.Day(); i++ {
		t := time.Date(m.date.Year(), m.date.Month(),
			i, 0, 0, 0, 0,
			m.date.Location())
		// Process noted days.
		if !m.config.NotedStyle.Blank() {
			if note.Exists(t, m.config.NoteDir) {
				sd[t.Format("2006-01-02")] = m.config.NotedStyle
			}
		}

		// Process holidays.
		if h, ok := m.holidays.Match(t); ok {
			sd[t.Format("2006-01-02")] = config.Style{Color: h.Color}
		}

		// Process keywords.
		if len(m.config.Keywords) != 0 {
			path := note.Path(t, m.config.NoteDir)
			f, err := os.Open(path)
			if err != nil && !errors.Is(err, fs.ErrNotExist) {
				continue
			}
			defer f.Close()

			if k, ok := m.config.Keywords.Match(f); ok {
				sd[t.Format("2006-01-02")] = config.Style{Color: k.Color}
			}
		}
	}

	msg.styledDays = sd
	return msg
}
