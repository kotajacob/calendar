// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package calendar

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"time"

	"git.sr.ht/~kota/calendar/config"
	"git.sr.ht/~kota/calendar/month"
	"git.sr.ht/~kota/calendar/preview"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type focus int

const (
	focusMonths = iota
	focusPreview
)

// editorFinishedMsg is a tea.Msg returned when the spawned editor process
// returns.
type editorFinishedMsg struct{ err error }

// Calendar is the Bubble Tea model for this calendar element. The calendar is a
// thin wrapper for the month elements. It creates and destroys them based on
// the size of the window.
type Calendar struct {
	today       time.Time
	selected    time.Time
	config      *config.Config
	style       lipgloss.Style
	months      []month.Month
	preview     preview.Model
	height      int
	width       int
	focus       focus
	initialized bool
}

// New creates a new calendar model.
func New(conf *config.Config) Calendar {
	now := time.Now()
	m := Calendar{
		today:    now,
		selected: now,
		style: lipgloss.NewStyle().
			PaddingLeft(conf.LeftPadding).
			PaddingRight(conf.RightPadding),
		months: []month.Month{
			month.New(now, now, now, true, conf),
		},
		config: conf,
	}
	m.SetFocus(focusMonths)
	return m
}

// Init the calendar in Bubble Tea.
func (c Calendar) Init() tea.Cmd {
	return nil
}

// Update the calendar in the Bubble Tea update loop.
func (c Calendar) Update(msg tea.Msg) (Calendar, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "h", "left":
			c.SetFocus(focusMonths)
		case "l", "right":
			c.SetFocus(focusMonths)
		case "enter":
			path := c.selected.Format(os.ExpandEnv(c.config.NotePath))
			cmd := tea.ExecProcess(
				exec.Command("vim", path),
				func(err error) tea.Msg {
					return editorFinishedMsg{err: err}
				})
			return c, cmd
		case "tab":
			c.ToggleFocus()
		}
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height

		if !c.initialized {
			note := loadNote(c.selected, c.config.NotePath)
			c.preview = preview.New(note, msg.Width, msg.Height, c.config)
			c.initialized = true
		}

		c = c.resize()
	case editorFinishedMsg:
		// Reload the note when the user exits their editor.
		c = c.Select(c.selected)
	}
	return c.propagate(msg)
}

// propagate an update to all children.
func (c Calendar) propagate(msg tea.Msg) (Calendar, tea.Cmd) {
	var cmds []tea.Cmd
	for i, month := range c.months {
		var cmd tea.Cmd
		c.months[i], cmd = month.Update(msg)
		cmds = append(cmds, cmd)
	}

	for _, month := range c.months {
		// If the month's selection changed we need to update our selection and
		// the selection of all other months to avoid letting the selection get
		// out of sync!
		if !month.Selected().Equal(c.selected) {
			c = c.Select(month.Selected())
		}
	}

	var cmd tea.Cmd
	c.preview, cmd = c.preview.Update(msg)
	cmds = append(cmds, cmd)

	return c, tea.Batch(cmds...)
}

// resize the number of months being displayed to fill the window size.
func (c Calendar) resize() Calendar {
	var want int
	switch {
	case c.height > 3*month.MonthHeight:
		want = 3
	default:
		want = 1
	}

	if len(c.months) == want {
		return c
	}

	// TODO: Revise this to find the "center" month once we have more than 2
	// different size options.
	switch want {
	case 3:
		last := month.New(
			lastMonth(c.selected),
			c.today,
			c.selected,
			true,
			c.config,
		)
		next := month.New(
			nextMonth(c.selected),
			c.today,
			c.selected,
			true,
			c.config,
		)
		c.months = append([]month.Month{last}, c.months...)
		c.months = append(c.months, next)
	default:
		c.months = []month.Month{month.New(
			c.selected,
			c.today,
			c.selected,
			true,
			c.config,
		)}
	}
	return c
}

// Select a different date. This updates the selection on the calendar, all of
// it's months, and the preview window. It also sets the focus to the months.
func (c Calendar) Select(t time.Time) Calendar {
	c.selected = t
	for i, month := range c.months {
		c.months[i] = month.Select(t)
	}
	c.preview = c.preview.SetContent(loadNote(t, c.config.NotePath))
	c.SetFocus(focusMonths)
	return c
}

// ToggleFocus sets the focus to the months or the preview window.
func (c *Calendar) SetFocus(f focus) {
	if f == focusPreview {
		c.preview.Focus()
		for id := range c.months {
			c.months[id].Unfocus()
		}
		c.focus = focusPreview
	} else {
		c.preview.Unfocus()
		for id := range c.months {
			c.months[id].Focus()
		}
		c.focus = focusMonths
	}
}

// ToggleFocus switches the focus between the months and the preview window.
func (c *Calendar) ToggleFocus() {
	if c.focus == focusMonths {
		c.SetFocus(focusPreview)
	} else {
		c.SetFocus(focusMonths)
	}
}

// SetToday sets the today value to a new time.
func (c *Calendar) SetToday(t time.Time) {
	c.today = t
	for id := range c.months {
		c.months[id].SetToday(t)
	}
}

// loadNote reads a note file for a given time.
// The given path should describe where the note would be located for this
// predefined time:
// January 2, 15:04:05, 2006, in time zone seven hours west of GMT
//
// Environment variable, such as $HOME may be used in the path and will be
// expanded appropriately. If the file is missing it is simply treated as an
// empty file. All other errors will return the error string itself (which is
// meant to be displayed to the user).
func loadNote(t time.Time, path string) string {
	formattedPath := t.Format(os.ExpandEnv(path))
	data, err := os.ReadFile(formattedPath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		data = []byte(err.Error())
	}
	return string(data)
}

// lastMonth returns a time representing the previous month from time t.
func lastMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month()-1, 1, 0, 0, 0, 0, t.Location())
}

// nextMonth returns a time representing the next month after time t.
func nextMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
}

// View renders the calendar in its current state.
func (c Calendar) View() string {
	if !c.initialized {
		return ""
	}

	// Build a slice of rendered months.
	var months []string
	for _, month := range c.months {
		months = append(months, month.View())
	}

	return c.style.Render(lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, months...),
		c.preview.View(),
	))
}
