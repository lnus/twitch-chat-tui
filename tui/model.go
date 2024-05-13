package tui

import (
	"strings"
	"ttui/chat"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gempir/go-twitch-irc/v4"
)

type MainModel struct {
	chatModels map[string]chat.ChatModel
	help       help.Model
	activeChat string
	keys       KeyMap
	tabs       []string
	textInput  textinput.Model
	width      int
	height     int
	isTyping   bool
}

func NewMainModel() MainModel {
	return MainModel{
		keys:       Keys,                            // Keybindings
		help:       help.New(),                      // Help menu model
		chatModels: make(map[string]chat.ChatModel), // Chat models
		tabs:       make([]string, 0),               // Tab representation of chat models
		activeChat: "",                              // Active chat, default to none
		textInput:  NewTextInput(),                  // Text input model
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

func (m MainModel) channelIndex() int {
	for i, tab := range m.tabs {
		if tab == m.activeChat {
			return i
		}
	}
	return -1
}

func (m MainModel) passUpdates(msg tea.Msg) []tea.Cmd {
	var cmds []tea.Cmd
	for _, chatModel := range m.chatModels {
		updatedModel, cmd := chatModel.Update(msg)
		m.chatModels[chatModel.Channel] = updatedModel.(chat.ChatModel)
		cmds = append(cmds, cmd)
	}
	return cmds
}

func (m *MainModel) addChannel(channel string) tea.Cmd {
	// Sanity check, just set active
	if m.channelExists(channel) {
		m.activeChat = channel
		return nil
	}

	// Create a new ChatModel with the username
	chatModel := chat.NewChatModel(twitch.NewAnonymousClient(), NewStyledSpinner(), channel)
	m.chatModels[channel] = chatModel
	m.activeChat = channel

	// Add this channel name to last of tab-list
	m.tabs = append(m.tabs, channel)

	return chatModel.Init()
}

// TODO: Abstract some of this logic
// Probably into receiver methods for MainModel
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	index := m.channelIndex()

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update all chat models
		cmds = append(cmds, m.passUpdates(msg)...)

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

				cmd = m.addChannel(username)
				cmds = append(cmds, cmd)
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

				// Tab switching
			case "left", "h": // Left arrow or h
				if index != -1 {
					// Go left and wrap around
					if index == 0 {
						m.activeChat = m.tabs[len(m.tabs)-1]
					} else {
						m.activeChat = m.tabs[index-1]
					}
				}
			case "right", "l": // Right arrow or l
				if index != -1 {
					// Go right and wrap around
					if index == len(m.tabs)-1 {
						m.activeChat = m.tabs[0]
					} else {
						m.activeChat = m.tabs[index+1]
					}
				}
			}
		}
	default:
		// FIXME: Update all chats, not only active
		// if m.channelExists(m.activeChat) {
		// 	updatedModel, cmd := m.chatModels[m.activeChat].Update(msg)
		// 	m.chatModels[m.activeChat] = updatedModel.(chat.ChatModel)
		// 	cmds = append(cmds, cmd)
		// }
		// pass updates to all chat models
		cmds = append(cmds, m.passUpdates(msg)...)

		// Iterate over all chat models and update them
		// For some reason this consolidates the messages
		// into one chat model? Strange.
		// for channel, chatModel := range m.chatModels {
		// 	updatedModel, cmd := chatModel.Update(msg)
		// 	m.chatModels[channel] = updatedModel.(chat.ChatModel)
		// 	cmds = append(cmds, cmd)
		// }
	}

	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	view := strings.Builder{}
	helpView := m.help.ShortHelpView(m.keys.ShortHelp())

	if len(m.chatModels) > 0 {
		// Iterate over chat models to make tabs
		row := renderTabString(m.tabs, m.activeChat)
		view.WriteString(row)
		view.WriteString("\n")

		// And render content within tab
		// TODO: Width & height should be more dynamic
		view.WriteString(windowStyle.
			Width(m.width - 4).
			Height(m.height - 12).
			Render(m.chatModels[m.activeChat].View()))
	}

	if m.isTyping {
		if len(m.chatModels) > 0 {
			view.WriteString("\n\n")
		}
		view.WriteString(m.textInput.View())
	}

	// Mini help display
	view.WriteString("\n\n" + helpView)

	// TODO: Rename docStyle
	return docStyle.Render(view.String())
}
