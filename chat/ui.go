// TODO:
// Might rework this file
package chat

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/gempir/go-twitch-irc/v4"
)

// TODO:
// This is very much a WIP
// But /shrug
func RoleString(user twitch.User) string {
	rolemap := map[string]string{
		"broadcaster": "◉",
		"moderator":   "⛨",
		"subscriber":  "✪",
		"partner":     "✓",
	}

	s := ""
	for k := range user.Badges {
		if rolemap[k] != "" {
			s += rolemap[k] + " "
		} else {
			// TODO: This is debug
			// s += k + " "
			continue
		}
	}

	return s
}

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

	return fullStyle.Render(userStyle.Render(RoleString(message.User)+message.User.Name) + ": " + contentStyle.Render(message.Message))
}
