package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gempir/go-twitch-irc/v4"
)

// This is the example I based this on
// https://github.com/charmbracelet/bubbletea/tree/master/examples/realtime

// Start async to listen for new messages, async because of tea.Cmd.
// Not fully 100% on how this works yet... hehe
func listenForActivity(sub chan twitch.PrivateMessage, client *twitch.Client) tea.Cmd {
	return func() tea.Msg {
		// TODO: This could probably be defined elsewhere /shrug
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

type model struct {
	sub      chan twitch.PrivateMessage
	client   *twitch.Client
	channel  string
	messages []string
	spinner  spinner.Model
}

func (m model) Init() tea.Cmd {
	m.client.Join(m.channel) // Join the channel first

	return tea.Batch(
		m.spinner.Tick,                     // Start ticking the spinner
		listenForActivity(m.sub, m.client), // Start listening to Twitch chat
		waitForActivity(m.sub),             // Start waiting for the first message
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case twitch.PrivateMessage:
		m.messages = append(m.messages, msg.User.Name+": "+msg.Message) // Append the message
		if len(m.messages) > 20 {                                       // Limit how many messages are displayed
			m.messages = m.messages[1:]
		}
		return m, waitForActivity(m.sub) // Continue waiting for the next message
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	s := "\nIn chat " + m.channel + ":\n\n"
	if len(m.messages) > 0 {
		s += strings.Join(m.messages, "\n")
	} else {
		s += fmt.Sprintf("%s No messages yet.", m.spinner.View())
	}
	s += "\n\nPress any key to exit\n"
	return s
}

func main() {
	client := twitch.NewAnonymousClient()

	p := tea.NewProgram(model{
		sub:     make(chan twitch.PrivateMessage),
		client:  client,
		spinner: spinner.New(),
		channel: "tarik", // Placeholder
	})

	// Run the UI
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
