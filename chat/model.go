package chat

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
	messages []string // TODO: Maybe limit the size of this lol
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

// TODO: Read up on this
// Still not 100% sure how this works
// Feels like bubbletea magic
func listenForActivity(sub chan twitch.PrivateMessage, client *twitch.Client) tea.Cmd {
	return func() tea.Msg {
		client.OnPrivateMessage(func(message twitch.PrivateMessage) {
			sub <- message
		})

		err := client.Connect()
		if err != nil {
			panic(err)
		}
		return nil
	}
}

func waitForActivity(sub chan twitch.PrivateMessage) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

func (m ChatModel) currentChannel(channel string) bool {
	return strings.EqualFold(m.channel, channel)
}

func (m ChatModel) Init() tea.Cmd {
	m.client.Join(m.channel) // Join the channel first
	return tea.Batch(
		m.spinner.Tick,                     // Start the spinner
		listenForActivity(m.sub, m.client), // Start accepting messages
		waitForActivity(m.sub),             // Wait to read the messages
	)
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case twitch.PrivateMessage:
		// FIXME: Hacky fix to only show messages from the current channel
		// Having race conditions here
		if !m.currentChannel(msg.Channel) {
			return m, waitForActivity(m.sub)
		}

		m.messages = append(m.messages, FormatMessage(msg))
		return m, waitForActivity(m.sub)
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m ChatModel) View() string {
	view := strings.Builder{}
	if len(m.messages) > 0 {
		// FIXME: This breaks on too many messages
		view.WriteString(strings.Join(m.messages, "\n"))
	} else {
		view.WriteString(fmt.Sprintf("%s No messages yet.", m.spinner.View()))
	}
	return view.String()
}
