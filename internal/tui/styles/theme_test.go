package styles

import (
	"testing"
)

func TestNewStyles(t *testing.T) {
	s := NewStyles(SpotifyTheme)

	if s.App.String() == "" && s.Header.String() == "" {
		t.Error("NewStyles() should create non-empty styles")
	}
}

func TestDefaultStyles(t *testing.T) {
	s := DefaultStyles()

	if s.App.String() == "" && s.Header.String() == "" {
		t.Error("DefaultStyles() should create non-empty styles")
	}
}

func TestSpotifyTheme(t *testing.T) {
	theme := SpotifyTheme

	if theme.Primary == nil {
		t.Error("SpotifyTheme.Primary should not be nil")
	}
	if theme.Background == nil {
		t.Error("SpotifyTheme.Background should not be nil")
	}
}
