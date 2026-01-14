package views

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zmb3/spotify/v2"

	"github.com/razobeckett/spottui/internal/tui/styles"
)

func createPlayerState(playing bool, progress, duration int, trackName, artistName string, volume int, shuffle bool, repeat string) *spotify.PlayerState {
	return &spotify.PlayerState{
		CurrentlyPlaying: spotify.CurrentlyPlaying{
			Playing:  playing,
			Progress: spotify.Numeric(progress),
			Item: &spotify.FullTrack{
				SimpleTrack: spotify.SimpleTrack{
					Name:     trackName,
					Duration: spotify.Numeric(duration),
					Artists:  []spotify.SimpleArtist{{Name: artistName}},
				},
			},
		},
		Device: spotify.PlayerDevice{
			Volume: spotify.Numeric(volume),
		},
		ShuffleState: shuffle,
		RepeatState:  repeat,
	}
}

func TestRenderNowPlaying_NilState(t *testing.T) {
	s := styles.DefaultStyles()
	result := RenderNowPlaying(nil, s, 80)
	assert.Contains(t, result, "Nothing playing")
}

func TestRenderNowPlaying_NilItem(t *testing.T) {
	s := styles.DefaultStyles()
	state := &spotify.PlayerState{}
	result := RenderNowPlaying(state, s, 80)
	assert.Contains(t, result, "Nothing playing")
}

func TestRenderNowPlaying_WithTrack(t *testing.T) {
	s := styles.DefaultStyles()
	state := createPlayerState(true, 60000, 180000, "Test Song", "Artist One", 75, true, "context")

	result := RenderNowPlaying(state, s, 100)

	assert.Contains(t, result, "Test Song")
	assert.Contains(t, result, "Artist One")
	assert.Contains(t, result, "75%")
	assert.Contains(t, result, "1:00")
	assert.Contains(t, result, "3:00")
}

func TestRenderNowPlaying_PlayPauseIcon(t *testing.T) {
	s := styles.DefaultStyles()

	playingState := createPlayerState(true, 0, 60000, "Track", "Artist", 50, false, "off")
	result := RenderNowPlaying(playingState, s, 80)
	assert.Contains(t, result, "")

	pausedState := createPlayerState(false, 0, 60000, "Track", "Artist", 50, false, "off")
	result = RenderNowPlaying(pausedState, s, 80)
	assert.Contains(t, result, "")
}

func TestRenderNowPlaying_RepeatStates(t *testing.T) {
	s := styles.DefaultStyles()

	offState := createPlayerState(false, 0, 60000, "Track", "Artist", 50, false, "off")
	result := RenderNowPlaying(offState, s, 80)
	assert.Contains(t, result, "⟳")

	contextState := createPlayerState(false, 0, 60000, "Track", "Artist", 50, false, "context")
	result = RenderNowPlaying(contextState, s, 80)
	assert.Contains(t, result, "⟳")

	trackState := createPlayerState(false, 0, 60000, "Track", "Artist", 50, false, "track")
	result = RenderNowPlaying(trackState, s, 80)
	assert.Contains(t, result, "⟳₁")
}

func TestRenderProgressBar(t *testing.T) {
	s := styles.DefaultStyles()

	result := renderProgressBar(0.5, 10, s)
	assert.Contains(t, result, "━")
	assert.Contains(t, result, "─")

	result = renderProgressBar(0, 10, s)
	filled := strings.Count(result, "━")
	assert.Equal(t, 0, filled)

	result = renderProgressBar(1.0, 10, s)
	empty := strings.Count(result, "─")
	assert.Equal(t, 0, empty)

	result = renderProgressBar(0.5, 0, s)
	assert.Equal(t, "", result)

	result = renderProgressBar(0.5, -5, s)
	assert.Equal(t, "", result)
}

func TestRenderProgressBar_Overflow(t *testing.T) {
	s := styles.DefaultStyles()
	result := renderProgressBar(1.5, 10, s)
	assert.NotContains(t, result, "─")
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		ms       int
		expected string
	}{
		{0, "0:00"},
		{1000, "0:01"},
		{60000, "1:00"},
		{61000, "1:01"},
		{90000, "1:30"},
		{180000, "3:00"},
		{3600000, "60:00"},
		{3661000, "61:01"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatDuration(tt.ms)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRenderNowPlaying_SmallWidth(t *testing.T) {
	s := styles.DefaultStyles()
	state := createPlayerState(false, 30000, 60000, "A Very Long Song Title", "Artist", 50, false, "off")

	result := RenderNowPlaying(state, s, 30)
	assert.NotEmpty(t, result)
}

func TestRenderNowPlaying_MultipleArtists(t *testing.T) {
	s := styles.DefaultStyles()
	state := &spotify.PlayerState{
		CurrentlyPlaying: spotify.CurrentlyPlaying{
			Playing:  true,
			Progress: 0,
			Item: &spotify.FullTrack{
				SimpleTrack: spotify.SimpleTrack{
					Name:     "Collab Song",
					Duration: 60000,
					Artists: []spotify.SimpleArtist{
						{Name: "Artist A"},
						{Name: "Artist B"},
						{Name: "Artist C"},
					},
				},
			},
		},
		Device: spotify.PlayerDevice{Volume: 50},
	}

	result := RenderNowPlaying(state, s, 100)
	assert.Contains(t, result, "Artist A")
	assert.Contains(t, result, "Artist B")
	assert.Contains(t, result, "Artist C")
}
