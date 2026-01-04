package cmd

import (
	"fmt"
	"os"

	"github.com/LarsEckart/huey/auth"
	"github.com/LarsEckart/huey/hue"
	"github.com/spf13/cobra"
)

// GroupsCmd lists all groups.
var GroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List all groups",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := auth.EnsureAuthenticated()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		client := hue.NewClient(cfg.BridgeIP, cfg.Username)
		groups, err := client.GetGroups()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
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
	},
}
