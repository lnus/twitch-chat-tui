package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gempir/go-twitch-irc/v4"
)

type ChatModel struct {
	sub      chan twitch.PrivateMessage
	client   *twitch.Client
	channel  string
	messages []string
	spinner  spinner.Model
}

func NewChatModel(client *twitch.Client, spinner spinner.Model, channel string) ChatModel {
	return ChatModel{
		sub:     make(chan twitch.PrivateMessage),
		client:  client,
		spinner: spinner,
		channel: channel,
	}
}

func (m ChatModel) Init() tea.Cmd {
	m.client.Join(m.channel) // Join the channel first
	return tea.Batch(
		m.spinner.Tick,
		listenForActivity(m.sub, m.client),
		waitForActivity(m.sub),
	)
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case twitch.PrivateMessage:
		m.messages = append(m.messages, FormatMessage(msg))
		if len(m.messages) > 20 {
			m.messages = m.messages[1:]
		}
		return m, waitForActivity(m.sub)
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m ChatModel) View() string {
	s := "\nIn chat " + m.channel + ":\n\n"
	if len(m.messages) > 0 {
		s += strings.Join(m.messages, "\n")
	} else {
		s += fmt.Sprintf("%s No messages yet.", m.spinner.View())
	}
	s += "\n\nPress any key to exit\n"
	return s
}
