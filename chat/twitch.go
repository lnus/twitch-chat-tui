package chat

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gempir/go-twitch-irc/v4"
)

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
