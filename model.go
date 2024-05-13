package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gempir/go-twitch-irc/v4"
)

type MainModel struct {
	chatModels []ChatModel
	textInput  textinput.Model
	activeChat int
	isTyping   bool
}

func NewMainModel() MainModel {
	ti := textinput.New()
	ti.Placeholder = "forsen"
	ti.Focus()
	ti.CharLimit = 25 // Twitch username limit
	ti.Width = 20

	return MainModel{
		textInput: ti,
		isTyping:  false,
	}
}

func (m MainModel) Init() tea.Cmd {
	return textinput.Blink // Cursor blink start
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "a":
			// TextInput pop up
			m.isTyping = true
			return m, nil
		case "enter":
			if m.isTyping {
				m.isTyping = false

				// Create a new ChatModel with the username
				chatModel := NewChatModel(twitch.NewAnonymousClient(), NewStyledSpinner(), m.textInput.Value())

				m.chatModels = append(m.chatModels, chatModel)

				// Set that chat as the active chat
				// This is a bit of a hack, but it works for now
				m.activeChat = len(m.chatModels) - 1

				return m, m.chatModels[m.activeChat].Init()
			}
		}
	}

	// Propogate update to activeChat
	if len(m.chatModels) > 0 {
		updatedModel, cmd := m.chatModels[m.activeChat].Update(msg)
		m.chatModels[m.activeChat] = updatedModel.(ChatModel)
		cmds = append(cmds, cmd)
	}

	// Might be a silly way to do this, but we're doing our best
	if m.isTyping {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	var s string

	if len(m.chatModels) == 0 && !m.isTyping {
		s += "No chats yet. Press 'a' to start typing."
	}

	if m.isTyping {
		s += m.textInput.View()
	}

	if len(m.chatModels) > 0 {
		s += m.chatModels[m.activeChat].View()
	}

	return s
}
