package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/gempir/go-twitch-irc/v4"
)

func NewStyledSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return s
}

func NewTextInput() textinput.Model {
	ti := textinput.New()     // Create new text input model
	ti.Placeholder = "forsen" // Default placeholder
	ti.CharLimit = 25         // Twitch username limit
	ti.Width = 20             // Width of the text input
	return ti
}

// Yoink from https://github.com/charmbracelet/bubbletea/blob/master/examples/tabs/main.go
func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.NormalBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	tabBorder        = lipgloss.NormalBorder()
	docStyle         = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor   = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle = lipgloss.NewStyle().Border(tabBorder, true)
	activeTabStyle   = inactiveTabStyle.Copy().BorderForeground(highlightColor)
	windowStyle      = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(0, 0).Align(lipgloss.Left).Border(lipgloss.NormalBorder())
)

func renderTabString(channels []string, active string) string {
	var renderedTabs []string

	for _, channel := range channels {
		var style lipgloss.Style
		isActive := active == channel
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}
		renderedTabs = append(renderedTabs, style.Render(channel))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	return row
}

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

	return fullStyle.Render("@" + message.Channel + "> " + userStyle.Render(RoleString(message.User)+message.User.Name) + ": " + contentStyle.Render(message.Message))
}
