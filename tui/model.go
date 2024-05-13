package tui

import (
	// import from local package chat/model
	"ttui/chat"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gempir/go-twitch-irc/v4"
)

// Keybinding stuff
// TODO: Move to a separate file
type keyMap struct {
	Append key.Binding
	Quit   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Append,
		k.Quit,
	}
}

var keys = keyMap{
	Append: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "append a new chat"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type MainModel struct {
	chatModels map[string]chat.ChatModel
	help       help.Model
	keys       keyMap
	activeChat string
	textInput  textinput.Model
	isTyping   bool
}

func NewMainModel() MainModel {
	ti := textinput.New()
	ti.Placeholder = "forsen"
	ti.CharLimit = 25 // Twitch username limit
	ti.Width = 20

	return MainModel{
		keys:       keys,                            // Keybindings
		help:       help.New(),                      // Help menu model
		chatModels: make(map[string]chat.ChatModel), // Chat models
		activeChat: "",                              // Active chat, default to none
		textInput:  ti,                              // Text input model
		isTyping:   false,                           // Typing mode
	}
}

func (m MainModel) Init() tea.Cmd {
	return textinput.Blink // Cursor blink start
}

func (m MainModel) channelExists(channel string) bool {
	_, exists := m.chatModels[channel]
	return exists
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.isTyping {
			// Handle text input updates
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			cmds = append(cmds, cmd)
			switch msg.String() {
			case "esc":
				m.isTyping = false
			case "enter":
				m.isTyping = false
				username := m.textInput.Value()

				// If channel already exists, just swap to that model
				if m.channelExists(username) {
					m.activeChat = username
				} else {
					// Create a new ChatModel with the username
					chatModel := chat.NewChatModel(twitch.NewAnonymousClient(), NewStyledSpinner(), username)
					m.chatModels[username] = chatModel
					m.activeChat = username
					cmds = append(cmds, chatModel.Init())
				}
			}
		} else {
			// Handle non-input key presses
			switch msg.String() {
			case "q", "ctrl+c":
				cmds = append(cmds, tea.Quit)
			case "a":
				// TextInput pop up
				m.isTyping = true
				cmds = append(cmds, m.textInput.Focus())
			}
		}
	default:
		if m.channelExists(m.activeChat) {
			updatedModel, cmd := m.chatModels[m.activeChat].Update(msg)
			m.chatModels[m.activeChat] = updatedModel.(chat.ChatModel)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	var s string
	helpView := m.help.ShortHelpView(m.keys.ShortHelp())

	if len(m.chatModels) == 0 && !m.isTyping {
		s += "No chats yet. Press 'a' to start typing."
	}

	if m.isTyping {
		s += m.textInput.View()
	}

	if len(m.chatModels) > 0 {
		s += m.chatModels[m.activeChat].View()
	}

	// Mini help display
	s += "\n\n" + helpView

	return s
}