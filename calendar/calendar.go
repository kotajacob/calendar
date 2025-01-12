// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package calendar

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"git.sr.ht/~kota/calendar/config"
	"git.sr.ht/~kota/calendar/date"
	"git.sr.ht/~kota/calendar/holiday"
	"git.sr.ht/~kota/calendar/keyword"
	"git.sr.ht/~kota/calendar/month"
	"git.sr.ht/~kota/calendar/note"
	"git.sr.ht/~kota/calendar/preview"
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// previewMode describes if the preview is shown, focused, or hidden.
type previewMode uint8

const (
	// previewModeShown is when the preview window is displayed, but not
	// currently focused.
	previewModeShown previewMode = iota
	// previewModeFocused is when the preview window is both displayed and
	// currently focused.
	previewModeFocused
	// previewModeHidden is when the preview window is hidden and thus cannot
	// be focused.
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
	holidays    holiday.Holidays
	keywords    keyword.Keywords
	height      int
	width       int
	initialized bool
}

// New creates a new calendar model.
func New(selected time.Time, conf *config.Config) Calendar {
	now := time.Now()
	holidays := holiday.Load(conf.HolidayLists)
	m := Calendar{
		today:    now,
		selected: selected,
		style: lipgloss.NewStyle().
			PaddingLeft(conf.LeftPadding).
			PaddingRight(conf.RightPadding),
		months: []month.Month{
			month.New(
				selected,
				now,
				selected,
				month.LayoutColumn,
				holidays,
				conf,
			),
		},
		holidays: holidays,
		config:   conf,
	}
	m.SetFocus(previewModeShown)
	return m
}

// Init the calendar in Bubble Tea.
func (c Calendar) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, m := range c.months {
		cmds = append(cmds, m.Init())
	}
	return tea.Batch(cmds...)
}

// Update the calendar in the Bubble Tea update loop.
func (c Calendar) Update(msg tea.Msg) (Calendar, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case c.config.KeySelectLeft.Contains(msg.String()) ||
			c.config.KeySelectRight.Contains(msg.String()) ||
			c.config.KeyLastSunday.Contains(msg.String()) ||
			c.config.KeyNextSunday.Contains(msg.String()) ||
			c.config.KeyNextSaturday.Contains(msg.String()):
			if c.previewMode != previewModeHidden {
				c.SetFocus(previewModeShown)
			}
		case c.config.KeyEditNote.Contains(msg.String()):
			path := filepath.Join(os.ExpandEnv(c.config.NoteDir), c.selected.Format("2006-01-02")) + ".md"
			cmd := tea.ExecProcess(
				exec.Command(c.config.Editor, path),
				func(err error) tea.Msg {
					return editorFinishedMsg{err: err}
				})
			return c, cmd
		case c.config.KeyFocusPreview.Contains(msg.String()):
			c.ToggleFocus()
		case c.config.KeyTogglePreview.Contains(msg.String()):
			c.TogglePreview()
			var cmd tea.Cmd
			c, cmd = c.resize()
			cmds = append(cmds, cmd)
		case c.config.KeyYankDate.Contains(msg.String()):
			clipboard.WriteAll(c.selected.Format("2006-01-02"))
		}
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height

		if !c.initialized {
			note := note.Load(c.selected, c.config.NoteDir)
			note = c.holidays.Prefix(c.selected, note)
			c.preview = preview.New(note, c.config)
			c.initialized = true
		}

		var cmd tea.Cmd
		c, cmd = c.resize()
		cmds = append(cmds, cmd)
	case editorFinishedMsg:
		// Reload the note when the user exits their editor.
		var cmd tea.Cmd
		c, cmd = c.Select(c.selected)
		cmds = append(cmds, cmd)
	}

	var cmd tea.Cmd
	c, cmd = c.propagate(msg)
	cmds = append(cmds, cmd)
	return c, tea.Batch(cmds...)
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
			var cmd tea.Cmd
			c, cmd = c.Select(month.Selected())
			cmds = append(cmds, cmd)
		}
	}

	var cmd tea.Cmd
	c.preview, cmd = c.preview.Update(msg)
	cmds = append(cmds, cmd)

	return c, tea.Batch(cmds...)
}

// Select a different date. This updates the selection on the calendar, all of
// it's months, and the preview window. It also sets the focus to the months.
func (c Calendar) Select(t time.Time) (Calendar, tea.Cmd) {
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

	var cmd tea.Cmd
	if offScreen {
		c, cmd = c.resize()
	}

	note := note.Load(t, c.config.NoteDir)
	note = c.holidays.Prefix(t, note)
	c.preview = c.preview.SetContent(note)
	if c.previewMode != previewModeHidden {
		c.SetFocus(previewModeShown)
	}
	return c, cmd
}

// SetFocus sets the focus to the months or the preview window.
func (c *Calendar) SetFocus(f previewMode) {
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

// renderMonths displays a grid of months.
func (c Calendar) renderMonths() string {
	var rows []string
	switch len(c.months) {
	case 12:
		banner := c.selected.Format("2006")
		if h, ok := c.holidays.Match(c.selected); ok {
			banner = h.Message + " " + banner
		}
		rows = append(rows, banner)
		for i := 0; i < 3; i++ {
			var column []string
			column = append(column, c.months[0+i*4].View()+strings.Repeat(" ",
				c.config.LeftPadding))
			column = append(column, c.months[1+i*4].View()+strings.Repeat(" ",
				c.config.LeftPadding))
			column = append(column, c.months[2+i*4].View()+strings.Repeat(" ",
				c.config.LeftPadding))
			column = append(column, c.months[3+i*4].View())

			rows = append(rows, lipgloss.JoinHorizontal(
				lipgloss.Center,
				column...,
			))
		}
	case 3:
		// The slice begins with an empty string to add a blank line of padding
		// at the top of the window. This is to make it line up with the
		// preview window, which contains a blank line or border at the top
		// (depending on focus).
		rows = append(rows, "")
		for _, month := range c.months {
			rows = append(rows, month.View())
		}
	case 1:
		rows = append(rows, c.months[0].View())
	}

	return lipgloss.JoinVertical(lipgloss.Center, rows...)
}

// renderPreview displays the preview window or returns a blank string.
func (c Calendar) renderPreview() string {
	if c.previewMode != previewModeHidden {
		return c.preview.View()
	}
	return ""
}

// View renders the calendar in its current state.
func (c Calendar) View() string {
	if !c.initialized {
		return ""
	}

	return c.style.Render(lipgloss.JoinHorizontal(
		lipgloss.Center,
		c.renderMonths(),
		c.renderPreview(),
	))
}
