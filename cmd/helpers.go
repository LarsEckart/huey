package cmd

import (
	"fmt"

	"github.com/LarsEckart/huey/auth"
	"github.com/LarsEckart/huey/hue"
)

func authenticatedClient() (*hue.Client, error) {
	cfg, err := auth.EnsureAuthenticated()
	if err != nil {
		return nil, fmt.Errorf("ensure authentication: %w", err)
	}

	return hue.NewClient(cfg.BridgeIP, cfg.Username), nil
}
