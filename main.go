package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gempir/go-twitch-irc/v4"
)

// twitchMsg is a message type to signal that a new message has been received from Twitch.
type twitchMsg struct {
	message twitch.PrivateMessage
}

// listenForActivity starts a goroutine to connect to Twitch and listen for messages.
func listenForActivity(sub chan twitchMsg, client *twitch.Client) tea.Cmd {
	return func() tea.Msg {
		go func() {
			client.OnPrivateMessage(func(message twitch.PrivateMessage) {
				sub <- twitchMsg{message: message}
			})
			err := client.Connect()
			if err != nil {
				fmt.Println("Failed to connect to Twitch:", err)
				os.Exit(1)
			}
		}()
		return nil
	}
}

func waitForActivity(sub chan twitchMsg) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

type model struct {
	sub      chan twitchMsg
	client   *twitch.Client
	messages []string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		listenForActivity(m.sub, m.client), // Start listening to Twitch chat
		waitForActivity(m.sub),             // Start waiting for the first message
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.quitting = true
		// TODO: This throws an error atm, not sexy
		// m.client.Disconnect()
		return m, tea.Quit
	case twitchMsg:
		m.messages = append(m.messages, msg.message.User.Name+": "+msg.message.Message) // Append new message
		if len(m.messages) > 10 {                                                       // Keep only the last 10 messages for display
			m.messages = m.messages[1:]
		}
		return m, waitForActivity(m.sub) // Continue waiting for the next message
	default:
		return m, nil
	}
}

func (m model) View() string {
	s := "\nRecent Twitch chat messages:\n\n"
	if len(m.messages) > 0 {
		s += strings.Join(m.messages, "\n")
	} else {
		s += "No messages yet."
	}
	s += "\n\nPress any key to exit\n"
	if m.quitting {
		s += "\n"
	}
	return s
}

func main() {
	client := twitch.NewAnonymousClient()

	p := tea.NewProgram(model{
		sub:    make(chan twitchMsg),
		client: client,
	})

	client.Join("tarik")

	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
