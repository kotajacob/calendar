// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package config

import (
	"errors"
	"io/fs"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/lipgloss"
	gap "github.com/muesli/go-app-paths"
)

// Config represents the toml configuration file.
type Config struct {
	TodayStyle           Style
	InactiveStyle        Style
	NotedStyle           Style
	NoteDir              string
	Editor               string
	LeftPadding          int
	RightPadding         int
	PreviewLeftMargin    int
	PreviewPadding       int
	PreviewMinWidth      int
	PreviewMaxWidth      int
	KeyQuit              Control
	KeyHelp              Control
	KeySelectLeft        Control
	KeySelectDown        Control
	KeySelectUp          Control
	KeySelectRight       Control
	KeyFocusPreview      Control
	KeyTogglePreview     Control
	KeyScrollPreviewDown Control
	KeyScrollPreviewUp   Control
	KeyEditNote          Control
	KeyYankDate          Control
	KeyLastSunday        Control
	KeyNextSunday        Control
	KeyNextSaturday      Control
	KeyMonthUp           Control
	KeyMonthDown         Control
	HolidayLists         []string
}

// Style represents how a type of date should be displayed.
type Style struct {
	Color  string
	Bold   bool
	Italic bool
}

// Export takes a base lipgloss.Style and makes the changes needed based on this configured Style.
func (s Style) Export(base lipgloss.Style) lipgloss.Style {
	if s.Color != "" {
		base = base.Foreground(lipgloss.Color(s.Color))
	}
	base = base.Bold(s.Bold)
	base = base.Italic(s.Italic)
	return base
}

// Blank returns true if the Style has no color and is not bold or italicized.
func (s Style) Blank() bool {
	if s.Color != "" {
		return false
	}
	if s.Bold {
		return false
	}
	if s.Italic {
		return false
	}
	return true
}

// Control is a slice of strings representing the keys bound to a given action.
type Control []string

// Contains reports if a key is bound for this control action.
func (c Control) Contains(key string) bool {
	for _, k := range c {
		if key == k {
			return true
		}
	}
	return false
}

// Default returns the default configuration.
func Default() *Config {
	return &Config{
		TodayStyle:        Style{Color: "2"},
		InactiveStyle:     Style{Color: "8"},
		NotedStyle:        Style{},
		LeftPadding:       2,
		RightPadding:      1,
		NoteDir:           "$HOME/.local/share/calendar",
		Editor:            "vi",
		PreviewLeftMargin: 3,
		PreviewPadding:    1,
		PreviewMinWidth:   40,
		PreviewMaxWidth:   80,
		KeyQuit:           []string{"ctrl+c", "q"},
		KeyHelp:           []string{"?"},
		KeySelectLeft:     []string{"left", "h"},
		KeySelectDown:     []string{"down", "j"},
		KeySelectUp:       []string{"up", "k"},
		KeySelectRight:    []string{"right", "l"},
		KeyFocusPreview:   []string{"tab"},
		KeyTogglePreview:  []string{"p"},
		KeyEditNote:       []string{"enter"},
		KeyYankDate:       []string{"y"},
		KeyLastSunday:     []string{"b", "H"},
		KeyNextSunday:     []string{"w"},
		KeyNextSaturday:   []string{"e", "L"},
		KeyMonthUp:        []string{"ctrl+u"},
		KeyMonthDown:      []string{"ctrl+d"},
		HolidayLists:      []string{""},
	}
}

// Load a configuration file from the user's config directory, the system config
// directory, or as a final fallback return default config settings.
func Load() (*Config, error) {
	conf := Default()
	if edvar, ok := os.LookupEnv("EDITOR"); ok {
		conf.Editor = edvar
	}
	if visvar, ok := os.LookupEnv("VISUAL"); ok {
		conf.Editor = visvar
	}

	scope := gap.NewScope(gap.User, "calendar")
	configPath, err := scope.ConfigPath("config.toml")
	if err != nil {
		return nil, err
	}

	_, err = toml.DecodeFile(configPath, conf)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	return conf, nil
}
