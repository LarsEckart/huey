package tui

import (
	"github.com/LarsEckart/huey/hue"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Tab represents which tab is active.
type Tab int

const (
	TabLights Tab = iota
	TabGroups
	TabScenes
)

// Mode represents the current interaction mode.
type Mode int

const (
	ModeNormal Mode = iota
	ModeRename
	ModeGroupInfo
	ModeDeleteConfirm
	ModeCreateGroupType
	ModeCreateGroupName
	ModeCreateGroupLights
	ModeCreateSceneGroup
	ModeCreateSceneName
	ModeDeleteSceneConfirm
)

// Model is the Bubble Tea model for the TUI.
type Model struct {
	client       *hue.Client
	lights       []hue.Light
	groups       []hue.Group
	scenes       []hue.Scene
	lightsLoaded bool
	groupsLoaded bool
	scenesLoaded bool
	activeTab    Tab
	lightCursor  int
	groupCursor  int
	sceneCursor  int
	err          error
	quitting     bool

	// Rename mode
	mode      Mode
	textInput textinput.Model
	renameID  string // ID of item being renamed

	// Group info mode
	infoGroupID string // ID of group being viewed

	// Delete confirmation mode
	deleteGroupID   string // ID of group to delete
	deleteGroupName string // Name of group to delete (for display)

	// Create group mode
	createGroupType     string          // "Room" or "Zone"
	createGroupName     string          // Name entered by user
	createGroupLights   []string        // Selected light IDs
	createLightCursor   int             // Cursor for light picker
	createLightSelected map[string]bool // Which lights are selected

	// Create scene mode
	createSceneGroupID string // Selected group ID
	createSceneName    string // Name entered by user
	createGroupCursor  int    // Cursor for group picker

	// Delete scene mode
	deleteSceneID   string // ID of scene to delete
	deleteSceneName string // Name of scene to delete (for display)
}

// New creates a new TUI model.
func New(client *hue.Client) Model {
	ti := textinput.New()
	ti.CharLimit = 32 // Hue names limited to 32 chars
	ti.Width = 24

	return Model{
		client:    client,
		activeTab: TabLights,
		mode:      ModeNormal,
		textInput: ti,
	}
}

// Init initializes the model and loads data.
func (m Model) Init() tea.Cmd {
	return tea.Batch(m.loadLights, m.loadGroups, m.loadScenes)
}

// Update handles messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle rename mode separately
	if m.mode == ModeRename {
		return m.updateRenameMode(msg)
	}

	// Handle group info mode
	if m.mode == ModeGroupInfo {
		return m.updateGroupInfoMode(msg)
	}

	// Handle delete confirmation mode
	if m.mode == ModeDeleteConfirm {
		return m.updateDeleteConfirmMode(msg)
	}

	// Handle create group modes
	if m.mode == ModeCreateGroupType {
		return m.updateCreateGroupTypeMode(msg)
	}
	if m.mode == ModeCreateGroupName {
		return m.updateCreateGroupNameMode(msg)
	}
	if m.mode == ModeCreateGroupLights {
		return m.updateCreateGroupLightsMode(msg)
	}

	// Handle create scene modes
	if m.mode == ModeCreateSceneGroup {
		return m.updateCreateSceneGroupMode(msg)
	}
	if m.mode == ModeCreateSceneName {
		return m.updateCreateSceneNameMode(msg)
	}

	// Handle delete scene confirmation
	if m.mode == ModeDeleteSceneConfirm {
		return m.updateDeleteSceneConfirmMode(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, keys.TabNext):
			switch m.activeTab {
			case TabLights:
				m.activeTab = TabGroups
				return m, m.loadGroups
			case TabGroups:
				m.activeTab = TabScenes
				return m, m.loadScenes
			case TabScenes:
				m.activeTab = TabLights
				return m, m.loadLights
			}

		case key.Matches(msg, keys.TabPrev):
			switch m.activeTab {
			case TabLights:
				m.activeTab = TabScenes
				return m, m.loadScenes
			case TabGroups:
				m.activeTab = TabLights
				return m, m.loadLights
			case TabScenes:
				m.activeTab = TabGroups
				return m, m.loadGroups
			}

		case key.Matches(msg, keys.Up):
			switch m.activeTab {
			case TabLights:
				if m.lightCursor > 0 {
					m.lightCursor--
				}
			case TabGroups:
				if m.groupCursor > 0 {
					m.groupCursor--
				}
			case TabScenes:
				if m.sceneCursor > 0 {
					m.sceneCursor--
				}
			}

		case key.Matches(msg, keys.Down):
			switch m.activeTab {
			case TabLights:
				if m.lightCursor < len(m.lights)-1 {
					m.lightCursor++
				}
			case TabGroups:
				if m.groupCursor < len(m.groups)-1 {
					m.groupCursor++
				}
			case TabScenes:
				if m.sceneCursor < len(m.scenes)-1 {
					m.sceneCursor++
				}
			}

		case key.Matches(msg, keys.Toggle):
			switch m.activeTab {
			case TabLights:
				if len(m.lights) > 0 {
					light := m.lights[m.lightCursor]
					return m, m.toggleLight(light.ID, light.On)
				}
			case TabGroups:
				if len(m.groups) > 0 {
					group := m.groups[m.groupCursor]
					return m, m.toggleGroup(group.ID, group.AnyOn)
				}
			case TabScenes:
				if len(m.scenes) > 0 {
					scene := m.scenes[m.sceneCursor]
					return m, m.activateScene(scene.ID, scene.Name)
				}
			}

		case key.Matches(msg, keys.Confirm):
			// Enter also activates scenes
			if m.activeTab == TabScenes && len(m.scenes) > 0 {
				scene := m.scenes[m.sceneCursor]
				return m, m.activateScene(scene.ID, scene.Name)
			}

		case key.Matches(msg, keys.Rename):
			// Enter rename mode
			if m.activeTab == TabLights && len(m.lights) > 0 {
				light := m.lights[m.lightCursor]
				m.mode = ModeRename
				m.renameID = light.ID
				m.textInput.SetValue(light.Name)
				m.textInput.Focus()
				m.textInput.CursorEnd()
				return m, textinput.Blink
			} else if m.activeTab == TabGroups && len(m.groups) > 0 {
				group := m.groups[m.groupCursor]
				m.mode = ModeRename
				m.renameID = group.ID
				m.textInput.SetValue(group.Name)
				m.textInput.Focus()
				m.textInput.CursorEnd()
				return m, textinput.Blink
			}

		case key.Matches(msg, keys.Info):
			// Enter group info mode (only available on groups tab)
			if m.activeTab == TabGroups && len(m.groups) > 0 {
				group := m.groups[m.groupCursor]
				m.mode = ModeGroupInfo
				m.infoGroupID = group.ID
				return m, nil
			}

		case key.Matches(msg, keys.Delete):
			// Enter delete confirmation mode (groups or scenes tab)
			if m.activeTab == TabGroups && len(m.groups) > 0 {
				group := m.groups[m.groupCursor]
				m.mode = ModeDeleteConfirm
				m.deleteGroupID = group.ID
				m.deleteGroupName = group.Name
				return m, nil
			}
			if m.activeTab == TabScenes && len(m.scenes) > 0 {
				scene := m.scenes[m.sceneCursor]
				m.mode = ModeDeleteSceneConfirm
				m.deleteSceneID = scene.ID
				m.deleteSceneName = scene.Name
				return m, nil
			}

		case key.Matches(msg, keys.Add):
			// Enter create group mode (only available on groups tab)
			if m.activeTab == TabGroups {
				m.mode = ModeCreateGroupType
				m.createGroupType = ""
				m.createGroupName = ""
				m.createGroupLights = nil
				m.createLightCursor = 0
				m.createLightSelected = make(map[string]bool)
				return m, nil
			}
			// Enter create scene mode (only available on scenes tab)
			if m.activeTab == TabScenes {
				m.mode = ModeCreateSceneGroup
				m.createSceneGroupID = ""
				m.createSceneName = ""
				m.createGroupCursor = 0
				return m, nil
			}
		}

	case lightsLoadedMsg:
		m.lights = msg.lights
		m.lightsLoaded = true
		m.err = nil

	case groupsLoadedMsg:
		m.groups = msg.groups
		m.groupsLoaded = true
		m.err = nil

	case lightToggledMsg:
		for i := range m.lights {
			if m.lights[i].ID == msg.id {
				m.lights[i].On = msg.newOn
				break
			}
		}
		m.err = nil
		// Refresh groups to update AnyOn/AllOn state
		return m, m.loadGroups

	case groupToggledMsg:
		for i := range m.groups {
			if m.groups[i].ID == msg.id {
				m.groups[i].AllOn = msg.newOn
				m.groups[i].AnyOn = msg.newOn
				break
			}
		}
		m.err = nil
		// Refresh lights to update individual light states
		return m, m.loadLights

	case lightRenamedMsg:
		for i := range m.lights {
			if m.lights[i].ID == msg.id {
				m.lights[i].Name = msg.newName
				break
			}
		}
		m.err = nil

	case groupRenamedMsg:
		for i := range m.groups {
			if m.groups[i].ID == msg.id {
				m.groups[i].Name = msg.newName
				break
			}
		}
		m.err = nil

	case groupDeletedMsg:
		// Remove group from list
		for i := range m.groups {
			if m.groups[i].ID == msg.id {
				m.groups = append(m.groups[:i], m.groups[i+1:]...)
				break
			}
		}
		// Adjust cursor if needed
		if m.groupCursor >= len(m.groups) && m.groupCursor > 0 {
			m.groupCursor--
		}
		m.err = nil

	case groupCreatedMsg:
		// Add group to list and select it
		m.groups = append(m.groups, msg.group)
		m.groupCursor = len(m.groups) - 1
		m.err = nil
		// Refresh to get accurate state
		return m, m.loadGroups

	case scenesLoadedMsg:
		m.scenes = msg.scenes
		m.scenesLoaded = true
		m.err = nil

	case sceneActivatedMsg:
		m.err = nil
		// Refresh lights to show the new state
		return m, m.loadLights

	case sceneCreatedMsg:
		// Add scene to list and select it
		m.scenes = append(m.scenes, msg.scene)
		m.sceneCursor = len(m.scenes) - 1
		m.err = nil
		// Refresh to get accurate scene list
		return m, m.loadScenes

	case sceneDeletedMsg:
		// Remove scene from list
		for i := range m.scenes {
			if m.scenes[i].ID == msg.id {
				m.scenes = append(m.scenes[:i], m.scenes[i+1:]...)
				break
			}
		}
		// Adjust cursor if needed
		if m.sceneCursor >= len(m.scenes) && m.sceneCursor > 0 {
			m.sceneCursor--
		}
		m.err = nil

	case errMsg:
		m.err = msg.err
	}

	return m, nil
}

