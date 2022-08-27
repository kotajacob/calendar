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

// Model is the Bubble Tea model for this calendar element. The calendar is a
// thin wrapper for the month elements. It creates and destroys them based on
// the size of the window.
type Model struct {
	today       time.Time
	selected    time.Time
	config      *config.Config
	style       lipgloss.Style
	months      []month.Model
	preview     preview.Model
	height      int
	width       int
	focus       focus
	initialized bool
}

// New creates a new calendar model.
func New(conf *config.Config) Model {
	now := time.Now()
	m := Model{
		today:    now,
		selected: now,
		style: lipgloss.NewStyle().
			PaddingLeft(conf.LeftPadding).
			PaddingRight(conf.RightPadding),
		months: []month.Model{
			month.New(now, now, now, true, conf),
		},
		config: conf,
	}
	m.SetFocus(focusMonths)
	return m
}

// Init the calendar in Bubble Tea.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update the calendar in the Bubble Tea update loop.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "h", "left":
			m.SetFocus(focusMonths)
		case "l", "right":
			m.SetFocus(focusMonths)
		case "enter":
			path := m.selected.Format(os.ExpandEnv(m.config.NotePath))
			cmd := tea.ExecProcess(exec.Command("vim", path), func(err error) tea.Msg {
				return editorFinishedMsg{err: err}
			})
			return m, cmd
		case "tab":
			m.ToggleFocus()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.initialized {
			note := loadNote(m.selected, m.config.NotePath)
			m.preview = preview.New(note, msg.Width, msg.Height, m.config)
			m.initialized = true
		}

		m = m.resize()
	case editorFinishedMsg:
		// Reload the note when the user exits their editor.
		m = m.Select(m.selected)
	}
	return m.propagate(msg)
}

// propagate an update to all children.
func (m Model) propagate(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	for i, month := range m.months {
		var cmd tea.Cmd
		m.months[i], cmd = month.Update(msg)
		cmds = append(cmds, cmd)
	}

	for _, month := range m.months {
		// If the month's selection changed we need to update our selection and
		// the selection of all other months to avoid letting the selection get
		// out of sync!
		if !month.Selected().Equal(m.selected) {
			m = m.Select(month.Selected())
		}
	}

	var cmd tea.Cmd
	m.preview, cmd = m.preview.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// resize the number of months being displayed to fill the window size.
func (m Model) resize() Model {
	var want int
	switch {
	case m.height > 3*month.MonthHeight:
		want = 3
	default:
		want = 1
	}

	if len(m.months) == want {
		return m
	}

	// TODO: Revise this to find the "center" month once we have more than 2
	// different size options.
	switch want {
	case 3:
		last := month.New(
			lastMonth(m.selected),
			m.today,
			m.selected,
			true,
			m.config,
		)
		next := month.New(
			nextMonth(m.selected),
			m.today,
			m.selected,
			true,
			m.config,
		)
		m.months = append([]month.Model{last}, m.months...)
		m.months = append(m.months, next)
	default:
		m.months = []month.Model{m.months[1]}
	}
	return m
}

// Select a different date. This updates the selection on the calendar, all of
// it's months, and the preview window. It also sets the focus to the months.
func (m Model) Select(t time.Time) Model {
	m.selected = t
	for i, month := range m.months {
		m.months[i] = month.Select(t)
	}
	m.preview = m.preview.SetContent(loadNote(t, m.config.NotePath))
	m.SetFocus(focusMonths)
	return m
}

// ToggleFocus sets the focus to the months or the preview window.
func (m *Model) SetFocus(f focus) {
	if f == focusPreview {
		m.preview.Focus()
		for id := range m.months {
			m.months[id].Unfocus()
		}
		m.focus = focusPreview
	} else {
		m.preview.Unfocus()
		for id := range m.months {
			m.months[id].Focus()
		}
		m.focus = focusMonths
	}
}

// ToggleFocus switches the focus between the months and the preview window.
func (m *Model) ToggleFocus() {
	if m.focus == focusMonths {
		m.SetFocus(focusPreview)
	} else {
		m.SetFocus(focusMonths)
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
func (m Model) View() string {
	if !m.initialized {
		return ""
	}

	// Build a slice of rendered months.
	var months []string
	for _, month := range m.months {
		months = append(months, month.View())
	}

	return m.style.Render(lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, months...),
		m.preview.View(),
	))
}
