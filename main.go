package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gempir/go-twitch-irc/v4"
)

type model struct {
	messageChan chan string // Channel for incoming messages
	messages    []string    // List of messages to display
}

func initialModel() model {
	return model{
		messages:    []string{},
		messageChan: make(chan string, 100), // Buffer 100 messages
	}
}

func listenTwitchMessages(messageChan chan string) tea.Cmd {
	return func() tea.Msg {
		client := twitch.NewAnonymousClient()

		client.OnPrivateMessage(func(message twitch.PrivateMessage) {
			messageChan <- message.Message
		})

		client.Join("forsen") // Join placeholder channel

		err := client.Connect()
		if err != nil {
			panic(err)
		}
		return nil
	}
}

func (m model) Init() tea.Cmd {
	// Start listening for messages
	return listenTwitchMessages(m.messageChan)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg: // Keypress event
		if msg.String() == "q" {
			return m, tea.Quit
		}
	case string: // Message from Twitch
		m.messages = append(m.messages, msg)
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	view := "Twitch chat:\n"
	for _, message := range m.messages {
		view += message + "\n"
	}
	view += "Press 'q' to quit"
	return view
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting program:", err)
	}
}
