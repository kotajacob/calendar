// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package month

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.sr.ht/~kota/calendar/config"
	"git.sr.ht/~kota/calendar/date"
	"git.sr.ht/~kota/calendar/holiday"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

const (
	MonthHeight = 8
	MonthWidth  = 20
)

var (
	monthstyle   = lipgloss.NewStyle().Height(MonthHeight)
	headingStyle = lipgloss.NewStyle().Width(MonthWidth).Align(lipgloss.Center)
	gridStyle    = lipgloss.NewStyle().Width(MonthWidth)
)

// Layout describes the arrangement of the month elements.
type Layout uint8

const (
	// LayoutColumn is when a column of month elements should be shown.
	LayoutColumn Layout = iota
	// LayoutGrid is when a grid of a full year of months should be shown.
	LayoutGrid
)

// Month is the Bubble Tea model for this month element.
type Month struct {
	date      time.Time
	today     time.Time
	selected  time.Time
	holidays  holiday.Holidays
	config    *config.Config
	id        string
	layout    Layout
	isFocused bool
}

// New creates a new month model.
func New(
	date, today, selected time.Time,
	layout Layout,
	holidays holiday.Holidays,
	conf *config.Config,
) Month {
	return Month{
		id:       date.Format("2006-01"),
		date:     date,
		today:    today,
		selected: selected,
		layout:   layout,
		holidays: holidays,
		config:   conf,
	}
}

// Init the month in Bubble Tea.
func (m Month) Init() tea.Cmd {
	return nil
}

// Updates the month in the Bubble Tea update loop.
func (m Month) Update(msg tea.Msg) (Month, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.isFocused {
			return m, nil
		}
		m.move(msg)
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			if m.isFocused {
				m.selected = m.selected.AddDate(0, 0, -7)
			}
		case tea.MouseWheelDown:
			if m.isFocused {
				m.selected = m.selected.AddDate(0, 0, 7)
			}
		case tea.MouseLeft:
			last := date.LastDay(m.date)
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
					m.selected = time.Date(
						year, month, day, 0, 0, 0, 0,
						m.date.Location(),
					)
					break
				}
			}
		}
	}

	return m, nil
}

// Date returns this month's date.
func (m Month) Date() time.Time {
	return m.date
}

// Selected returns the current selection (from the perspective of this month).
func (m Month) Selected() time.Time {
	return m.selected
}

// Select updates the selected time.
func (m Month) Select(t time.Time) Month {
	m.selected = t
	return m
}

// Focus the preview.
func (m *Month) Focus() {
	m.isFocused = true
}

// Unfocus the preview.
func (m *Month) Unfocus() {
	m.isFocused = false
}

// SetToday sets the today value to a new time.
func (m *Month) SetToday(t time.Time) {
	m.today = t
}

// View renders the month in its current state.
func (m Month) View() string {
	h := headingStyle.Render(m.heading())
	g := gridStyle.Render(m.grid())

	return monthstyle.Render(lipgloss.JoinVertical(lipgloss.Top, h, g))
}

// heading prints the month and optionally year centered with the weekday list
// below it.
func (m Month) heading() string {
	var heading strings.Builder
	heading.WriteString(m.date.Month().String())
	if m.layout == LayoutColumn {
		heading.WriteString(" ")
		heading.WriteString(strconv.Itoa(m.date.Year()))
	}
	heading.WriteString("\n")
	heading.WriteString("Su Mo Tu We Th Fr Sa")

	style := headingStyle.Copy()
	if !date.SameMonth(m.date, m.selected) {
		style.Inherit(
			m.config.InactiveStyle.Export(lipgloss.NewStyle()),
		)
	}
	return style.Render(heading.String())
}

// grid prints the out the date grid for a given month.
func (m Month) grid() string {
	first := date.FirstDay(m.date)
	last := date.LastDay(m.date)

	var b strings.Builder
	// Insert blank padding until first day.
	for i := 0; i < int(first.Weekday()); i++ {
		b.WriteString("   ")
	}

	// Render the grid of days.
	for i := 1; i <= last.Day(); i++ {
		day := lipgloss.NewStyle()
		// Selected or inactive.
		if date.SameMonth(m.date, m.selected) {
			if i == m.selected.Day() {
				day = day.Copy().Reverse(true)
			}
		} else {
			day = day.Inherit(
				m.config.InactiveStyle.Export(lipgloss.NewStyle()),
			)
		}
		// Render holidays.
		if h, ok := m.holidays.Match(time.Date(
			m.date.Year(), m.date.Month(), i, 0, 0, 0, 0,
			m.date.Location(),
		)); ok {
			day = day.Copy().Foreground(lipgloss.Color(h.Color))
		}
		// Render today.
		if date.SameMonth(m.date, m.today) && i == m.today.Day() {
			day = m.config.TodayStyle.Export(day.Copy())
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

// String prints out the month's data for debugging.
func (m Month) String() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintln("date:", m.date))
	b.WriteString(fmt.Sprintln("today:", m.today))
	b.WriteString(fmt.Sprintln("selected:", m.selected))
	b.WriteString(fmt.Sprintln("id:", m.id))
	b.WriteString(fmt.Sprintln("layout:", m.layout))
	b.WriteString(fmt.Sprintln("is focused:", m.isFocused))
	return b.String()
}
