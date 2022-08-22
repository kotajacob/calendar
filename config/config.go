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
	Today    string
	Inactive string
}

// Load a configuration file from the user's config directory, the system config
// directory, or as a final fallback return default config settings.
func Load() (*Config, error) {
	conf := Config{
		Today:    "2",
		Inactive: "8",
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
