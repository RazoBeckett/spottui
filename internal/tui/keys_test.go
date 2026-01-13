package tui

import (
	"testing"

	"github.com/charmbracelet/bubbles/key"
)

func TestDefaultKeyMap(t *testing.T) {
	km := DefaultKeyMap()

	tests := []struct {
		name    string
		binding key.Binding
		keys    []string
	}{
		{"Up", km.Up, []string{"up", "k"}},
		{"Down", km.Down, []string{"down", "j"}},
		{"Left", km.Left, []string{"left", "h"}},
		{"Right", km.Right, []string{"right", "l"}},
		{"Enter", km.Enter, []string{"enter"}},
		{"Back", km.Back, []string{"esc", "backspace"}},
		{"PlayPause", km.PlayPause, []string{" "}},
		{"Next", km.Next, []string{"n", ">"}},
		{"Prev", km.Prev, []string{"p", "<"}},
		{"VolumeUp", km.VolumeUp, []string{"+", "="}},
		{"VolumeDown", km.VolumeDown, []string{"-", "_"}},
		{"Shuffle", km.Shuffle, []string{"s"}},
		{"Repeat", km.Repeat, []string{"r"}},
		{"SeekBackward", km.SeekBackward, []string{"["}},
		{"SeekForward", km.SeekForward, []string{"]"}},
		{"Help", km.Help, []string{"?"}},
		{"Quit", km.Quit, []string{"q", "ctrl+c"}},
		{"Search", km.Search, []string{"/"}},
		{"GlobalSearch", km.GlobalSearch, []string{"S"}},
		{"Refresh", km.Refresh, []string{"ctrl+r"}},
		{"Devices", km.Devices, []string{"d"}},
		{"History", km.History, []string{"H"}},
		{"Artist", km.Artist, []string{"A"}},
		{"Like", km.Like, []string{"l"}},
		{"Lyrics", km.Lyrics, []string{"L"}},
		{"AddToPlaylist", km.AddToPlaylist, []string{"a"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bindingKeys := tt.binding.Keys()
			for _, expectedKey := range tt.keys {
				found := false
				for _, actualKey := range bindingKeys {
					if actualKey == expectedKey {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("%s binding missing key %q, has %v", tt.name, expectedKey, bindingKeys)
				}
			}
		})
	}
}

func TestShortHelp(t *testing.T) {
	km := DefaultKeyMap()
	help := km.ShortHelp()

	if len(help) == 0 {
		t.Error("ShortHelp() returned empty slice")
	}

	if len(help) > 10 {
		t.Errorf("ShortHelp() returned %d bindings, should be concise", len(help))
	}
}

func TestFullHelp(t *testing.T) {
	km := DefaultKeyMap()
	help := km.FullHelp()

	if len(help) == 0 {
		t.Error("FullHelp() returned empty slice")
	}

	totalBindings := 0
	for _, group := range help {
		totalBindings += len(group)
	}

	if totalBindings < 10 {
		t.Errorf("FullHelp() returned only %d bindings, expected more", totalBindings)
	}
}
