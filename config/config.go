// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package config

import (
	"errors"
	"io/fs"

	"github.com/BurntSushi/toml"
	gap "github.com/muesli/go-app-paths"
)

// Config represents the toml configuration file.
type Config struct {
	TodayColor        string
	InactiveColor     string
	NotePath          string
	LeftPadding       int
	RightPadding      int
	PreviewLeftMargin int
	PreviewPadding    int
	PreviewMinWidth   int
	PreviewMaxWidth   int
}

// Load a configuration file from the user's config directory, the system config
// directory, or as a final fallback return default config settings.
func Load() (*Config, error) {
	conf := Config{
		TodayColor:        "2",
		InactiveColor:     "8",
		LeftPadding:       2,
		RightPadding:      1,
		NotePath:          "$HOME/.local/share/calendar/2006-01-02.md",
		PreviewLeftMargin: 3,
		PreviewPadding:    1,
		PreviewMinWidth:   40,
		PreviewMaxWidth:   80,
	}

	scope := gap.NewScope(gap.User, "calendar")
	configPath, err := scope.ConfigPath("config.toml")
	if err != nil {
		return nil, err
	}

	_, err = toml.DecodeFile(configPath, &conf)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	return &conf, nil
}
