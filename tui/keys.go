package tui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Append key.Binding
	Quit   key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Append,
		k.Quit,
	}
}

var Keys = KeyMap{
	Append: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "append a new chat"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
