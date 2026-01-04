package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/LarsEckart/huey/auth"
	"github.com/LarsEckart/huey/hue"
	"github.com/spf13/cobra"
)

// LightsCmd lists all lights.
var LightsCmd = &cobra.Command{
	Use:   "lights",
	Short: "List all lights",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := auth.EnsureAuthenticated()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		client := hue.NewClient(cfg.BridgeIP, cfg.Username)
		lights, err := client.GetLights()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Sort by ID for consistent output
		sort.Slice(lights, func(i, j int) bool {
			return lights[i].ID < lights[j].ID
		})

		for _, light := range lights {
			status := "off"
			if light.On {
				status = "on"
			}
			fmt.Printf("%s  %-20s  %s\n", light.ID, light.Name, status)
		}
	},
}
