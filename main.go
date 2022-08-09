// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package main

import (
	"log"
	"time"

	"git.sr.ht/~kota/calendar/month"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

type model struct {
	width        int
	height       int
	currentMonth tea.Model
}

func newModel() model {
	now := time.Now()
	return model{
		currentMonth: month.NewMonth(now, now, now, true),
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
			return m.propagate(msg)
		case "l", "right":
			return m.propagate(msg)
		case "j", "down":
			return m.propagate(msg)
		case "k", "up":
			return m.propagate(msg)
		}
	case tea.MouseMsg:
		return m.propagate(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m *model) propagate(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Propagate to all children.
	var c tea.Cmd
	m.currentMonth, c = m.currentMonth.Update(msg)
	return m, c
}

func (m model) View() string {
	// Render a calendar for the current month.
	return zone.Scan(lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		m.currentMonth.View(),
	))
}

func main() {
	zone.NewGlobal()
	p := tea.NewProgram(
		newModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if err := p.Start(); err != nil {
		log.Fatalf("calendar has crashed: %v\n", err)
	}
}
