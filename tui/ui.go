package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func NewStyledSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return s
}

func NewTextInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "forsen"
	ti.CharLimit = 25 // Twitch username limit
	ti.Width = 20
	return ti
}