// updateGroupInfoMode handles input in group info mode.
func (m Model) updateGroupInfoMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Cancel), key.Matches(msg, keys.Quit):
			m.mode = ModeNormal
			m.infoGroupID = ""
			return m, nil
		}
	}
	return m, nil
}

// updateDeleteConfirmMode handles input in delete confirmation mode.
func (m Model) updateDeleteConfirmMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Yes), key.Matches(msg, keys.Confirm):
			// Confirm delete
			id := m.deleteGroupID
			m.mode = ModeNormal
			m.deleteGroupID = ""
			m.deleteGroupName = ""
			return m, m.deleteGroup(id)

		case key.Matches(msg, keys.No), key.Matches(msg, keys.Cancel):
			// Cancel delete
			m.mode = ModeNormal
			m.deleteGroupID = ""
			m.deleteGroupName = ""
			return m, nil
		}
	}
	return m, nil
}

// updateDeleteSceneConfirmMode handles input in delete scene confirmation mode.
func (m Model) updateDeleteSceneConfirmMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Yes), key.Matches(msg, keys.Confirm):
			// Confirm delete
			id := m.deleteSceneID
			m.mode = ModeNormal
			m.deleteSceneID = ""
			m.deleteSceneName = ""
			return m, m.deleteScene(id)

		case key.Matches(msg, keys.No), key.Matches(msg, keys.Cancel):
			// Cancel delete
			m.mode = ModeNormal
			m.deleteSceneID = ""
			m.deleteSceneName = ""
			return m, nil
		}
	}
	return m, nil
}

