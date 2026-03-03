package cmd

import (
	"fmt"

	"github.com/LarsEckart/huey/hue"
	"github.com/spf13/cobra"
)

// ScenesCmd lists all scenes.
var ScenesCmd = &cobra.Command{
	Use:   "scenes",
	Short: "List all scenes",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := authenticatedClient()
		if err != nil {
			return err
		}

		scenes, err := client.GetScenes()
		if err != nil {
			return fmt.Errorf("get scenes: %w", err)
		}

		groups, err := client.GetGroups()
		if err != nil {
			return fmt.Errorf("get groups: %w", err)
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

		return nil
	},
}
