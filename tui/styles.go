package tui

import "github.com/charmbracelet/lipgloss"

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
