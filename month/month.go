// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package month

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	headingStyle = lipgloss.NewStyle().Width(20).Align(lipgloss.Center)
	gridStyle    = lipgloss.NewStyle().Width(20)
	nowStyle     = lipgloss.NewStyle().Reverse(true)
)

type Month struct {
	Date     time.Time
	ShowYear bool
}

func NewMonth(t time.Time, showYear bool) Month {
	return Month{
		Date:     t,
		ShowYear: showYear,
	}
}

func (m Month) View() string {
	h := headingStyle.Render(heading(m.Date, m.ShowYear))
	g := gridStyle.Render(grid(m.Date))
	return lipgloss.JoinVertical(lipgloss.Top, h, g)
}

// heading prints the month and optionally year centered with the weekday list
// below it.
func heading(t time.Time, showYear bool) string {
	var heading strings.Builder
	heading.WriteString(t.Month().String())
	if showYear {
		heading.WriteString(" ")
		heading.WriteString(strconv.Itoa(t.Year()))
	}
	heading.WriteString("\n")
	heading.WriteString("Su Mo Tu We Th Fr Sa")

	return heading.String()
}

// grid prints the out the date grid for a given month.
func grid(t time.Time) string {
	first_day := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	last_day := time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location())

	var b strings.Builder
	// Insert blank padding until first day.
	for i := 0; i < int(first_day.Weekday()); i++ {
		b.WriteString("   ")
	}

	// Render the grid of days.
	for i := first_day.Day(); i <= last_day.Day(); i++ {
		if i == t.Day() {
			b.WriteString(nowStyle.Render(fmt.Sprintf("%2.d", i)))
		} else {
			b.WriteString(fmt.Sprintf("%2.d", i))
		}
		b.WriteString(" ")
		if (i+int(first_day.Weekday()))%7 == 0 {
			b.WriteString("\n")
		}
	}
	return b.String()
}
