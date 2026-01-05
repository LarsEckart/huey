package tui

import "github.com/charmbracelet/bubbles/key"

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
