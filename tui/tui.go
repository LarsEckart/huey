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
)

// Mode represents the current interaction mode.
type Mode int

const (
	ModeNormal Mode = iota
	ModeRename
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
	TabNext key.Binding
	TabPrev key.Binding
	Quit    key.Binding
	Confirm key.Binding
	Cancel  key.Binding
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
		key.WithKeys("enter", " "),
		key.WithHelp("enter", "toggle"),
	),
	Rename: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rename"),
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
}

// Model is the Bubble Tea model for the TUI.
type Model struct {
	client      *hue.Client
	lights      []hue.Light
	groups      []hue.Group
	activeTab   Tab
	lightCursor int
	groupCursor int
	err         error
	quitting    bool

	// Rename mode
	mode      Mode
	textInput textinput.Model
	renameID  string // ID of item being renamed
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

// Init initializes the model and loads data.
func (m Model) Init() tea.Cmd {
	return tea.Batch(m.loadLights, m.loadGroups)
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

// Update handles messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle rename mode separately
	if m.mode == ModeRename {
		return m.updateRenameMode(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, keys.TabNext):
			if m.activeTab == TabLights {
				m.activeTab = TabGroups
			} else {
				m.activeTab = TabLights
			}

		case key.Matches(msg, keys.TabPrev):
			if m.activeTab == TabGroups {
				m.activeTab = TabLights
			} else {
				m.activeTab = TabGroups
			}

		case key.Matches(msg, keys.Up):
			if m.activeTab == TabLights && m.lightCursor > 0 {
				m.lightCursor--
			} else if m.activeTab == TabGroups && m.groupCursor > 0 {
				m.groupCursor--
			}

		case key.Matches(msg, keys.Down):
			if m.activeTab == TabLights && m.lightCursor < len(m.lights)-1 {
				m.lightCursor++
			} else if m.activeTab == TabGroups && m.groupCursor < len(m.groups)-1 {
				m.groupCursor++
			}

		case key.Matches(msg, keys.Toggle):
			if m.activeTab == TabLights && len(m.lights) > 0 {
				light := m.lights[m.lightCursor]
				return m, m.toggleLight(light.ID, light.On)
			} else if m.activeTab == TabGroups && len(m.groups) > 0 {
				group := m.groups[m.groupCursor]
				return m, m.toggleGroup(group.ID, group.AnyOn)
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
		}

	case lightsLoadedMsg:
		m.lights = msg.lights
		m.err = nil

	case groupsLoadedMsg:
		m.groups = msg.groups
		m.err = nil

	case lightToggledMsg:
		for i := range m.lights {
			if m.lights[i].ID == msg.id {
				m.lights[i].On = msg.newOn
				break
			}
		}
		m.err = nil

	case groupToggledMsg:
		for i := range m.groups {
			if m.groups[i].ID == msg.id {
				m.groups[i].AllOn = msg.newOn
				m.groups[i].AnyOn = msg.newOn
				break
			}
		}
		m.err = nil

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

	case errMsg:
		m.err = msg.err
	}

	return m, nil
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

	s := titleStyle.Render("huey - Hue Light Control") + "\n\n"

	// Render tabs
	s += m.renderTabs() + "\n\n"

	if m.err != nil {
		s += fmt.Sprintf("Error: %v\n\n", m.err)
	}

	// Render active tab content
	if m.activeTab == TabLights {
		s += m.renderLights()
	} else {
		s += m.renderGroups()
	}

	// Render help based on mode
	if m.mode == ModeRename {
		s += "\n" + helpStyle.Render("enter confirm • esc cancel")
	} else {
		s += "\n" + helpStyle.Render("↑/↓ navigate • enter toggle • r rename • tab switch • q quit")
	}

	return s
}

func (m Model) renderTabs() string {
	var lightsTab, groupsTab string

	if m.activeTab == TabLights {
		lightsTab = tabActiveStyle.Render("Lights")
		groupsTab = tabInactiveStyle.Render("Groups")
	} else {
		lightsTab = tabInactiveStyle.Render("Lights")
		groupsTab = tabActiveStyle.Render("Groups")
	}

	return lightsTab + "  " + groupsTab
}

func (m Model) renderLights() string {
	if len(m.lights) == 0 {
		return "Loading lights...\n"
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

		line := fmt.Sprintf("%s%-3s %s %s", cursor, light.ID, name, status)
		s += style.Render(line) + "\n"
	}

	return s
}

func (m Model) renderGroups() string {
	if len(m.groups) == 0 {
		return "Loading groups...\n"
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

		line := fmt.Sprintf("%s%-3s %s %-8s %s", cursor, group.ID, name, groupType, status)
		s += style.Render(line) + "\n"
	}

	return s
}

// Run starts the TUI.
func Run(client *hue.Client) error {
	p := tea.NewProgram(New(client))
	_, err := p.Run()
	return err
}
