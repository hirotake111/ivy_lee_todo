package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	mainStyle           = lipgloss.NewStyle().MarginLeft(2)
	errorStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	docStyle            = lipgloss.NewStyle().Margin(1, 2)
	titleStyle          = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFDF5")).
				Background(lipgloss.Color("#25A065")).
				Padding(0, 1).Render
	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"})
)
