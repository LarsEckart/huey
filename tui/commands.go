package tui

import (
	"github.com/LarsEckart/huey/hue"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) loadLights() tea.Msg {
	lights, err := m.client.GetLights()
	if err != nil {
		return errMsg{err: err}
	}
	return lightsLoadedMsg{lights: lights}
}

func (m Model) loadGroups() tea.Msg {
	groups, err := m.client.GetGroups()
	if err != nil {
		return errMsg{err: err}
	}
	return groupsLoadedMsg{groups: groups}
}

func (m Model) loadScenes() tea.Msg {
	scenes, err := m.client.GetScenes()
	if err != nil {
		return errMsg{err: err}
	}
	return scenesLoadedMsg{scenes: scenes}
}

func (m Model) activateScene(id, name string) tea.Cmd {
	return func() tea.Msg {
		if err := m.client.ActivateScene(id); err != nil {
			return errMsg{err: err}
		}
		return sceneActivatedMsg{id: id, name: name}
	}
}

func (m Model) createScene(name, groupID string) tea.Cmd {
	return func() tea.Msg {
		id, err := m.client.CreateScene(name, groupID)
		if err != nil {
			return errMsg{err: err}
		}
		return sceneCreatedMsg{
			scene: hue.Scene{
				ID:    id,
				Name:  name,
				Group: groupID,
				Type:  "GroupScene",
			},
		}
	}
}

func (m Model) deleteScene(id string) tea.Cmd {
	return func() tea.Msg {
		if err := m.client.DeleteScene(id); err != nil {
			return errMsg{err: err}
		}
		return sceneDeletedMsg{id: id}
	}
}

func (m Model) toggleLight(id string, currentOn bool) tea.Cmd {
	return func() tea.Msg {
		newOn := !currentOn
		state := hue.LightState{On: &newOn}
		if err := m.client.SetLightState(id, state); err != nil {
			return errMsg{err: err}
		}
		return lightToggledMsg{id: id, newOn: newOn}
	}
}

func (m Model) toggleGroup(id string, anyOn bool) tea.Cmd {
	return func() tea.Msg {
		// If any light is on, turn all off. Otherwise turn all on.
		newOn := !anyOn
		action := hue.GroupAction{On: &newOn}
		if err := m.client.SetGroupState(id, action); err != nil {
			return errMsg{err: err}
		}
		return groupToggledMsg{id: id, newOn: newOn}
	}
}

func (m Model) renameLight(id, name string) tea.Cmd {
	return func() tea.Msg {
		if err := m.client.RenameLight(id, name); err != nil {
			return errMsg{err: err}
		}
		return lightRenamedMsg{id: id, newName: name}
	}
}

func (m Model) renameGroup(id, name string) tea.Cmd {
	return func() tea.Msg {
		if err := m.client.RenameGroup(id, name); err != nil {
			return errMsg{err: err}
		}
		return groupRenamedMsg{id: id, newName: name}
	}
}

func (m Model) deleteGroup(id string) tea.Cmd {
	return func() tea.Msg {
		if err := m.client.DeleteGroup(id); err != nil {
			return errMsg{err: err}
		}
		return groupDeletedMsg{id: id}
	}
}

func (m Model) createGroup(name, groupType string, lightIDs []string) tea.Cmd {
	return func() tea.Msg {
		id, err := m.client.CreateGroup(name, groupType, lightIDs)
		if err != nil {
			return errMsg{err: err}
		}
		return groupCreatedMsg{
			group: hue.Group{
				ID:     id,
				Name:   name,
				Type:   groupType,
				Lights: lightIDs,
				AllOn:  false,
				AnyOn:  false,
			},
		}
	}
}
