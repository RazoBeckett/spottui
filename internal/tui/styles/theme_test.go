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

	if theme.Primary.Dark == "" {
		t.Error("SpotifyTheme.Primary.Dark should not be empty")
	}
	if theme.Primary.Light == "" {
		t.Error("SpotifyTheme.Primary.Light should not be empty")
	}
}
