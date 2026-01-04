package tui

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/LarsEckart/huey/hue"
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
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

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

// keyMap defines key bindings.
type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Toggle key.Binding
	Quit   key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Toggle: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter", "toggle"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// Model is the Bubble Tea model for the TUI.
type Model struct {
	client   *hue.Client
	lights   []hue.Light
	cursor   int
	err      error
	quitting bool
}

// New creates a new TUI model.
func New(client *hue.Client) Model {
	return Model{
		client: client,
	}
}

// lightsLoadedMsg is sent when lights are loaded.
type lightsLoadedMsg struct {
	lights []hue.Light
}

// errMsg is sent when an error occurs.
type errMsg struct {
	err error
}

// lightToggledMsg is sent when a light is toggled.
type lightToggledMsg struct {
	id    string
	newOn bool
}

// Init initializes the model and loads lights.
func (m Model) Init() tea.Cmd {
	return m.loadLights
}

func (m Model) loadLights() tea.Msg {
	lights, err := m.client.GetLights()
	if err != nil {
		return errMsg{err: err}
	}

	// Sort by ID numerically for natural order
	sort.Slice(lights, func(i, j int) bool {
		iID, _ := strconv.Atoi(lights[i].ID)
		jID, _ := strconv.Atoi(lights[j].ID)
		return iID < jID
	})

	return lightsLoadedMsg{lights: lights}
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

// Update handles messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}

		case key.Matches(msg, keys.Down):
			if m.cursor < len(m.lights)-1 {
				m.cursor++
			}

		case key.Matches(msg, keys.Toggle):
			if len(m.lights) > 0 {
				light := m.lights[m.cursor]
				return m, m.toggleLight(light.ID, light.On)
			}
		}

	case lightsLoadedMsg:
		m.lights = msg.lights
		m.err = nil

	case lightToggledMsg:
		// Update local state
		for i := range m.lights {
			if m.lights[i].ID == msg.id {
				m.lights[i].On = msg.newOn
				break
			}
		}
		m.err = nil

	case errMsg:
		m.err = msg.err
	}

	return m, nil
}

// View renders the UI.
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	s := titleStyle.Render("huey - Hue Light Control") + "\n\n"

	if m.err != nil {
		s += fmt.Sprintf("Error: %v\n\n", m.err)
	}

	if len(m.lights) == 0 {
		s += "Loading lights...\n"
	} else {
		for i, light := range m.lights {
			// Cursor indicator
			cursor := "  "
			style := normalStyle
			if i == m.cursor {
				cursor = "> "
				style = selectedStyle
			}

			// Status indicator
			var status string
			if light.On {
				status = onStyle.Render("● on")
			} else {
				status = offStyle.Render("○ off")
			}

			line := fmt.Sprintf("%s%-3s %-24s %s", cursor, light.ID, light.Name, status)
			s += style.Render(line) + "\n"
		}
	}

	s += "\n" + helpStyle.Render("↑/↓ navigate • enter toggle • q quit")

	return s
}

// Run starts the TUI.
func Run(client *hue.Client) error {
	p := tea.NewProgram(New(client))
	_, err := p.Run()
	return err
}
