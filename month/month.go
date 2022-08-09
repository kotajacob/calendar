// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package month

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

var (
	headingStyle = lipgloss.NewStyle().Width(20).Align(lipgloss.Center)
	gridStyle    = lipgloss.NewStyle().Width(20)
)

type Month struct {
	id       string
	date     time.Time
	today    time.Time
	selected time.Time
	showYear bool
}

func NewMonth(date, today, selected time.Time, showYear bool) Month {
	return Month{
		id:       zone.NewPrefix(),
		date:     date,
		today:    today,
		selected: selected,
		showYear: showYear,
	}
}

func (m Month) Init() tea.Cmd {
	return nil
}

func (m Month) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
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
		for i := 1; i < last.Day(); i++ {
			if zone.Get(m.id + strconv.Itoa(i)).InBounds(msg) {
				v := i - m.selected.Day()
				m.selected = m.selected.AddDate(0, 0, v)
				break
			}
		}
	}

	return m, nil
}

func (m Month) View() string {
	h := headingStyle.Render(m.heading())
	g := gridStyle.Render(m.grid())
	return lipgloss.JoinVertical(lipgloss.Top, h, g)
}

// heading prints the month and optionally year centered with the weekday list
// below it.
func (m Month) heading() string {
	var heading strings.Builder
	heading.WriteString(m.date.Month().String())
	if m.showYear {
		heading.WriteString(" ")
		heading.WriteString(strconv.Itoa(m.date.Year()))
	}
	heading.WriteString("\n")
	heading.WriteString("Su Mo Tu We Th Fr Sa")

	return heading.String()
}

// grid prints the out the date grid for a given month.
func (m Month) grid() string {
	first := firstDay(m.date)
	last := lastDay(m.date)

	var b strings.Builder
	// Insert blank padding until first day.
	for i := 0; i < int(first.Weekday()); i++ {
		b.WriteString("   ")
	}

	// Render the grid of days.
	for i := first.Day(); i <= last.Day(); i++ {
		day := lipgloss.NewStyle()
		if i == m.today.Day() {
			day = day.Copy().Foreground(lipgloss.Color("5"))
		}
		if i == m.selected.Day() {
			day = day.Copy().Reverse(true)
		}
		b.WriteString(
			day.Render(zone.Mark(
				m.id+strconv.Itoa(i),
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
