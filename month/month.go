// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package month

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.sr.ht/~kota/calendar/config"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

const MonthHeight = 8
const MonthWidth = 20

var (
	monthstyle   = lipgloss.NewStyle().Height(MonthHeight)
	headingStyle = lipgloss.NewStyle().Width(MonthWidth).Align(lipgloss.Center)
	gridStyle    = lipgloss.NewStyle().Width(MonthWidth)
)

// Model is the Bubble Tea model for this month element.
type Model struct {
	date      time.Time
	today     time.Time
	selected  time.Time
	config    *config.Config
	id        string
	showYear  bool
	isFocused bool
}

// New creates a new month model.
func New(
	date, today, selected time.Time,
	showYear bool,
	conf *config.Config,
) Model {
	return Model{
		id:       date.Format("2006-01"),
		date:     date,
		today:    today,
		selected: selected,
		showYear: showYear,
		config:   conf,
	}
}

// Init the month in Bubble Tea.
func (m Model) Init() tea.Cmd {
	return nil
}

// Updates the month in the Bubble Tea update loop.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.isFocused {
			return m, nil
		}
		switch msg.String() {
		case "h", "left":
			m.selected = m.selected.AddDate(0, 0, -1)
		case "l", "right":
			m.selected = m.selected.AddDate(0, 0, 1)
		case "j", "down":
			m.selected = m.selected.AddDate(0, 0, 7)
		case "k", "up":
			m.selected = m.selected.AddDate(0, 0, -7)
		}
	case tea.MouseMsg:
		if msg.Type != tea.MouseLeft {
			return m, nil
		}

		last := lastDay(m.date)
		for day := 1; day <= last.Day(); day++ {
			if zone.Get(m.id + "-" + strconv.Itoa(day)).InBounds(msg) {
				t, err := time.ParseInLocation(
					"2006-01",
					m.id,
					m.date.Location(),
				)
				if err != nil {
					return m, nil
				}
				year := t.Year()
				month := t.Month()

				m.selected = time.Date(year, month, day, 0, 0, 0, 0, m.date.Location())
				break
			}
		}
	}

	return m, nil
}

// Selected returns the current selection (from the perspective of this month).
func (m Model) Selected() time.Time {
	return m.selected
}

// Select updates the selected time.
func (m Model) Select(t time.Time) Model {
	m.selected = t
	return m
}

// Focus the preview.
func (m *Model) Focus() {
	m.isFocused = true
}

// Unfocus the preview.
func (m *Model) Unfocus() {
	m.isFocused = false
}

// View renders the month in its current state.
func (m Model) View() string {
	h := headingStyle.Render(m.heading())
	g := gridStyle.Render(m.grid())

	return monthstyle.Render(lipgloss.JoinVertical(lipgloss.Top, h, g))
}

// heading prints the month and optionally year centered with the weekday list
// below it.
func (m Model) heading() string {
	var heading strings.Builder
	heading.WriteString(m.date.Month().String())
	if m.showYear {
		heading.WriteString(" ")
		heading.WriteString(strconv.Itoa(m.date.Year()))
	}
	heading.WriteString("\n")
	heading.WriteString("Su Mo Tu We Th Fr Sa")

	style := headingStyle.Copy()
	if !sameMonth(m.date, m.selected) {
		style.Inherit(
			lipgloss.NewStyle().Foreground(lipgloss.Color(m.config.InactiveColor)),
		)
	}
	return style.Render(heading.String())
}

// grid prints the out the date grid for a given month.
func (m Model) grid() string {
	first := firstDay(m.date)
	last := lastDay(m.date)

	var b strings.Builder
	// Insert blank padding until first day.
	for i := 0; i < int(first.Weekday()); i++ {
		b.WriteString("   ")
	}

	// Render the grid of days.
	for i := 1; i <= last.Day(); i++ {
		day := lipgloss.NewStyle()
		if sameMonth(m.date, m.selected) {
			if i == m.selected.Day() {
				day = day.Copy().Reverse(true)
			}
		} else {
			day = day.Inherit(
				lipgloss.NewStyle().Foreground(lipgloss.Color(m.config.InactiveColor)),
			)
		}
		if sameMonth(m.date, m.today) && i == m.today.Day() {
			day = day.Copy().Foreground(lipgloss.Color(m.config.TodayColor))
		}
		b.WriteString(
			day.Render(zone.Mark(
				m.id+"-"+strconv.Itoa(i),
				fmt.Sprintf("%2.d", i),
			)),
		)
		b.WriteString(" ")
		if (i+int(first.Weekday()))%7 == 0 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

// firstDay returns a time representing the first day of the month for time t.
func firstDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// lastDay returns a time representing the last day of the month for time t.
func lastDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location())
}

// sameMonth returns true if both times are in the same month and year.
func sameMonth(x, y time.Time) bool {
	if x.Year() == y.Year() && int(x.Month()) == int(y.Month()) {
		return true
	}
	return false
}
