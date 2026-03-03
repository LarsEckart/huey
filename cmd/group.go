package cmd

import (
	"fmt"

	"github.com/LarsEckart/huey/hue"
	"github.com/spf13/cobra"
)

var (
	groupFlagOn     bool
	groupFlagOff    bool
	groupFlagToggle bool
	groupFlagName   string
	groupFlagDelete bool
)

// GroupCmd controls a single group.
var GroupCmd = &cobra.Command{
	Use:   "group <id>",
	Short: "Control a single group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID := args[0]

		if groupFlagDelete {
			return deleteGroup(groupID)
		}

		if groupFlagName != "" {
			return renameGroup(groupID, groupFlagName)
		}

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
			return showGroup(groupID)
		}

		if flagCount > 1 {
			return fmt.Errorf("use only one of --on, --off, or --toggle")
		}

		client, err := authenticatedClient()
		if err != nil {
			return err
		}

		var targetOn bool
		if groupFlagToggle {
			groups, err := client.GetGroups()
			if err != nil {
				return fmt.Errorf("get groups: %w", err)
			}

			found := false
			for _, g := range groups {
				if g.ID == groupID {
					targetOn = !g.AnyOn
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("group %s not found", groupID)
			}
		} else {
			targetOn = groupFlagOn
		}

		action := hue.GroupAction{On: &targetOn}
		if err := client.SetGroupState(groupID, action); err != nil {
			return fmt.Errorf("set group state: %w", err)
		}

		status := "off"
		if targetOn {
			status = "on"
		}
		fmt.Printf("Group %s turned %s\n", groupID, status)
		return nil
	},
}

func showGroup(groupID string) error {
	client, err := authenticatedClient()
	if err != nil {
		return err
	}

	groups, err := client.GetGroups()
	if err != nil {
		return fmt.Errorf("get groups: %w", err)
	}

	var group *hue.Group
	for _, g := range groups {
		if g.ID == groupID {
			group = &g
			break
		}
	}

	if group == nil {
		return fmt.Errorf("group %s not found", groupID)
	}

	lights, err := client.GetLights()
	if err != nil {
		return fmt.Errorf("get lights: %w", err)
	}

	lightByID := make(map[string]hue.Light)
	for _, l := range lights {
		lightByID[l.ID] = l
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
	fmt.Printf("Lights:\n")
	for _, lightID := range group.Lights {
		light, ok := lightByID[lightID]
		if ok {
			state := "off"
			if light.On {
				state = "on"
			}
			fmt.Printf("  %s. %s (%s)\n", lightID, light.Name, state)
		} else {
			fmt.Printf("  %s. (unknown)\n", lightID)
		}
	}

	return nil
}

func renameGroup(groupID, name string) error {
	client, err := authenticatedClient()
	if err != nil {
		return err
	}

	if err := client.RenameGroup(groupID, name); err != nil {
		return fmt.Errorf("rename group: %w", err)
	}

	fmt.Printf("Group %s renamed to %q\n", groupID, name)
	return nil
}

func deleteGroup(groupID string) error {
	client, err := authenticatedClient()
	if err != nil {
		return err
	}

	if err := client.DeleteGroup(groupID); err != nil {
		return fmt.Errorf("delete group: %w", err)
	}

	fmt.Printf("Group %s deleted\n", groupID)
	return nil
}

func init() {
	GroupCmd.Flags().BoolVar(&groupFlagOn, "on", false, "Turn all lights in group on")
	GroupCmd.Flags().BoolVar(&groupFlagOff, "off", false, "Turn all lights in group off")
	GroupCmd.Flags().BoolVar(&groupFlagToggle, "toggle", false, "Toggle group state")
	GroupCmd.Flags().StringVar(&groupFlagName, "name", "", "Rename the group")
	GroupCmd.Flags().BoolVar(&groupFlagDelete, "delete", false, "Delete the group")
}
