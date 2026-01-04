package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/LarsEckart/huey/hue"
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

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Padding(0, 1)

	tabActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Padding(0, 1)

	tabInactiveStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252")).
				Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	onStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("220"))

	offStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	partialStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("208"))

	typeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229"))
)

// keyMap defines key bindings.
type keyMap struct {
	Up      key.Binding
	Down    key.Binding
	Toggle  key.Binding
	Rename  key.Binding
	Delete  key.Binding
	Add     key.Binding
	Info    key.Binding
	TabNext key.Binding
	TabPrev key.Binding
	Quit    key.Binding
	Confirm key.Binding
	Cancel  key.Binding
	Yes     key.Binding
	No      key.Binding
	Room    key.Binding
	Zone    key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓", "down"),
	),
	Toggle: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "toggle"),
	),
	Rename: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rename"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add"),
	),
	Info: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "info"),
	),
	TabNext: key.NewBinding(
		key.WithKeys("tab", "l"),
		key.WithHelp("tab", "next tab"),
	),
	TabPrev: key.NewBinding(
		key.WithKeys("shift+tab", "h"),
		key.WithHelp("shift+tab", "prev tab"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	Yes: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "yes"),
	),
	No: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "no"),
	),
	Room: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "room"),
	),
	Zone: key.NewBinding(
		key.WithKeys("z"),
		key.WithHelp("z", "zone"),
	),
}

// Model is the Bubble Tea model for the TUI.
type Model struct {
	client        *hue.Client
	lights        []hue.Light
	groups        []hue.Group
	scenes        []hue.Scene
	lightsLoaded  bool
	groupsLoaded  bool
	scenesLoaded  bool
	activeTab     Tab
	lightCursor   int
	groupCursor   int
	sceneCursor   int
	err           error
	quitting      bool

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
	createGroupType     string   // "Room" or "Zone"
	createGroupName     string   // Name entered by user
	createGroupLights   []string // Selected light IDs
	createLightCursor   int      // Cursor for light picker
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

// Messages
type lightsLoadedMsg struct {
	lights []hue.Light
}

type groupsLoadedMsg struct {
	groups []hue.Group
}

type errMsg struct {
	err error
}

type lightToggledMsg struct {
	id    string
	newOn bool
}

type groupToggledMsg struct {
	id    string
	newOn bool
}

type lightRenamedMsg struct {
	id      string
	newName string
}

type groupRenamedMsg struct {
	id      string
	newName string
}

type groupDeletedMsg struct {
	id string
}

type groupCreatedMsg struct {
	group hue.Group
}

type scenesLoadedMsg struct {
	scenes []hue.Scene
}

type sceneActivatedMsg struct {
	id   string
	name string
}

type sceneCreatedMsg struct {
	scene hue.Scene
}

type sceneDeletedMsg struct {
	id string
}

// Init initializes the model and loads data.
func (m Model) Init() tea.Cmd {
	return tea.Batch(m.loadLights, m.loadGroups, m.loadScenes)
}

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
		s += fmt.Sprintf("Error: %v\n\n", m.err)
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
	if m.mode == ModeRename {
		s += "\n" + helpStyle.Render("enter confirm • esc cancel")
	} else if m.mode == ModeDeleteConfirm || m.mode == ModeDeleteSceneConfirm {
		s += "\n" + helpStyle.Render("y/enter delete • n/esc cancel")
	} else {
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
	if !m.lightsLoaded {
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
	if !m.groupsLoaded {
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
	if !m.scenesLoaded {
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
			s += fmt.Sprintf("  (unknown light)\n")
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

// Run starts the TUI.
func Run(client *hue.Client) error {
	p := tea.NewProgram(New(client))
	_, err := p.Run()
	return err
}
