package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gempir/go-twitch-irc/v4"
)

func main() {
	client := twitch.NewAnonymousClient()

	p := tea.NewProgram(initializeModel(client, NewStyledSpinner(), "forsen"), tea.WithAltScreen())

	// Run the UI
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
