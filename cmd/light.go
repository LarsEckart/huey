package cmd

import (
	"fmt"
	"os"

	"github.com/LarsEckart/huey/auth"
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
	Run: func(cmd *cobra.Command, args []string) {
		lightID := args[0]

		// Handle rename separately
		if flagName != "" {
			renameLight(lightID, flagName)
			return
		}

		// Validate flags - exactly one must be set
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
			// No flag: show light info
			showLight(lightID)
			return
		}

		if flagCount > 1 {
			fmt.Fprintln(os.Stderr, "Error: use only one of --on, --off, or --toggle")
			os.Exit(1)
		}

		cfg, err := auth.EnsureAuthenticated()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		client := hue.NewClient(cfg.BridgeIP, cfg.Username)

		// Determine target state
		var targetOn bool
		if flagToggle {
			light, err := client.GetLight(lightID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			targetOn = !light.On
		} else {
			targetOn = flagOn
		}

		// Set the state
		state := hue.LightState{On: &targetOn}
		if err := client.SetLightState(lightID, state); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		status := "off"
		if targetOn {
			status = "on"
		}
		fmt.Printf("Light %s turned %s\n", lightID, status)
	},
}

func showLight(lightID string) {
	cfg, err := auth.EnsureAuthenticated()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	client := hue.NewClient(cfg.BridgeIP, cfg.Username)
	light, err := client.GetLight(lightID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
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
}

func renameLight(lightID, name string) {
	cfg, err := auth.EnsureAuthenticated()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	client := hue.NewClient(cfg.BridgeIP, cfg.Username)
	if err := client.RenameLight(lightID, name); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Light %s renamed to %q\n", lightID, name)
}

func init() {
	LightCmd.Flags().BoolVar(&flagOn, "on", false, "Turn light on")
	LightCmd.Flags().BoolVar(&flagOff, "off", false, "Turn light off")
	LightCmd.Flags().BoolVar(&flagToggle, "toggle", false, "Toggle light state")
	LightCmd.Flags().StringVar(&flagName, "name", "", "Rename the light")
}
