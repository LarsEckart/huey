package tui

import (
	"fmt"

	"github.com/LarsEckart/huey/hue"
)

// View renders the UI.
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	// Group info mode has its own view
	if m.mode == ModeGroupInfo {
		return m.renderGroupInfo()
	}

	// Create group modes have their own views
	if m.mode == ModeCreateGroupType {
		return m.renderCreateGroupType()
	}
	if m.mode == ModeCreateGroupName {
		return m.renderCreateGroupName()
	}
	if m.mode == ModeCreateGroupLights {
		return m.renderCreateGroupLights()
	}

	// Create scene modes have their own views
	if m.mode == ModeCreateSceneGroup {
		return m.renderCreateSceneGroup()
	}
	if m.mode == ModeCreateSceneName {
		return m.renderCreateSceneName()
	}

	s := titleStyle.Render("huey - Hue Light Control") + "\n\n"

	// Render tabs
	s += m.renderTabs() + "\n\n"

	if m.err != nil {
		s += errorStyle.Render("Cannot reach Hue bridge. Are you on the same network?") + "\n\n"
	}

	// Render active tab content
	switch m.activeTab {
	case TabLights:
		s += m.renderLights()
	case TabGroups:
		s += m.renderGroups()
	case TabScenes:
		s += m.renderScenes()
	}

	// Render delete confirmation if active
	if m.mode == ModeDeleteConfirm {
		s += "\n" + inputStyle.Render(fmt.Sprintf("Delete %q? (y/n)", m.deleteGroupName))
	}
	if m.mode == ModeDeleteSceneConfirm {
		s += "\n" + inputStyle.Render(fmt.Sprintf("Delete %q? (y/n)", m.deleteSceneName))
	}

	// Render help based on mode and tab
	switch m.mode {
	case ModeRename:
		s += "\n" + helpStyle.Render("enter confirm • esc cancel")
	case ModeDeleteConfirm, ModeDeleteSceneConfirm:
		s += "\n" + helpStyle.Render("y/enter delete • n/esc cancel")
	default:
		switch m.activeTab {
		case TabLights:
			s += "\n" + helpStyle.Render("↑/↓ navigate • space toggle • r rename • tab switch • q quit")
		case TabGroups:
			s += "\n" + helpStyle.Render("↑/↓ navigate • space toggle • a add • r rename • d delete • i info • tab switch • q quit")
		case TabScenes:
			s += "\n" + helpStyle.Render("↑/↓ navigate • space activate • a add • d delete • tab switch • q quit")
		}
	}

	return s
}

func (m Model) renderTabs() string {
	var lightsTab, groupsTab, scenesTab string

	switch m.activeTab {
	case TabLights:
		lightsTab = tabActiveStyle.Render("Lights")
		groupsTab = tabInactiveStyle.Render("Groups")
		scenesTab = tabInactiveStyle.Render("Scenes")
	case TabGroups:
		lightsTab = tabInactiveStyle.Render("Lights")
		groupsTab = tabActiveStyle.Render("Groups")
		scenesTab = tabInactiveStyle.Render("Scenes")
	case TabScenes:
		lightsTab = tabInactiveStyle.Render("Lights")
		groupsTab = tabInactiveStyle.Render("Groups")
		scenesTab = tabActiveStyle.Render("Scenes")
	}

	return lightsTab + "  " + groupsTab + "  " + scenesTab
}

func (m Model) renderLights() string {
	if !m.lightsLoaded && m.err == nil {
		return "Loading lights...\n"
	}
	if len(m.lights) == 0 {
		return "No lights found.\n"
	}

	var s string
	for i, light := range m.lights {
		cursor := "  "
		style := normalStyle
		isSelected := i == m.lightCursor
		if isSelected {
			cursor = "> "
			style = selectedStyle
		}

		var status string
		if light.On {
			status = onStyle.Render("● on")
		} else {
			status = offStyle.Render("○ off")
		}

		// Show text input if renaming this light
		var name string
		if m.mode == ModeRename && isSelected && m.renameID == light.ID {
			name = m.textInput.View()
		} else {
			name = fmt.Sprintf("%-24s", light.Name)
		}

		line := fmt.Sprintf("%s%s %s", cursor, name, status)
		s += style.Render(line) + "\n"
	}

	return s
}

func (m Model) renderGroups() string {
	if !m.groupsLoaded && m.err == nil {
		return "Loading groups...\n"
	}
	if len(m.groups) == 0 {
		return "No groups found.\n"
	}

	var s string
	for i, group := range m.groups {
		cursor := "  "
		style := normalStyle
		isSelected := i == m.groupCursor
		if isSelected {
			cursor = "> "
			style = selectedStyle
		}

		var status string
		if group.AllOn {
			status = onStyle.Render("● all on")
		} else if group.AnyOn {
			status = partialStyle.Render("◐ some on")
		} else {
			status = offStyle.Render("○ all off")
		}

		groupType := typeStyle.Render(fmt.Sprintf("(%s)", group.Type))

		// Show text input if renaming this group
		var name string
		if m.mode == ModeRename && isSelected && m.renameID == group.ID {
			name = m.textInput.View()
		} else {
			name = fmt.Sprintf("%-20s", group.Name)
		}

		line := fmt.Sprintf("%s%s %-8s %s", cursor, name, groupType, status)
		s += style.Render(line) + "\n"
	}

	return s
}

