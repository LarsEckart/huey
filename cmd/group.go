package cmd

import (
	"fmt"
	"os"

	"github.com/LarsEckart/huey/auth"
	"github.com/LarsEckart/huey/hue"
	"github.com/spf13/cobra"
)

var (
	groupFlagOn     bool
	groupFlagOff    bool
	groupFlagToggle bool
)

// GroupCmd controls a single group.
var GroupCmd = &cobra.Command{
	Use:   "group <id>",
	Short: "Control a single group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		groupID := args[0]

		// Validate flags - exactly one must be set
		flagCount := 0
		if groupFlagOn {
			flagCount++
		}
		if groupFlagOff {
			flagCount++
		}
		if groupFlagToggle {
			flagCount++
		}

		if flagCount == 0 {
			// No flag: show group info
			showGroup(groupID)
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
		if groupFlagToggle {
			groups, err := client.GetGroups()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			// Find the group to check current state
			var found bool
			for _, g := range groups {
				if g.ID == groupID {
					// Toggle: if any light is on, turn all off; otherwise turn all on
					targetOn = !g.AnyOn
					found = true
					break
				}
			}
			if !found {
				fmt.Fprintf(os.Stderr, "Error: group %s not found\n", groupID)
				os.Exit(1)
			}
		} else {
			targetOn = groupFlagOn
		}

		// Set the state
		action := hue.GroupAction{On: &targetOn}
		if err := client.SetGroupState(groupID, action); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		status := "off"
		if targetOn {
			status = "on"
		}
		fmt.Printf("Group %s turned %s\n", groupID, status)
	},
}

func showGroup(groupID string) {
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

	// Find the group
	var group *hue.Group
	for _, g := range groups {
		if g.ID == groupID {
			group = &g
			break
		}
	}

	if group == nil {
		fmt.Fprintf(os.Stderr, "Error: group %s not found\n", groupID)
		os.Exit(1)
	}

	var status string
	if group.AllOn {
		status = "all on"
	} else if group.AnyOn {
		status = "some on"
	} else {
		status = "all off"
	}

	fmt.Printf("ID:     %s\n", group.ID)
	fmt.Printf("Name:   %s\n", group.Name)
	fmt.Printf("Type:   %s\n", group.Type)
	fmt.Printf("State:  %s\n", status)
	fmt.Printf("Lights: %v\n", group.Lights)
}

func init() {
	GroupCmd.Flags().BoolVar(&groupFlagOn, "on", false, "Turn all lights in group on")
	GroupCmd.Flags().BoolVar(&groupFlagOff, "off", false, "Turn all lights in group off")
	GroupCmd.Flags().BoolVar(&groupFlagToggle, "toggle", false, "Toggle group state")
}
