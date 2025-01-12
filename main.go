// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package main

import (
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"git.sr.ht/~kota/calendar/calendar"
	"git.sr.ht/~kota/calendar/config"
	"git.sr.ht/~kota/calendar/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

// Version is overwritten at build time in the Makefile.
var Version = "development version"

// tickerMsg is a tea.Msg which should be returned with the current time.
type tickerMsg time.Time

// model is the top level Bubble Tea model for the whole program.
type model struct {
	config   *config.Config
	mode     mode
	calendar calendar.Calendar
	help     help.Help
	width    int
	height   int
}

// Init the model in Bubble Tea.
func (m model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		m.calendar.Init(),
		fiveMinutes(),
	}
	return tea.Batch(cmds...)
}

// mode describes if the calendar or the help menu should be shown.
type mode uint8

const (
	modeCalendar mode = iota
	modeHelp
)

// fiveMinutes starts a timer which will return a tickerMsg after five minutes.
// This is used to update the "today" value on the calendar.
func fiveMinutes() tea.Cmd {
	now := time.Now()
	tomorrow := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute()+5,
		0,
		0,
		now.Location(),
	)
	d := tomorrow.Sub(now)

	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickerMsg(t)
	})
}

// Updates the model in the Bubble Tea update loop.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case m.config.KeyHelp.Contains(msg.String()):
			m.mode = modeHelp
		case m.config.KeyQuit.Contains(msg.String()):
			if m.mode == modeHelp {
				m.mode = modeCalendar
			} else {
				return m, tea.Quit
			}

		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tickerMsg:
		// Update the "today" value and kick off another timer.
		m.calendar.SetToday(time.Now())
		return m, fiveMinutes()
	}
	return m.propagate(msg)
}

// propagate an update to all children.
func (m model) propagate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var c tea.Cmd
	if m.mode == modeHelp {
		m.help, c = m.help.Update(msg)
	} else {
		m.calendar, c = m.calendar.Update(msg)
	}
	return m, c
}

// View renders the model in its current state.
func (m model) View() string {
	if m.mode == modeHelp {
		// Display the help menu.
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			m.help.View(),
		)
	}

	// Display the calendar.
	return zone.Scan(lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		m.calendar.View(),
	))
}

// parseArgs reads the program's arguments to parse a starting selected time.
func parseArgs(args []string, now time.Time) time.Time {
	switch len(args) {
	case 1:
		// Early exit if there were no arguments.
		return now
	case 2:
		// Argument is either "day", "timestamp", or "monthname".
		timestamp, err := time.Parse("2006-01-02", args[1])
		if err == nil {
			return timestamp
		}
		monthnameTime, err := time.Parse("January", args[1])
		if err == nil {
			return time.Date(
				now.Year(),
				monthnameTime.Month(),
				now.Day(),
				0, 0, 0, 0,
				now.Location(),
			)
		}
		day, err := strconv.Atoi(args[1])
		if err == nil {
			return time.Date(
				now.Year(),
				now.Month(),
				day,
				0, 0, 0, 0,
				now.Location(),
			)
		}
		return now
	case 3:
		day, err := strconv.Atoi(args[1])
		if err != nil {
			return now
		}
		month, err := strconv.Atoi(args[2])
		if err != nil {
			return now
		}
		return time.Date(
			now.Year(),
			time.Month(month),
			day,
			0, 0, 0, 0,
			now.Location(),
		)
	case 4:
		day, err := strconv.Atoi(args[1])
		if err != nil {
			return now
		}
		month, err := strconv.Atoi(args[2])
		if err != nil {
			return now
		}
		year, err := strconv.Atoi(args[3])
		if err != nil {
			return now
		}
		return time.Date(
			year,
			time.Month(month),
			day,
			0, 0, 0, 0,
			now.Location(),
		)
	}

	// Return now as a fallback.
	return now
}

func main() {
	log.SetPrefix("")
	log.SetFlags(0)
	log.SetOutput(io.Discard)
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			log.Fatalf("failed setting up debug logging: %v\n", err)
		}
		defer f.Close()
	}

	conf, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v\n", err)
	}

	selected := parseArgs(os.Args, time.Now())
	zone.NewGlobal()
	p := tea.NewProgram(
		model{
			calendar: calendar.New(selected, conf),
			help:     help.New(Version),
			config:   conf,
		},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		log.Fatalf("calendar has crashed: %v\n", err)
	}
}
