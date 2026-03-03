package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var sceneFlagDelete bool

// SceneCmd activates a single scene.
var SceneCmd = &cobra.Command{
	Use:   "scene <id>",
	Short: "Activate or delete a scene",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sceneID := args[0]

		client, err := authenticatedClient()
		if err != nil {
			return err
		}

		scene, err := client.GetScene(sceneID)
		if err != nil {
			return fmt.Errorf("get scene: %w", err)
		}

		if sceneFlagDelete {
			if err := client.DeleteScene(sceneID); err != nil {
				return fmt.Errorf("delete scene: %w", err)
			}
			fmt.Printf("Deleted scene %q\n", scene.Name)
			return nil
		}

		if err := client.ActivateScene(sceneID); err != nil {
			return fmt.Errorf("activate scene: %w", err)
		}

		fmt.Printf("Activated scene %q\n", scene.Name)
		return nil
	},
}

func init() {
	SceneCmd.Flags().BoolVar(&sceneFlagDelete, "delete", false, "Delete the scene")
}
