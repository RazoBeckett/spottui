package styles

import "github.com/charmbracelet/lipgloss"

// Theme defines the complete color scheme for the application
type Theme struct {
	// Primary colors
	Primary   lipgloss.AdaptiveColor
	Secondary lipgloss.AdaptiveColor
	Accent    lipgloss.AdaptiveColor

	// Semantic colors
	Success lipgloss.AdaptiveColor
	Warning lipgloss.AdaptiveColor
	Error   lipgloss.AdaptiveColor

	// UI colors
	Background lipgloss.AdaptiveColor
	Surface    lipgloss.AdaptiveColor
	Border     lipgloss.AdaptiveColor
	Text       lipgloss.AdaptiveColor
	TextMuted  lipgloss.AdaptiveColor
}

// SpotifyTheme is a Spotify-inspired color scheme
var SpotifyTheme = Theme{
	Primary:    lipgloss.AdaptiveColor{Light: "#1DB954", Dark: "#1DB954"}, // Spotify green
	Secondary:  lipgloss.AdaptiveColor{Light: "#191414", Dark: "#FFFFFF"},
	Accent:     lipgloss.AdaptiveColor{Light: "#1ED760", Dark: "#1ED760"},
	Success:    lipgloss.AdaptiveColor{Light: "#1DB954", Dark: "#1DB954"},
	Warning:    lipgloss.AdaptiveColor{Light: "#F59B23", Dark: "#F59B23"},
	Error:      lipgloss.AdaptiveColor{Light: "#E91429", Dark: "#E91429"},
	Background: lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#121212"},
	Surface:    lipgloss.AdaptiveColor{Light: "#F8F8F8", Dark: "#282828"},
	Border:     lipgloss.AdaptiveColor{Light: "#E0E0E0", Dark: "#404040"},
	Text:       lipgloss.AdaptiveColor{Light: "#191414", Dark: "#FFFFFF"},
	TextMuted:  lipgloss.AdaptiveColor{Light: "#6A6A6A", Dark: "#B3B3B3"},
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
	Error   lipgloss.Style
	Success lipgloss.Style
	Spinner lipgloss.Style
	Muted   lipgloss.Style
}

// NewStyles creates styled components from a theme
func NewStyles(theme Theme) Styles {
	return Styles{
		App: lipgloss.NewStyle().
			Background(theme.Background),

		Header: lipgloss.NewStyle().
			Foreground(theme.Primary).
			Bold(true).
			Padding(0, 1).
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
			PaddingLeft(2),

		ListItemActive: lipgloss.NewStyle().
			Foreground(theme.Primary).
			Bold(true).
			PaddingLeft(0),

		ListTitle: lipgloss.NewStyle().
			Foreground(theme.Primary).
			Bold(true).
			Padding(0, 1).
			MarginBottom(1),

		NowPlaying: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Primary).
			Padding(1, 2).
			MarginTop(1),

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
			Foreground(theme.TextMuted),
	}
}

// DefaultStyles returns styles with the Spotify theme
func DefaultStyles() Styles {
	return NewStyles(SpotifyTheme)
}
