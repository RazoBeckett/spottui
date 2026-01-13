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
		return s.NowPlaying.Width(max(width-4, 20)).Render(" Nothing playing")
	}

	track := state.Item

	artists := make([]string, len(track.Artists))
	for i, a := range track.Artists {
		artists[i] = a.Name
	}
	artistStr := strings.Join(artists, ", ")

	playIcon := ""
	if state.Playing {
		playIcon = ""
	}

	shuffleIcon := ""
	if state.ShuffleState {
		shuffleIcon = "  "
	}

	repeatIcon := ""
	switch state.RepeatState {
	case "context":
		repeatIcon = "  "
	case "track":
		repeatIcon = "  "
	}

	volumeStr := fmt.Sprintf(" %d%%", state.Device.Volume)

	availableWidth := max(width-8, 30)
	progressWidth := max(availableWidth-14, 10)
	progress := float64(state.Progress) / float64(track.Duration)
	progressBar := renderProgressBar(progress, progressWidth, s)

	currentTime := formatDuration(int(state.Progress))
	totalTime := formatDuration(int(track.Duration))
	timeStr := fmt.Sprintf("%s / %s", currentTime, totalTime)

	title := s.TrackTitle.Render(track.Name)
	artist := s.TrackArtist.Render(artistStr)

	leftContent := lipgloss.JoinHorizontal(lipgloss.Center,
		playIcon, " ", title, " - ", artist, shuffleIcon, repeatIcon,
	)

	leftWidth := lipgloss.Width(leftContent)
	volumeWidth := lipgloss.Width(volumeStr)
	spacerWidth := max(availableWidth-leftWidth-volumeWidth, 1)
	spacer := strings.Repeat(" ", spacerWidth)

	line1 := leftContent + spacer + s.Muted.Render(volumeStr)

	timeDisplay := s.PlaybackTime.Render(timeStr)
	timeWidth := lipgloss.Width(timeDisplay)
	barWidth := lipgloss.Width(progressBar)
	line2SpacerWidth := max(availableWidth-barWidth-timeWidth, 2)
	line2Spacer := strings.Repeat(" ", line2SpacerWidth)

	line2 := progressBar + line2Spacer + timeDisplay

	content := lipgloss.JoinVertical(lipgloss.Left, line1, line2)

	return s.NowPlaying.Width(max(width-4, 20)).Render(content)
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
