package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gempir/go-twitch-irc/v4"
)

func main() {
	client := twitch.NewAnonymousClient()
	channel := "forsen" // TODO: Make dynamic per tab
	chatModelOnly := false

	var p *tea.Program

	if chatModelOnly {
		p = tea.NewProgram(NewChatModel(client, NewStyledSpinner(), channel), tea.WithAltScreen())
	} else {
		p = tea.NewProgram(NewMainModel(), tea.WithAltScreen())
	}

	// Run the UI
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
