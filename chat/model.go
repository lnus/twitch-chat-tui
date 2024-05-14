package chat

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gempir/go-twitch-irc/v4"
)

type ChatModel struct {
	sub          chan twitch.PrivateMessage
	client       *twitch.Client
	Channel      string
	messages     []string
	spinner      spinner.Model
	MessageCount int
}

func NewChatModel(client *twitch.Client, spinner spinner.Model, channel string) ChatModel {
	sub := make(chan twitch.PrivateMessage)
	model := ChatModel{
		sub:     sub,
		client:  client,
		spinner: spinner,
		Channel: channel,
	}

	// Start listening for messages
	go listenForMessages(sub, client)

	return model
}

func listenForMessages(sub chan twitch.PrivateMessage, client *twitch.Client) {
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		sub <- message
	})

	if err := client.Connect(); err != nil {
		if err != twitch.ErrClientDisconnected {
			panic(err)
		}
	}
}

func waitForActivity(sub chan twitch.PrivateMessage) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

func (m ChatModel) currentChannel(channel string) bool {
	return strings.EqualFold(m.Channel, channel)
}

func (m ChatModel) Destroy() {
	m.client.Disconnect()
}

func (m ChatModel) Init() tea.Cmd {
	m.client.Join(m.Channel) // Join the channel first
	return tea.Batch(
		m.spinner.Tick,         // Start the spinner
		waitForActivity(m.sub), // Wait to read the messages
	)
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case twitch.PrivateMessage:
		if m.currentChannel(msg.Channel) {
			m.messages = append(m.messages, FormatMessage(msg))
			m.MessageCount++
		}
		cmds = append(cmds, waitForActivity(m.sub))
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m ChatModel) View() string {
	view := strings.Builder{}
	if len(m.messages) > 0 {
		view.WriteString(strings.Join(m.messages, "\n"))
	} else {
		view.WriteString(fmt.Sprintf("%s No messages yet.", m.spinner.View()))
	}

	return view.String()
}
