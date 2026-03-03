package cmd

import (
	"fmt"
	"slices"
	"strings"

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
	RunE: func(cmd *cobra.Command, args []string) error {
		if createGroupName == "" {
			return fmt.Errorf("--name is required")
		}

		groupType := cases.Title(language.English).String(strings.ToLower(createGroupType))
		if !slices.Contains([]string{"Room", "Zone"}, groupType) {
			return fmt.Errorf("--type must be 'room' or 'zone'")
		}

		var lightIDs []string
		if createGroupLights != "" {
			for lightID := range strings.SplitSeq(createGroupLights, ",") {
				lightIDs = append(lightIDs, strings.TrimSpace(lightID))
			}
		}

		client, err := authenticatedClient()
		if err != nil {
			return err
		}

		id, err := client.CreateGroup(createGroupName, groupType, lightIDs)
		if err != nil {
			return fmt.Errorf("create group: %w", err)
		}

		fmt.Printf("Created %s %q (ID: %s)\n", groupType, createGroupName, id)
		return nil
	},
}

func init() {
	GroupCreateCmd.Flags().StringVar(&createGroupName, "name", "", "Name of the group (required)")
	GroupCreateCmd.Flags().StringVar(&createGroupType, "type", "zone", "Type: 'room' or 'zone'")
	GroupCreateCmd.Flags().StringVar(&createGroupLights, "lights", "", "Comma-separated light IDs (e.g., '1,2,3')")
}
