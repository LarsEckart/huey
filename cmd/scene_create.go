package cmd

import (
	"fmt"
	"os"

	"github.com/LarsEckart/huey/auth"
	"github.com/LarsEckart/huey/hue"
	"github.com/spf13/cobra"
)

var (
	sceneCreateName  string
	sceneCreateGroup string
)

// SceneCreateCmd creates a new scene.
var SceneCreateCmd = &cobra.Command{
	Use:   "scene-create",
	Short: "Create a new scene (captures current light states)",
	Long:  "Creates a scene that saves the current state of all lights in a group.",
	Run: func(cmd *cobra.Command, args []string) {
		if sceneCreateName == "" {
			fmt.Fprintln(os.Stderr, "Error: --name is required")
			os.Exit(1)
		}
		if sceneCreateGroup == "" {
			fmt.Fprintln(os.Stderr, "Error: --group is required")
			os.Exit(1)
		}

		cfg, err := auth.EnsureAuthenticated()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		client := hue.NewClient(cfg.BridgeIP, cfg.Username)

		id, err := client.CreateScene(sceneCreateName, sceneCreateGroup)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Created scene %q (ID: %s)\n", sceneCreateName, id)
	},
}

func init() {
	SceneCreateCmd.Flags().StringVar(&sceneCreateName, "name", "", "Scene name (required)")
	SceneCreateCmd.Flags().StringVar(&sceneCreateGroup, "group", "", "Group ID to capture (required)")
}
