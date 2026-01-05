package tui

import (
	"github.com/LarsEckart/huey/hue"
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

// Run starts the TUI.
func Run(client *hue.Client) error {
	p := tea.NewProgram(New(client))
	_, err := p.Run()
	return err
}
