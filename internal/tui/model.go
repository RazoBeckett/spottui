package tui

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zmb3/spotify/v2"

	"github.com/razobeckett/spottui/internal/tui/styles"
	"github.com/razobeckett/spottui/internal/tui/views"
)

// View represents the current screen
type View int

const (
	ViewLoading View = iota
	ViewPlaylists
	ViewTracks
	ViewHelp
)

// Model is the root application state
type Model struct {
	// Core state
	view     View
	prevView View // For returning from overlays
	width    int
	height   int
	err      error

	// Spotify client
	client *spotify.Client
	ctx    context.Context

	// Sub-models (embedded Bubble Tea components)
	spinner   spinner.Model
	playlists list.Model
	tracks    list.Model

	// Data
	currentUser      *spotify.PrivateUser
	playlistsData    []spotify.SimplePlaylist
	selectedPlaylist *spotify.SimplePlaylist
	tracksData       []spotify.PlaylistTrack
	playbackState    *spotify.PlayerState

	// Styling and keybindings
	styles styles.Styles
	keys   KeyMap
}

// NewModel creates the initial application model
func NewModel(client *spotify.Client) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return Model{
		view:    ViewLoading,
		client:  client,
		ctx:     context.Background(),
		spinner: s,
		styles:  styles.DefaultStyles(),
		keys:    DefaultKeyMap(),
	}
}

// Init runs initial commands (fetch user data, playlists)
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.fetchInitialData(),
	)
}

// Update handles all messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.view == ViewPlaylists && m.playlists.FilterState() == list.Filtering {
			var cmd tea.Cmd
			m.playlists, cmd = m.playlists.Update(msg)
			return m, cmd
		}
		if m.view == ViewTracks && m.tracks.FilterState() == list.Filtering {
			var cmd tea.Cmd
			m.tracks, cmd = m.tracks.Update(msg)
			return m, cmd
		}
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Update list dimensions
		listHeight := msg.Height - 12 // Leave room for header, player, help
		if listHeight < 5 {
			listHeight = 5
		}
		if len(m.playlists.Items()) > 0 || m.view == ViewPlaylists {
			m.playlists.SetSize(msg.Width-4, listHeight)
		}
		if len(m.tracks.Items()) > 0 || m.view == ViewTracks {
			m.tracks.SetSize(msg.Width-4, listHeight)
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case UserDataMsg:
		m.currentUser = msg.User
		m.playlistsData = msg.Playlists
		listHeight := m.height - 12
		if listHeight < 5 {
			listHeight = 5
		}
		m.playlists = views.CreatePlaylistList(msg.Playlists, m.styles, m.width-4, listHeight)
		m.view = ViewPlaylists
		return m, tea.Batch(m.pollPlaybackState(), m.schedulePlaybackPoll())

	case TracksLoadedMsg:
		m.tracksData = msg.Tracks
		listHeight := m.height - 12
		if listHeight < 5 {
			listHeight = 5
		}
		currentTrack := ""
		if m.playbackState != nil && m.playbackState.Item != nil {
			currentTrack = string(m.playbackState.Item.URI)
		}
		m.tracks = views.CreateTrackList(msg.Tracks, m.styles, currentTrack, m.width-4, listHeight)
		if m.selectedPlaylist != nil {
			m.tracks.Title = m.selectedPlaylist.Name
		}
		m.view = ViewTracks
		return m, nil

	case PlaybackStateMsg:
		m.playbackState = msg.State
		// Update track list delegate to show currently playing
		if m.view == ViewTracks && m.playbackState != nil && m.playbackState.Item != nil {
			currentTrack := string(m.playbackState.Item.URI)
			delegate := views.TrackDelegate{Styles: m.styles, CurrentTrack: currentTrack}
			m.tracks.SetDelegate(delegate)
		}
		return m, nil

	case PollPlaybackMsg:
		return m, tea.Batch(m.pollPlaybackState(), m.schedulePlaybackPoll())

	case ErrMsg:
		m.err = msg.Err
		return m, nil
	}

	// Delegate to active view's sub-model
	switch m.view {
	case ViewPlaylists:
		var cmd tea.Cmd
		m.playlists, cmd = m.playlists.Update(msg)
		cmds = append(cmds, cmd)

	case ViewTracks:
		var cmd tea.Cmd
		m.tracks, cmd = m.tracks.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the UI
func (m Model) View() string {
	if m.err != nil {
		return m.styles.Error.Render("Error: " + m.err.Error() + "\n\nPress q to quit.")
	}

	switch m.view {
	case ViewLoading:
		return m.renderLoading()
	case ViewPlaylists:
		return m.renderPlaylists()
	case ViewTracks:
		return m.renderTracks()
	case ViewHelp:
		return m.renderHelp()
	default:
		return "Unknown view"
	}
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global keybindings (work in all views)
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, m.keys.Help):
		if m.view == ViewHelp {
			m.view = m.prevView
		} else {
			m.prevView = m.view
			m.view = ViewHelp
		}
		return m, nil

	case key.Matches(msg, m.keys.Back):
		if m.view == ViewHelp {
			m.view = m.prevView
			return m, nil
		}
		if m.view == ViewTracks {
			m.view = ViewPlaylists
			return m, nil
		}

	// Playback controls (global)
	case key.Matches(msg, m.keys.PlayPause):
		return m, m.togglePlayback()

	case key.Matches(msg, m.keys.Next):
		return m, m.nextTrack()

	case key.Matches(msg, m.keys.Prev):
		return m, m.prevTrack()

	case key.Matches(msg, m.keys.VolumeUp):
		return m, m.volumeUp()

	case key.Matches(msg, m.keys.VolumeDown):
		return m, m.volumeDown()

	case key.Matches(msg, m.keys.Shuffle):
		return m, m.toggleShuffle()

	case key.Matches(msg, m.keys.Repeat):
		return m, m.cycleRepeat()

	case key.Matches(msg, m.keys.Refresh):
		return m, m.pollPlaybackState()
	}

	// View-specific keybindings
	switch m.view {
	case ViewPlaylists:
		return m.handlePlaylistKeys(msg)

	case ViewTracks:
		return m.handleTrackKeys(msg)
	}

	return m, nil
}

