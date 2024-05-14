package chat

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gempir/go-twitch-irc/v4"
)

type ChatModel struct {
	sub          chan twitch.PrivateMessage
	client       *twitch.Client
	Channel      string
	viewport     viewport.Model
	messages     []string
	spinner      spinner.Model
	MessageCount int
}

func NewChatModel(client *twitch.Client, spinner spinner.Model, channel string) ChatModel {
	sub := make(chan twitch.PrivateMessage)
	model := ChatModel{
		sub:      sub,
		client:   client,
		spinner:  spinner,
		Channel:  channel,
		viewport: viewport.New(0, 0), // Init viewport to (0,0), see update
	}

	model.viewport.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("62"))

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

// Force re-render message, makes viewport size good
func (m ChatModel) SetViewportSize(width int, height int) tea.Cmd {
	return func() tea.Msg {
		return tea.WindowSizeMsg{Width: width, Height: height}
	}
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

// TODO: Update this
func (m ChatModel) renderViewportInfo() string {
	top := "false"
	bottom := "false"

	if m.viewport.AtTop() {
		top = "true"
	}

	if m.viewport.AtBottom() {
		bottom = "true"
	}

	return fmt.Sprintf("Total %d, Visible %d, At top %s, At bottom %s", m.viewport.TotalLineCount(), m.viewport.VisibleLineCount(), top, bottom)
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// TODO: Dynamic?
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 12

		// Re-render the viewport
		cmds = append(cmds, viewport.Sync(m.viewport))

	case tea.KeyMsg:
		// FIXME: WHY IS THIS NOT WORKING?
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)

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
	if len(m.messages) > 0 {
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
	} else {
		m.viewport.SetContent(fmt.Sprintf("%s No messages yet.", m.spinner.View()))
	}

	return (m.viewport.View() + "\n" + m.renderViewportInfo())
}