func (m Model) renderScenes() string {
	if !m.scenesLoaded && m.err == nil {
		return "Loading scenes...\n"
	}
	if len(m.scenes) == 0 {
		return "No scenes found.\n"
	}

	// Build group lookup for display
	groupByID := make(map[string]hue.Group)
	for _, g := range m.groups {
		groupByID[g.ID] = g
	}

	var s string
	for i, scene := range m.scenes {
		cursor := "  "
		style := normalStyle
		if i == m.sceneCursor {
			cursor = "> "
			style = selectedStyle
		}

		groupName := "(no group)"
		if g, ok := groupByID[scene.Group]; ok {
			groupName = g.Name
		}

		line := fmt.Sprintf("%s%-24s %s", cursor, scene.Name, typeStyle.Render(fmt.Sprintf("[%s]", groupName)))
		s += style.Render(line) + "\n"
	}

	return s
}

func (m Model) renderGroupInfo() string {
	// Find the group
	var group *hue.Group
	for i := range m.groups {
		if m.groups[i].ID == m.infoGroupID {
			group = &m.groups[i]
			break
		}
	}

	if group == nil {
		return "Group not found\n\n" + helpStyle.Render("esc back")
	}

	// Build light lookup
	lightByID := make(map[string]hue.Light)
	for _, l := range m.lights {
		lightByID[l.ID] = l
	}

	s := titleStyle.Render(fmt.Sprintf("Group: %s", group.Name)) + "\n"
	s += typeStyle.Render(fmt.Sprintf("Type: %s", group.Type)) + "\n\n"

	s += "Lights:\n"
	for _, lightID := range group.Lights {
		light, ok := lightByID[lightID]
		if ok {
			var status string
			if light.On {
				status = onStyle.Render("● on")
			} else {
				status = offStyle.Render("○ off")
			}
			s += fmt.Sprintf("  %-24s %s\n", light.Name, status)
		} else {
			s += "  (unknown light)\n"
		}
	}

	s += "\n" + helpStyle.Render("esc back")
	return s
}

func (m Model) renderCreateGroupType() string {
	s := titleStyle.Render("Create Group") + "\n\n"
	s += "Select type:\n\n"
	s += "  [r] Room — A light can only be in one room\n"
	s += "  [z] Zone — A light can be in multiple zones\n"
	s += "\n" + helpStyle.Render("r room • z zone • esc cancel")
	return s
}

func (m Model) renderCreateGroupName() string {
	s := titleStyle.Render(fmt.Sprintf("Create %s", m.createGroupType)) + "\n\n"
	s += "Enter name:\n\n"
	s += "  " + m.textInput.View() + "\n"
	s += "\n" + helpStyle.Render("enter confirm • esc cancel")
	return s
}

func (m Model) renderCreateGroupLights() string {
	s := titleStyle.Render(fmt.Sprintf("Create %s: %s", m.createGroupType, m.createGroupName)) + "\n\n"
	s += "Select lights (space to toggle):\n\n"

	for i, light := range m.lights {
		cursor := "  "
		if i == m.createLightCursor {
			cursor = "> "
		}

		checkbox := "[ ]"
		if m.createLightSelected[light.ID] {
			checkbox = "[✓]"
		}

		line := fmt.Sprintf("%s%s %s %s", cursor, checkbox, light.ID, light.Name)
		if i == m.createLightCursor {
			s += selectedStyle.Render(line) + "\n"
		} else {
			s += normalStyle.Render(line) + "\n"
		}
	}

	s += "\n" + helpStyle.Render("↑/↓ navigate • space toggle • enter create • esc cancel")
	return s
}

func (m Model) renderCreateSceneGroup() string {
	s := titleStyle.Render("Create Scene") + "\n\n"
	s += "Select group to capture:\n\n"

	for i, group := range m.groups {
		cursor := "  "
		style := normalStyle
		if i == m.createGroupCursor {
			cursor = "> "
			style = selectedStyle
		}

		line := fmt.Sprintf("%s%s", cursor, group.Name)
		s += style.Render(line) + "\n"
	}

	s += "\n" + helpStyle.Render("↑/↓ navigate • enter select • esc cancel")
	return s
}

func (m Model) renderCreateSceneName() string {
	// Find the group name for display
	groupName := m.createSceneGroupID
	for _, g := range m.groups {
		if g.ID == m.createSceneGroupID {
			groupName = g.Name
			break
		}
	}

	s := titleStyle.Render(fmt.Sprintf("Create Scene for %s", groupName)) + "\n\n"
	s += "Enter scene name:\n\n"
	s += "  " + m.textInput.View() + "\n"
	s += "\n" + helpStyle.Render("enter confirm • esc cancel")
	return s
}
