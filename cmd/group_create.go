package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/LarsEckart/huey/auth"
	"github.com/LarsEckart/huey/hue"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	createGroupName   string
	createGroupType   string
	createGroupLights string
)

// GroupCreateCmd creates a new group.
var GroupCreateCmd = &cobra.Command{
	Use:   "group-create",
	Short: "Create a new group (room or zone)",
	Run: func(cmd *cobra.Command, args []string) {
		if createGroupName == "" {
			fmt.Fprintln(os.Stderr, "Error: --name is required")
			os.Exit(1)
		}

		// Normalize and validate type
		groupType := cases.Title(language.English).String(strings.ToLower(createGroupType))
		if groupType != "Room" && groupType != "Zone" {
			fmt.Fprintln(os.Stderr, "Error: --type must be 'room' or 'zone'")
			os.Exit(1)
		}

		// Parse light IDs
		var lightIDs []string
		if createGroupLights != "" {
			lightIDs = strings.Split(createGroupLights, ",")
			// Trim whitespace
			for i := range lightIDs {
				lightIDs[i] = strings.TrimSpace(lightIDs[i])
			}
		}

		cfg, err := auth.EnsureAuthenticated()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		client := hue.NewClient(cfg.BridgeIP, cfg.Username)
		id, err := client.CreateGroup(createGroupName, groupType, lightIDs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Created %s %q (ID: %s)\n", groupType, createGroupName, id)
	},
}

func init() {
	GroupCreateCmd.Flags().StringVar(&createGroupName, "name", "", "Name of the group (required)")
	GroupCreateCmd.Flags().StringVar(&createGroupType, "type", "zone", "Type: 'room' or 'zone'")
	GroupCreateCmd.Flags().StringVar(&createGroupLights, "lights", "", "Comma-separated light IDs (e.g., '1,2,3')")
}