// updateCreateGroupTypeMode handles type selection (Room/Zone).
func (m Model) updateCreateGroupTypeMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Room):
			m.createGroupType = "Room"
			m.mode = ModeCreateGroupName
			m.textInput.SetValue("")
			m.textInput.Focus()
			return m, textinput.Blink

		case key.Matches(msg, keys.Zone):
			m.createGroupType = "Zone"
			m.mode = ModeCreateGroupName
			m.textInput.SetValue("")
			m.textInput.Focus()
			return m, textinput.Blink

		case key.Matches(msg, keys.Cancel):
			m.mode = ModeNormal
			return m, nil
		}
	}
	return m, nil
}

// updateCreateGroupNameMode handles name entry.
func (m Model) updateCreateGroupNameMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Confirm):
			name := m.textInput.Value()
			if name == "" {
				return m, nil // Don't proceed with empty name
			}
			m.createGroupName = name
			m.textInput.Blur()
			m.mode = ModeCreateGroupLights
			m.createLightCursor = 0
			m.createLightSelected = make(map[string]bool)
			return m, nil

		case key.Matches(msg, keys.Cancel):
			m.textInput.Blur()
			m.mode = ModeNormal
			return m, nil
		}
	}

	// Forward other messages to text input
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// updateCreateGroupLightsMode handles light selection.
func (m Model) updateCreateGroupLightsMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Up):
			if m.createLightCursor > 0 {
				m.createLightCursor--
			}

		case key.Matches(msg, keys.Down):
			if m.createLightCursor < len(m.lights)-1 {
				m.createLightCursor++
			}

		case key.Matches(msg, keys.Toggle):
			// Toggle light selection
			if len(m.lights) > 0 {
				light := m.lights[m.createLightCursor]
				m.createLightSelected[light.ID] = !m.createLightSelected[light.ID]
			}

		case key.Matches(msg, keys.Confirm):
			// Collect selected light IDs
			var lightIDs []string
			for _, light := range m.lights {
				if m.createLightSelected[light.ID] {
					lightIDs = append(lightIDs, light.ID)
				}
			}
			m.mode = ModeNormal
			return m, m.createGroup(m.createGroupName, m.createGroupType, lightIDs)

		case key.Matches(msg, keys.Cancel):
			m.mode = ModeNormal
			return m, nil
		}
	}
	return m, nil
}

