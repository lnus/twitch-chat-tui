// TODO:
// Might rework this file
package chat

import (
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

	return fullStyle.Render("@" + message.Channel + "> " + userStyle.Render(message.User.Name) + ": " + contentStyle.Render(message.Message))
}
