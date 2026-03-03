package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// LightsCmd lists all lights.
var LightsCmd = &cobra.Command{
	Use:   "lights",
	Short: "List all lights",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := authenticatedClient()
		if err != nil {
			return err
		}

		lights, err := client.GetLights()
		if err != nil {
			return fmt.Errorf("get lights: %w", err)
		}

		for _, light := range lights {
			status := "off"
			if light.On {
				status = "on"
			}
			fmt.Printf("%s  %-20s  %s\n", light.ID, light.Name, status)
		}

		return nil
	},
}