// updateCreateSceneGroupMode handles group selection for scene creation.
func (m Model) updateCreateSceneGroupMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Up):
			if m.createGroupCursor > 0 {
				m.createGroupCursor--
			}

		case key.Matches(msg, keys.Down):
			if m.createGroupCursor < len(m.groups)-1 {
				m.createGroupCursor++
			}

		case key.Matches(msg, keys.Confirm):
			if len(m.groups) > 0 {
				group := m.groups[m.createGroupCursor]
				m.createSceneGroupID = group.ID
				m.mode = ModeCreateSceneName
				m.textInput.SetValue("")
				m.textInput.Focus()
				return m, textinput.Blink
			}

		case key.Matches(msg, keys.Cancel):
			m.mode = ModeNormal
			return m, nil
		}
	}
	return m, nil
}

// updateCreateSceneNameMode handles name entry for scene creation.
func (m Model) updateCreateSceneNameMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Confirm):
			name := m.textInput.Value()
			if name == "" {
				return m, nil // Don't proceed with empty name
			}
			m.createSceneName = name
			m.textInput.Blur()
			m.mode = ModeNormal
			return m, m.createScene(m.createSceneName, m.createSceneGroupID)

		case key.Matches(msg, keys.Cancel):
			m.textInput.Blur()
			m.mode = ModeNormal
			return m, nil
		}
	}

	// Forward other messages to text input
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// updateRenameMode handles input in rename mode.
func (m Model) updateRenameMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Confirm):
			// Submit rename
			newName := m.textInput.Value()
			m.mode = ModeNormal
			m.textInput.Blur()

			if m.activeTab == TabLights {
				return m, m.renameLight(m.renameID, newName)
			}
			return m, m.renameGroup(m.renameID, newName)

		case key.Matches(msg, keys.Cancel):
			// Cancel rename
			m.mode = ModeNormal
			m.textInput.Blur()
			return m, nil
		}
	}

	// Forward other messages to text input
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the UI.
// Run starts the TUI.
func Run(client *hue.Client) error {
	p := tea.NewProgram(New(client))
	_, err := p.Run()
	return err
}
