package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// GroupsCmd lists all groups.
var GroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List all groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := authenticatedClient()
		if err != nil {
			return err
		}

		groups, err := client.GetGroups()
		if err != nil {
			return fmt.Errorf("get groups: %w", err)
		}

		for _, group := range groups {
			var status string
			if group.AllOn {
				status = "all on"
			} else if group.AnyOn {
				status = "some on"
			} else {
				status = "all off"
			}
			fmt.Printf("%s  %-20s  %-10s  %s\n", group.ID, group.Name, group.Type, status)
		}

		return nil
	},
}
