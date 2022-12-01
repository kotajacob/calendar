// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package config

import (
	"errors"
	"io/fs"
	"os"

	"github.com/BurntSushi/toml"
	gap "github.com/muesli/go-app-paths"
)

// Config represents the toml configuration file.
type Config struct {
	TodayColor           string
	InactiveColor        string
	NotePath             string
	Editor               string
	LeftPadding          int
	RightPadding         int
	PreviewLeftMargin    int
	PreviewPadding       int
	PreviewMinWidth      int
	PreviewMaxWidth      int
	KeyQuit              Control
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
		TodayColor:        "2",
		InactiveColor:     "8",
		LeftPadding:       2,
		RightPadding:      1,
		NotePath:          "$HOME/.local/share/calendar/2006-01-02.md",
		Editor:            "vi",
		PreviewLeftMargin: 3,
		PreviewPadding:    1,
		PreviewMinWidth:   40,
		PreviewMaxWidth:   80,
		KeyQuit:           []string{"ctrl+c", "q"},
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