func (m Model) handlePlaylistKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Enter) {
		if item, ok := m.playlists.SelectedItem().(views.PlaylistItem); ok {
			m.selectedPlaylist = &item.Playlist
			return m, m.fetchTracks(item.Playlist.ID)
		}
	}

	var cmd tea.Cmd
	m.playlists, cmd = m.playlists.Update(msg)
	return m, cmd
}

func (m Model) handleTrackKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Enter) {
		if item, ok := m.tracks.SelectedItem().(views.TrackItem); ok {
			return m, m.playTrack(item.Track)
		}
	}

	var cmd tea.Cmd
	m.tracks, cmd = m.tracks.Update(msg)
	return m, cmd
}

// Render functions

func (m Model) renderLoading() string {
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		m.styles.Spinner.Render(m.spinner.View()+" Loading your Spotify data..."),
	)
}

func (m Model) renderPlaylists() string {
	header := m.renderHeader()
	content := m.playlists.View()
	player := views.RenderNowPlaying(m.playbackState, m.styles, m.width)
	help := m.renderHelpBar()

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		content,
		player,
		help,
	)
}

func (m Model) renderTracks() string {
	header := m.renderHeader()
	content := m.tracks.View()
	player := views.RenderNowPlaying(m.playbackState, m.styles, m.width)
	help := m.renderHelpBar()

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		content,
		player,
		help,
	)
}

func (m Model) renderHeader() string {
	userName := "Spotify User"
	if m.currentUser != nil && m.currentUser.DisplayName != "" {
		userName = m.currentUser.DisplayName
	}

	title := m.styles.Header.Render("♫ SpotTUI")
	user := m.styles.Muted.Render(" - " + userName)

	return title + user
}

func (m Model) renderHelpBar() string {
	return m.styles.HelpBar.Render(
		"↑/↓ navigate • enter select • esc back • space play/pause • n/p next/prev • +/- volume • ? help • q quit",
	)
}

func (m Model) renderHelp() string {
	help := `
╭─────────────────────────────────────╮
│           SpotTUI Help              │
├─────────────────────────────────────┤
│  Navigation                         │
│  ↑/k      Move up                   │
│  ↓/j      Move down                 │
│  enter    Select item               │
│  esc      Go back                   │
│  /        Filter list               │
│                                     │
│  Playback                           │
│  space    Play/Pause                │
│  n/>      Next track                │
│  p/<      Previous track            │
│  +/=      Volume up                 │
│  -        Volume down               │
│  s        Toggle shuffle            │
│  r        Cycle repeat mode         │
│                                     │
│  General                            │
│  ?        Toggle help               │
│  ctrl+r   Refresh                   │
│  q        Quit                      │
╰─────────────────────────────────────╯

Press ? or esc to close this help screen.
`
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		m.styles.Dialog.Render(help),
	)
}
