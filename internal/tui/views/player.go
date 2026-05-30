package views

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/zmb3/spotify/v2"

	"github.com/razobeckett/spottui/internal/tui/styles"
)

// RenderNowPlaying renders the now playing component. progressMs is the
// locally-tracked playback position used for an optimistic, smooth progress bar.
func RenderNowPlaying(state *spotify.PlayerState, progressMs int, s styles.Styles, width int) string {
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

	var shuffleIcon string
	if state.ShuffleState {
		shuffleIcon = s.ActiveIcon.Render("⤮ ")
	} else {
		shuffleIcon = s.MutedIcon.Render("⤮ ")
	}

	var repeatIcon string
	switch state.RepeatState {
	case "context":
		repeatIcon = s.ActiveIcon.Render("⟳ ")
	case "track":
		repeatIcon = s.ActiveIcon.Render("⟳₁")
	default:
		repeatIcon = s.MutedIcon.Render("⟳ ")
	}

	volumeStr := ""
	if state.Device.ID != "" {
		volumeStr = fmt.Sprintf("   %d%%", state.Device.Volume)
	}

	availableWidth := max(width-8, 30)
	progressWidth := max(availableWidth-14, 10)

	// Clamp the optimistic position to the track duration.
	posMs := progressMs
	if posMs < 0 {
		posMs = 0
	}
	if d := int(track.Duration); d > 0 && posMs > d {
		posMs = d
	}

	var progress float64
	if track.Duration > 0 {
		progress = float64(posMs) / float64(track.Duration)
	}
	progressBar := renderProgressBar(progress, progressWidth, s)

	currentTime := formatDuration(posMs)
	totalTime := formatDuration(int(track.Duration))
	timeStr := fmt.Sprintf("%s / %s", currentTime, totalTime)

	title := s.TrackTitle.Render(track.Name)
	artist := s.TrackArtist.Render(artistStr)

	rightContent := shuffleIcon + repeatIcon + s.MutedIcon.Render(volumeStr)
	rightWidth := lipgloss.Width(rightContent)

	leftContent := lipgloss.JoinHorizontal(lipgloss.Center,
		playIcon, " ", title, " - ", artist,
	)
	leftWidth := lipgloss.Width(leftContent)

	spacerWidth := max(availableWidth-leftWidth-rightWidth, 1)
	spacer := strings.Repeat(" ", spacerWidth)

	line1 := leftContent + spacer + rightContent

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
