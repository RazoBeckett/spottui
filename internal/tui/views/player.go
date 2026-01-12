package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/zmb3/spotify/v2"

	"github.com/razobeckett/spottui/internal/tui/styles"
)

// RenderNowPlaying renders the now playing component
func RenderNowPlaying(state *spotify.PlayerState, s styles.Styles, width int) string {
	if state == nil || state.Item == nil {
		return s.NowPlaying.Width(width - 4).Render("♫ Nothing playing")
	}

	track := state.Item

	// Build artist string
	artists := make([]string, len(track.Artists))
	for i, a := range track.Artists {
		artists[i] = a.Name
	}
	artistStr := strings.Join(artists, ", ")

	// Play/pause indicator
	playIcon := "▶"
	if state.Playing {
		playIcon = "󰏤"
	}

	// Shuffle indicator
	shuffleIcon := ""
	if state.ShuffleState {
		shuffleIcon = " 󰒝 "
	}

	// Repeat indicator
	repeatIcon := ""
	switch state.RepeatState {
	case "context":
		repeatIcon = " 󰑖 "
	case "track":
		repeatIcon = " 󰑘 "
	}

	// Volume indicator
	volumeStr := fmt.Sprintf("󰕾 %d%%", state.Device.Volume)

	// Progress bar
	progressWidth := width - 24
	if progressWidth < 10 {
		progressWidth = 10
	}
	progress := float64(state.Progress) / float64(track.Duration)
	progressBar := renderProgressBar(progress, progressWidth, s)

	// Time display
	currentTime := formatDuration(int(state.Progress))
	totalTime := formatDuration(int(track.Duration))
	timeStr := fmt.Sprintf("%s / %s", currentTime, totalTime)

	// Build the player UI
	title := s.TrackTitle.Render(track.Name)
	artist := s.TrackArtist.Render(artistStr)

	// Left side: play icon, track info, shuffle, repeat
	leftContent := lipgloss.JoinHorizontal(lipgloss.Center,
		playIcon, " ", title, " - ", artist, shuffleIcon, repeatIcon,
	)

	// Calculate spacing for right-aligned volume
	leftWidth := lipgloss.Width(leftContent)
	volumeWidth := lipgloss.Width(volumeStr)
	availableWidth := width - 8 // account for padding
	spacerWidth := availableWidth - leftWidth - volumeWidth
	if spacerWidth < 1 {
		spacerWidth = 1
	}
	spacer := strings.Repeat(" ", spacerWidth)

	line1 := leftContent + spacer + s.Muted.Render(volumeStr)

	timeDisplay := s.PlaybackTime.Render(timeStr)
	timeWidth := lipgloss.Width(timeDisplay)
	barWidth := lipgloss.Width(progressBar)
	line2SpacerWidth := availableWidth - barWidth - timeWidth
	if line2SpacerWidth < 2 {
		line2SpacerWidth = 2
	}
	line2Spacer := strings.Repeat(" ", line2SpacerWidth)

	line2 := progressBar + line2Spacer + timeDisplay

	content := lipgloss.JoinVertical(lipgloss.Left, line1, line2)

	return s.NowPlaying.Width(width - 4).Render(content)
}

func renderProgressBar(progress float64, width int, s styles.Styles) string {
	if width <= 0 {
		return ""
	}

	filled := int(float64(width) * progress)
	if filled > width {
		filled = width
	}
	empty := width - filled

	bar := strings.Repeat("━", filled) + strings.Repeat("─", empty)
	return s.ProgressBar.Render(bar)
}

func formatDuration(ms int) string {
	seconds := ms / 1000
	minutes := seconds / 60
	seconds = seconds % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
