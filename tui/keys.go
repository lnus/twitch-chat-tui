package tui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Append key.Binding
	Quit   key.Binding
	Left   key.Binding
	Right  key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Append,
		k.Quit,
		k.Left,
		k.Right,
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
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "go left in tabs"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("l/→", "go right in tabs"),
	),
}
