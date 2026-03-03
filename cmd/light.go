package cmd

import (
	"fmt"

	"github.com/LarsEckart/huey/hue"
	"github.com/spf13/cobra"
)

var (
	flagOn     bool
	flagOff    bool
	flagToggle bool
	flagName   string
)

// LightCmd controls a single light.
var LightCmd = &cobra.Command{
	Use:   "light <id>",
	Short: "Control a single light",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		lightID := args[0]

		if flagName != "" {
			return renameLight(lightID, flagName)
		}

		flagCount := 0
		if flagOn {
			flagCount++
		}
		if flagOff {
			flagCount++
		}
		if flagToggle {
			flagCount++
		}

		if flagCount == 0 {
			return showLight(lightID)
		}

		if flagCount > 1 {
			return fmt.Errorf("use only one of --on, --off, or --toggle")
		}

		client, err := authenticatedClient()
		if err != nil {
			return err
		}

		var targetOn bool
		if flagToggle {
			light, err := client.GetLight(lightID)
			if err != nil {
				return fmt.Errorf("get light: %w", err)
			}
			targetOn = !light.On
		} else {
			targetOn = flagOn
		}

		state := hue.LightState{On: &targetOn}
		if err := client.SetLightState(lightID, state); err != nil {
			return fmt.Errorf("set light state: %w", err)
		}

		status := "off"
		if targetOn {
			status = "on"
		}
		fmt.Printf("Light %s turned %s\n", lightID, status)
		return nil
	},
}

func showLight(lightID string) error {
	client, err := authenticatedClient()
	if err != nil {
		return err
	}

	light, err := client.GetLight(lightID)
	if err != nil {
		return fmt.Errorf("get light: %w", err)
	}

	status := "off"
	if light.On {
		status = "on"
	}

	fmt.Printf("ID:         %s\n", light.ID)
	fmt.Printf("Name:       %s\n", light.Name)
	fmt.Printf("Type:       %s\n", light.Type)
	fmt.Printf("State:      %s\n", status)
	fmt.Printf("Brightness: %d\n", light.Brightness)
	return nil
}

func renameLight(lightID, name string) error {
	client, err := authenticatedClient()
	if err != nil {
		return err
	}

	if err := client.RenameLight(lightID, name); err != nil {
		return fmt.Errorf("rename light: %w", err)
	}

	fmt.Printf("Light %s renamed to %q\n", lightID, name)
	return nil
}

func init() {
	LightCmd.Flags().BoolVar(&flagOn, "on", false, "Turn light on")
	LightCmd.Flags().BoolVar(&flagOff, "off", false, "Turn light off")
	LightCmd.Flags().BoolVar(&flagToggle, "toggle", false, "Toggle light state")
	LightCmd.Flags().StringVar(&flagName, "name", "", "Rename the light")
}
