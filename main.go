package main

import (
	"fmt"
	"os"
	"ttui/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(tui.NewMainModel(), tea.WithAltScreen())

	// Run the UI
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
