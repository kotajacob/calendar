// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package main

import (
	"log"
	"time"

	"git.sr.ht/~kota/calendar/month"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	today        time.Time
	selected     time.Time
	width        int
	height       int
	currentMonth month.Month
}

func newModel() model {
	now := time.Now()
	return model{
		today:    now,
		selected: now,
		currentMonth: month.Month{
			Date:     now,
			Today:    now,
			Selected: now,
			ShowYear: true,
		},
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
		case "h", "left":
			m.selected = m.selected.AddDate(0, 0, -1)
			m.currentMonth.Selected = m.selected
		case "l", "right":
			m.selected = m.selected.AddDate(0, 0, 1)
			m.currentMonth.Selected = m.selected
		case "j", "down":
			m.selected = m.selected.AddDate(0, 0, 7)
			m.currentMonth.Selected = m.selected
		case "k", "up":
			m.selected = m.selected.AddDate(0, 0, -7)
			m.currentMonth.Selected = m.selected
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m model) View() string {
	// Render a calendar for the current month.
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		m.currentMonth.View(),
	)
}

func main() {
	p := tea.NewProgram(
		newModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if err := p.Start(); err != nil {
		log.Fatalf("calendar has crashed: %v\n", err)
	}
}
