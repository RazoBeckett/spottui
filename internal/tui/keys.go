package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines all keybindings for the application
type KeyMap struct {
	// Navigation
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Enter key.Binding
	Back  key.Binding

	// Playback
	PlayPause  key.Binding
	Next       key.Binding
	Prev       key.Binding
	VolumeUp   key.Binding
	VolumeDown key.Binding
	Shuffle    key.Binding
	Repeat     key.Binding

	// App
	Help         key.Binding
	Quit         key.Binding
	Search       key.Binding
	GlobalSearch key.Binding
	Refresh      key.Binding
	Devices      key.Binding
	History      key.Binding
	Artist       key.Binding
	Like         key.Binding
	Lyrics       key.Binding
}

// DefaultKeyMap returns the default keybindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc", "backspace"),
			key.WithHelp("esc", "back"),
		),
		PlayPause: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "play/pause"),
		),
		Next: key.NewBinding(
			key.WithKeys("n", ">"),
			key.WithHelp("n/>", "next track"),
		),
		Prev: key.NewBinding(
			key.WithKeys("p", "<"),
			key.WithHelp("p/<", "prev track"),
		),
		VolumeUp: key.NewBinding(
			key.WithKeys("+", "="),
			key.WithHelp("+", "volume up"),
		),
		VolumeDown: key.NewBinding(
			key.WithKeys("-", "_"),
			key.WithHelp("-", "volume down"),
		),
		Shuffle: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "shuffle"),
		),
		Repeat: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "repeat"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
		GlobalSearch: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "global search"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("ctrl+r", "refresh"),
		),
		Devices: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "devices"),
		),
		History: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "recently played"),
		),
		Artist: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "view artist"),
		),
		Like: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "like"),
		),
		Lyrics: key.NewBinding(
			key.WithKeys("L"),
			key.WithHelp("L", "lyrics"),
		),
	}
}

// ShortHelp returns keybindings to show in the mini help view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Back, k.PlayPause, k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Enter, k.Back, k.Search, k.GlobalSearch},
		{k.PlayPause, k.Next, k.Prev},
		{k.VolumeUp, k.VolumeDown, k.Shuffle, k.Repeat},
		{k.Help, k.Quit, k.Refresh, k.Devices, k.History, k.Artist, k.Like, k.Lyrics},
	}
}
