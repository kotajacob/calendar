// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package main

import (
	"log"
	"os"
	"time"

	"git.sr.ht/~kota/calendar/month"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type model struct {
	currentMonth month.Month
}

func newModel(currentMonth month.Month) model {
	return model{
		currentMonth: currentMonth,
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
	termWidth, termHeight, _ := term.GetSize(int(os.Stdout.Fd()))

	// Render a calendar for the current month.
	s := m.currentMonth.View()
	return lipgloss.Place(
		termWidth,
		termHeight,
		lipgloss.Center,
		lipgloss.Center,
		s,
	)
}

func main() {
	p := tea.NewProgram(
		newModel(month.NewMonth(time.Now(), true)),
		tea.WithAltScreen(),
	)
	if err := p.Start(); err != nil {
		log.Fatalf("calendar has crashed: %v\n", err)
	}
}
