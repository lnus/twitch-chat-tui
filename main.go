package main

import (
	"fmt"
	"os"
	"ttui/chat"
	"ttui/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gempir/go-twitch-irc/v4"
)

func main() {
	client := twitch.NewAnonymousClient()
	channel := "forsen" // TODO: Make dynamic per tab
	chatModelOnly := false

	var p *tea.Program

	if chatModelOnly {
		p = tea.NewProgram(chat.NewChatModel(client, tui.NewStyledSpinner(), channel), tea.WithAltScreen())
	} else {
		p = tea.NewProgram(tui.NewMainModel(), tea.WithAltScreen())
	}

	// Run the UI
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
