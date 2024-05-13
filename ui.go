package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/gempir/go-twitch-irc/v4"
)

func FormatMessage(message twitch.PrivateMessage) string {
	// Add lipgloss styling to make the username the users color
	userStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(message.User.Color))

	contentStyle := lipgloss.NewStyle().
		Bold(false)

	fullStyle := lipgloss.NewStyle()
	if message.FirstMessage {
		fullStyle.Background(lipgloss.Color("201"))
	}

	return fullStyle.Render(userStyle.Render(message.User.Name) + ": " + contentStyle.Render(message.Message))
}

func NewStyledSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return s
}

func RenderView(m model) string {
	s := "\nIn chat " + m.channel + ":\n\n"
	if len(m.messages) > 0 {
		s += strings.Join(m.messages, "\n")
	} else {
		s += fmt.Sprintf("%s No messages yet.", m.spinner.View())
	}
	s += "\n\nPress any key to exit\n"
	return s
}
