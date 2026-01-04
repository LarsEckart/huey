package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Config holds the Hue bridge connection settings.
type Config struct {
	BridgeIP string `json:"bridge_ip"`
	Username string `json:"username"`
}

// Path returns the config file path: ~/.config/huey/config.json
func Path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "huey", "config.json"), nil
}

// Load reads the config from disk.
// Returns empty Config (not error) if file doesn't exist.
func Load() (*Config, error) {
	path, err := Path()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{}, nil
		}
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// Save writes the config to disk, creating directories as needed.
func (config *Config) Save() error {
	path, err := Path()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

// IsConfigured returns true if both bridge IP and username are set.
func (config *Config) IsConfigured() bool {
	return config.BridgeIP != "" && config.Username != ""
}
