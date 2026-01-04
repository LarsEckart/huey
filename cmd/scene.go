package cmd

import (
	"fmt"
	"os"

	"github.com/LarsEckart/huey/auth"
	"github.com/LarsEckart/huey/hue"
	"github.com/spf13/cobra"
)

var sceneFlagDelete bool

// SceneCmd activates a single scene.
var SceneCmd = &cobra.Command{
	Use:   "scene <id>",
	Short: "Activate or delete a scene",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sceneID := args[0]

		cfg, err := auth.EnsureAuthenticated()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		client := hue.NewClient(cfg.BridgeIP, cfg.Username)

		// Get scene info first for confirmation message
		scene, err := client.GetScene(sceneID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if sceneFlagDelete {
			if err := client.DeleteScene(sceneID); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Deleted scene %q\n", scene.Name)
			return
		}

		if err := client.ActivateScene(sceneID); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Activated scene %q\n", scene.Name)
	},
}

func init() {
	SceneCmd.Flags().BoolVar(&sceneFlagDelete, "delete", false, "Delete the scene")
}
