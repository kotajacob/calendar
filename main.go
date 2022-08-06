// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const month_width = 20

var nowStyle = lipgloss.NewStyle().Reverse(true)

type model struct {
	now time.Time
}

func newModel() model {
	now := time.Now()

	return model{
		now: now,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	// Render a calendar for the current month.
	b.WriteString(month(m.now, true))
	b.WriteString("Su Mo Tu We Th Fr Sa\n")
	b.WriteString(grid(m.now))

	return b.String()
}

// month prints the month and optionally the year heading for a given time.
func month(t time.Time, year bool) string {
	var month strings.Builder
	month.WriteString(t.Month().String())
	if year {
		month.WriteString(" ")
		month.WriteString(strconv.Itoa(t.Year()))
	}

	left_len := (month_width - len(month.String())) / 2
	var left strings.Builder
	for i := 0; i < left_len; i++ {
		left.WriteString(" ")
	}

	right_len := month_width - (left_len + len(month.String()))
	var right strings.Builder
	for i := 0; i < right_len; i++ {
		right.WriteString(" ")
	}
	return left.String() + month.String() + right.String() + "\n"
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

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatalf("calendar has crashed: %v\n", err)
	}
}
