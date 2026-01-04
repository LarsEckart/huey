package cmd

import (
	"fmt"
	"os"

	"github.com/LarsEckart/huey/auth"
	"github.com/LarsEckart/huey/hue"
	"github.com/spf13/cobra"
)

// ScenesCmd lists all scenes.
var ScenesCmd = &cobra.Command{
	Use:   "scenes",
	Short: "List all scenes",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := auth.EnsureAuthenticated()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		client := hue.NewClient(cfg.BridgeIP, cfg.Username)

		scenes, err := client.GetScenes()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Get groups to show group names
		groups, err := client.GetGroups()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		groupByID := make(map[string]hue.Group)
		for _, g := range groups {
			groupByID[g.ID] = g
		}

		for _, scene := range scenes {
			groupName := "(no group)"
			if g, ok := groupByID[scene.Group]; ok {
				groupName = g.Name
			}
			fmt.Printf("%-20s %-24s [%s]\n", scene.ID, scene.Name, groupName)
		}
	},
}
