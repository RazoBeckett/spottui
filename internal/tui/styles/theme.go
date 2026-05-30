package styles

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// Theme defines the complete color scheme for the application
type Theme struct {
	// Primary colors
	Primary   color.Color
	Secondary color.Color
	Accent    color.Color

	// Semantic colors
	Success color.Color
	Warning color.Color
	Error   color.Color

	// UI colors
	Background color.Color
	Surface    color.Color
	Border     color.Color
	Text       color.Color
	TextMuted  color.Color
}

// SpotifyTheme is a Spotify-inspired color scheme (dark variant)
var SpotifyTheme = Theme{
	Primary:    lipgloss.Color("#1DB954"), // Spotify green
	Secondary:  lipgloss.Color("#FFFFFF"),
	Accent:     lipgloss.Color("#1ED760"),
	Success:    lipgloss.Color("#1DB954"),
	Warning:    lipgloss.Color("#F59B23"),
	Error:      lipgloss.Color("#E91429"),
	Background: lipgloss.Color("#121212"),
	Surface:    lipgloss.Color("#282828"),
	Border:     lipgloss.Color("#404040"),
	Text:       lipgloss.Color("#FFFFFF"),
	TextMuted:  lipgloss.Color("#B3B3B3"),
}

// Styles contains all UI component styles
type Styles struct {
	// App chrome
	App       lipgloss.Style
	Header    lipgloss.Style
	StatusBar lipgloss.Style
	HelpBar   lipgloss.Style

	// List components
	List           lipgloss.Style
	ListItem       lipgloss.Style
	ListItemActive lipgloss.Style
	ListTitle      lipgloss.Style

	// Player components
	NowPlaying   lipgloss.Style
	ProgressBar  lipgloss.Style
	TrackTitle   lipgloss.Style
	TrackArtist  lipgloss.Style
	PlaybackTime lipgloss.Style

	// Dialog/overlay
	Dialog       lipgloss.Style
	DialogTitle  lipgloss.Style
	DialogButton lipgloss.Style

	// Misc
	Error      lipgloss.Style
	Success    lipgloss.Style
	Spinner    lipgloss.Style
	Muted      lipgloss.Style
	ActiveIcon lipgloss.Style
	MutedIcon  lipgloss.Style

	// Gradient endpoints (primary -> accent) for logo and progress bar.
	GradStart color.Color
	GradEnd   color.Color
}

// NewStyles creates styled components from a theme
func NewStyles(theme Theme) Styles {
	return Styles{
		App: lipgloss.NewStyle().
			Background(theme.Background),

		Header: lipgloss.NewStyle().
			Foreground(theme.Primary).
			Bold(true).
			Padding(0, 2).
			MarginBottom(1),

		StatusBar: lipgloss.NewStyle().
			Foreground(theme.TextMuted).
			Background(theme.Surface).
			Padding(0, 1),

		HelpBar: lipgloss.NewStyle().
			Foreground(theme.TextMuted).
			Padding(0, 1),

		List: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Border).
			Padding(1, 2),

		ListItem: lipgloss.NewStyle().
			Foreground(theme.Text).
			Padding(0, 2),

		ListItemActive: lipgloss.NewStyle().
			Foreground(theme.Primary).
			Bold(true).
			PaddingLeft(2),

		ListTitle: lipgloss.NewStyle().
			Foreground(theme.Primary).
			Bold(true),

		NowPlaying: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Primary).
			Padding(1, 2),

		ProgressBar: lipgloss.NewStyle().
			Foreground(theme.Primary),

		TrackTitle: lipgloss.NewStyle().
			Foreground(theme.Text).
			Bold(true),

		TrackArtist: lipgloss.NewStyle().
			Foreground(theme.TextMuted),

		PlaybackTime: lipgloss.NewStyle().
			Foreground(theme.TextMuted).
			Width(12).
			Align(lipgloss.Right),

		Dialog: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Primary).
			Padding(1, 2).
			Width(60),

		DialogTitle: lipgloss.NewStyle().
			Foreground(theme.Primary).
			Bold(true).
			MarginBottom(1),

		DialogButton: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF")).
			Background(theme.Primary).
			Padding(0, 3).
			MarginRight(1),

		Error: lipgloss.NewStyle().
			Foreground(theme.Error).
			Bold(true),

		Success: lipgloss.NewStyle().
			Foreground(theme.Success),

		Spinner: lipgloss.NewStyle().
			Foreground(theme.Primary),

		Muted: lipgloss.NewStyle().
			Foreground(theme.TextMuted).
			Padding(0, 2),

		ActiveIcon: lipgloss.NewStyle().
			Foreground(theme.Primary),

		MutedIcon: lipgloss.NewStyle().
			Foreground(theme.TextMuted),

		GradStart: theme.Primary,
		GradEnd:   theme.Accent,
	}
}

// DefaultStyles returns styles with the Spotify theme
func DefaultStyles() Styles {
	return NewStyles(SpotifyTheme)
}
