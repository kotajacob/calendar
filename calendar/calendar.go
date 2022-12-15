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
	"git.sr.ht/~kota/calendar/date"
	"git.sr.ht/~kota/calendar/month"
	"git.sr.ht/~kota/calendar/preview"
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// previewMode describes if the preview is shown, focused, or hidden.
type previewMode uint8

const (
	previewModeShown previewMode = iota
	previewModeFocused
	previewModeHidden
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
	preview     preview.Preview
	previewMode previewMode
	height      int
	width       int
	initialized bool
}

// New creates a new calendar model.
func New(selected time.Time, conf *config.Config) Calendar {
	now := time.Now()
	m := Calendar{
		today:    now,
		selected: selected,
		style: lipgloss.NewStyle().
			PaddingLeft(conf.LeftPadding).
			PaddingRight(conf.RightPadding),
		months: []month.Month{
			month.New(selected, now, selected, true, conf),
		},
		config: conf,
	}
	m.SetFocus(previewModeShown)
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
		switch {
		case c.config.KeySelectLeft.Contains(msg.String()) ||
			c.config.KeySelectRight.Contains(msg.String()) ||
			c.config.KeyLastSunday.Contains(msg.String()) ||
			c.config.KeyNextSunday.Contains(msg.String()) ||
			c.config.KeyNextSaturday.Contains(msg.String()):
			c.SetFocus(previewModeShown)
		case c.config.KeyEditNote.Contains(msg.String()):
			path := c.selected.Format(os.ExpandEnv(c.config.NotePath))
			cmd := tea.ExecProcess(
				exec.Command("vim", path),
				func(err error) tea.Msg {
					return editorFinishedMsg{err: err}
				})
			return c, cmd
		case c.config.KeyFocusPreview.Contains(msg.String()):
			c.ToggleFocus()
		case c.config.KeyTogglePreview.Contains(msg.String()):
			c.TogglePreview()
		case c.config.KeyYankDate.Contains(msg.String()):
			clipboard.WriteAll(c.selected.Format("2006-01-02"))
		}
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height

		if !c.initialized {
			note := loadNote(c.selected, c.config.NotePath)
			c.preview = preview.New(note, c.config)
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

// Select a different date. This updates the selection on the calendar, all of
// it's months, and the preview window. It also sets the focus to the months.
func (c Calendar) Select(t time.Time) Calendar {
	c.selected = t

	// If the selection has moved "off-screen" we need to rebuild the month
	// list.
	offScreen := true
	for i, m := range c.months {
		if date.SameMonth(m.Date(), t) {
			offScreen = false
		}
		c.months[i] = m.Select(t)
	}

	if offScreen {
		c = c.resize()
	}

	c.preview = c.preview.SetContent(loadNote(t, c.config.NotePath))
	c.SetFocus(previewModeShown)
	return c
}

// SetFocus sets the focus to the months or the preview window, but will not
// enable the preview if it was hidden.
func (c *Calendar) SetFocus(f previewMode) {
	if c.previewMode == previewModeHidden {
		return
	}

	if f == previewModeFocused {
		c.preview.Focus()
		for id := range c.months {
			c.months[id].Unfocus()
		}
	} else {
		c.preview.Unfocus()
		for id := range c.months {
			c.months[id].Focus()
		}
	}
	c.previewMode = f
}

// ToggleFocus switches the focus between the months and the preview window.
func (c *Calendar) ToggleFocus() {
	if c.previewMode == previewModeShown {
		c.SetFocus(previewModeFocused)
	} else {
		c.SetFocus(previewModeShown)
	}
}

// TogglePreview switches the focus between the months and the preview window.
func (c *Calendar) TogglePreview() {
	c.preview.Unfocus()
	for id := range c.months {
		c.months[id].Focus()
	}
	if c.previewMode == previewModeShown ||
		c.previewMode == previewModeFocused {
		c.previewMode = previewModeHidden
	} else {
		c.previewMode = previewModeShown
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

// View renders the calendar in its current state.
func (c Calendar) View() string {
	if !c.initialized {
		return ""
	}

	// Build a slice of rendered months.
	// The slice begins with an empty string to add a blank line of padding at
	// the top of the window. This is to make it line up with the preview
	// window, which contains a blank line or a a border at the top (depending
	// on focus).
	months := []string{""}
	for _, month := range c.months {
		months = append(months, month.View())
	}

	var preview string
	if c.previewMode != previewModeHidden {
		preview = c.preview.View()
	}
	r := c.style.Render(lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, months...),
		preview,
	))
	return r
}
