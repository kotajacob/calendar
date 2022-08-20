// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package calendar

import (
	"time"

	"git.sr.ht/~kota/calendar/month"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model is the Bubble Tea model for this calendar element. The calendar is a
// thin wrapper for the month elements. It creates and destroys them based on
// the size of the window.
type Model struct {
	width  int
	height int

	today    time.Time
	selected time.Time
	months   []month.Model
}

func New() Model {
	now := time.Now()
	return Model{
		today:    now,
		selected: now,
		months: []month.Model{
			month.New(now, now, now, true),
		},
	}
}

// Init the calendar in Bubble Tea.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update the calendar in the Bubble Tea update loop.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m = m.resize()
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
		last := month.New(lastMonth(m.selected), m.today, m.selected, true)
		next := month.New(nextMonth(m.selected), m.today, m.selected, true)
		m.months = append([]month.Model{last}, m.months...)
		m.months = append(m.months, next)
	default:
		m.months = []month.Model{m.months[1]}
	}
	return m
}

func (m Model) Select(t time.Time) Model {
	m.selected = t
	for i, month := range m.months {
		m.months[i] = month.Select(t)
	}
	return m
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
	var strs []string
	for _, month := range m.months {
		strs = append(strs, month.View())
	}

	return lipgloss.JoinVertical(lipgloss.Center, strs...)
}
