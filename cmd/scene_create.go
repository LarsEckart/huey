package cmd

import (
	"fmt"

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
	RunE: func(cmd *cobra.Command, args []string) error {
		if sceneCreateName == "" {
			return fmt.Errorf("--name is required")
		}
		if sceneCreateGroup == "" {
			return fmt.Errorf("--group is required")
		}

		client, err := authenticatedClient()
		if err != nil {
			return err
		}

		id, err := client.CreateScene(sceneCreateName, sceneCreateGroup)
		if err != nil {
			return fmt.Errorf("create scene: %w", err)
		}

		fmt.Printf("Created scene %q (ID: %s)\n", sceneCreateName, id)
		return nil
	},
}

func init() {
	SceneCreateCmd.Flags().StringVar(&sceneCreateName, "name", "", "Scene name (required)")
	SceneCreateCmd.Flags().StringVar(&sceneCreateGroup, "group", "", "Group ID to capture (required)")
}
