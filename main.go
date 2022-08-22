// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package main

import (
	"log"

	"git.sr.ht/~kota/calendar/calendar"
	"git.sr.ht/~kota/calendar/config"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

// model is the top level Bubble Tea model for the whole program.
type model struct {
	width  int
	height int

	calendar calendar.Model
	config   *config.Config
}

// Init the model in Bubble Tea.
func (m model) Init() tea.Cmd {
	return nil
}

// Updates the model in the Bubble Tea update loop.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m.propagate(msg)
}

// propagate an update to all children.
func (m model) propagate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var c tea.Cmd
	m.calendar, c = m.calendar.Update(msg)
	return m, c
}

// View renders the model in its current state.
func (m model) View() string {
	// Render a calendar for the current month.
	return zone.Scan(lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		m.calendar.View(),
	))
}

func main() {
	log.SetPrefix("")
	log.SetFlags(0)

	conf, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v\n", err)
	}

	zone.NewGlobal()
	p := tea.NewProgram(
		model{
			calendar: calendar.New(conf),
			config:   conf,
		},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if err := p.Start(); err != nil {
		log.Fatalf("calendar has crashed: %v\n", err)
	}
}
