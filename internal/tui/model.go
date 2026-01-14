package tui

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
	ViewAlbum
	ViewArtist
	ViewDevices
	ViewSearch
	ViewHistory
	ViewHelp
	ViewLyrics
	ViewAddToPlaylist
)

const (
	MinWidth      = 60
	MinHeight     = 15
	HorizontalPad = 4
	ListHeightSub = 12
	MinListHeight = 5
)

type Model struct {
	// Core state
	view     View
	prevView View // For returning from overlays
	width    int
	height   int

	errMsg    string
	showError bool

	notifyMsg  string
	showNotify bool

	// Spotify client
	client *spotify.Client
	ctx    context.Context

	// Sub-models (embedded Bubble Tea components)
	spinner         spinner.Model
	playlists       list.Model
	tracks          list.Model
	albumTracks     list.Model
	artistTopTracks list.Model
	artistAlbums    list.Model
	devices         list.Model
	searchInput     textinput.Model
	searchResults   list.Model
	historyTracks   list.Model

	// Data
	currentUser         *spotify.PrivateUser
	playlistsData       []spotify.SimplePlaylist
	selectedPlaylist    *spotify.SimplePlaylist
	selectedAlbum       *spotify.SimpleAlbum
	selectedArtist      *spotify.FullArtist
	tracksData          []spotify.PlaylistTrack
	albumTracksData     []spotify.SimpleTrack
	artistTopTracksData []spotify.FullTrack
	artistAlbumsData    []spotify.SimpleAlbum
	playbackState       *spotify.PlayerState
	devicesData         []spotify.PlayerDevice
	searchTracksData    []spotify.FullTrack
	searchAlbumsData    []spotify.SimpleAlbum
	searchPlaylistsData []spotify.SimplePlaylist
	searchArtistsData   []spotify.FullArtist
	historyData         []spotify.RecentlyPlayedItem
	searching           bool
	artistViewMode      string
	lyricsData          string
	lyricsSynced        []SyncedLyricLine
	lyricsIsSynced      bool
	lyricsTrackName     string
	lyricsArtistName    string
	lyricsScrollOffset  int
	fetchingLyrics      bool

	localProgress  int
	lastProgressAt time.Time
	isPlaying      bool

	pendingSeek     int
	seekPending     bool
	lastSeekRequest time.Time

	addToPlaylistTrack spotify.ID
	addToPlaylistList  list.Model

	fetching     bool
	fetchingDots int

	// Styling and keybindings
	styles styles.Styles
	keys   KeyMap
	cache  *Cache
}

type ProgressTickMsg struct{}
type SeekTickMsg struct{}
type FetchingTickMsg struct{}

func NewModel(client *spotify.Client) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	ti := textinput.New()
	ti.Placeholder = "Search tracks..."
	ti.CharLimit = 100
	ti.Width = 40

	return Model{
		view:        ViewLoading,
		client:      client,
		ctx:         context.Background(),
		spinner:     s,
		searchInput: ti,
		styles:      styles.DefaultStyles(),
		keys:        DefaultKeyMap(),
		cache:       NewCache(),
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
		if m.view == ViewAlbum && m.albumTracks.FilterState() == list.Filtering {
			var cmd tea.Cmd
			m.albumTracks, cmd = m.albumTracks.Update(msg)
			return m, cmd
		}
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		listWidth := m.width - HorizontalPad
		listHeight := m.height - ListHeightSub
		if listHeight < MinListHeight {
			listHeight = MinListHeight
		}
		lists := []*list.Model{
			&m.playlists, &m.tracks, &m.albumTracks,
			&m.artistTopTracks, &m.artistAlbums, &m.devices,
			&m.searchResults, &m.historyTracks, &m.addToPlaylistList,
		}
		for _, l := range lists {
			if len(l.Items()) > 0 {
				l.SetSize(listWidth, listHeight)
			}
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case UserDataMsg:
		m.currentUser = msg.User
		m.playlistsData = msg.Playlists
		m.playlists = views.CreatePlaylistList(msg.Playlists, msg.LikedSongsTotal, m.styles, m.listWidth(), m.listHeight())
		m.view = ViewPlaylists
		return m, tea.Batch(m.pollPlaybackState(), m.schedulePlaybackPoll())

	case TracksLoadedMsg:
		m.fetching = false
		m.tracksData = msg.Tracks
		currentTrack := ""
		if m.playbackState != nil && m.playbackState.Item != nil {
			currentTrack = string(m.playbackState.Item.URI)
		}
		m.tracks = views.CreateTrackList(msg.Tracks, m.styles, currentTrack, m.listWidth(), m.listHeight())
		if m.selectedPlaylist != nil {
			m.tracks.Title = m.selectedPlaylist.Name
		} else {
			m.tracks.Title = "Liked Songs"
		}
		return m, nil

	case AlbumTracksLoadedMsg:
		m.fetching = false
		m.albumTracksData = msg.Tracks
		currentTrack := ""
		if m.playbackState != nil && m.playbackState.Item != nil {
			currentTrack = string(m.playbackState.Item.URI)
		}
		m.albumTracks = views.CreateAlbumTrackList(msg.Tracks, m.styles, currentTrack, m.listWidth(), m.listHeight())
		if m.selectedAlbum != nil {
			m.albumTracks.Title = m.selectedAlbum.Name
		}
		return m, nil

	case PlaybackStateMsg:
		m.playbackState = msg.State
		if msg.State != nil {
			m.localProgress = int(msg.State.Progress)
			m.lastProgressAt = time.Now()
			m.isPlaying = msg.State.Playing
		}
		if m.view == ViewTracks && m.playbackState != nil && m.playbackState.Item != nil {
			currentTrack := string(m.playbackState.Item.URI)
			delegate := views.TrackDelegate{Styles: m.styles, CurrentTrack: currentTrack}
			m.tracks.SetDelegate(delegate)
		}
		return m, nil

	case ProgressTickMsg:
		if m.isPlaying && m.view == ViewLyrics {
			elapsed := time.Since(m.lastProgressAt)
			m.localProgress += int(elapsed.Milliseconds())
			m.lastProgressAt = time.Now()
		}
		return m, m.scheduleProgressTick()

	case SeekTickMsg:
		if m.seekPending && time.Since(m.lastSeekRequest) >= 100*time.Millisecond {
			m.seekPending = false
			return m, m.executeSeek(m.pendingSeek)
		}
		if m.seekPending {
			return m, m.scheduleSeekTick()
		}
		return m, nil

	case FetchingTickMsg:
		if m.fetching {
			m.fetchingDots = (m.fetchingDots + 1) % 4
			return m, m.scheduleFetchingTick()
		}
		return m, nil

	case PollPlaybackMsg:
		return m, tea.Batch(m.pollPlaybackState(), m.schedulePlaybackPoll())

	case DevicesLoadedMsg:
		m.fetching = false
		m.devicesData = msg.Devices
		m.devices = views.CreateDeviceList(msg.Devices, m.styles, m.listWidth(), m.listHeight())
		return m, nil

	case SearchResultsMsg:
		m.fetching = false
		m.searching = false
		m.searchTracksData = msg.Tracks
		m.searchAlbumsData = msg.Albums
		m.searchPlaylistsData = msg.Playlists
		m.searchArtistsData = msg.Artists
		currentTrack := ""
		if m.playbackState != nil && m.playbackState.Item != nil {
			currentTrack = string(m.playbackState.Item.URI)
		}
		m.searchResults = views.CreateSearchResultsList(msg.Tracks, msg.Albums, msg.Playlists, msg.Artists, m.styles, currentTrack, m.listWidth(), m.listHeight())
		return m, nil

	case HistoryLoadedMsg:
		m.fetching = false
		m.historyData = msg.Items
		currentTrack := ""
		if m.playbackState != nil && m.playbackState.Item != nil {
			currentTrack = string(m.playbackState.Item.URI)
		}
		m.historyTracks = views.CreateHistoryList(msg.Items, m.styles, currentTrack, m.listWidth(), m.listHeight())
		return m, nil

	case ArtistLoadedMsg:
		m.fetching = false
		m.selectedArtist = msg.Artist
		m.artistTopTracksData = msg.TopTracks
		m.artistAlbumsData = msg.Albums
		m.artistViewMode = "tracks"
		currentTrack := ""
		if m.playbackState != nil && m.playbackState.Item != nil {
			currentTrack = string(m.playbackState.Item.URI)
		}
		m.artistTopTracks = views.CreateArtistTopTracksList(msg.TopTracks, m.styles, currentTrack, m.listWidth(), m.listHeight())
		m.artistAlbums = views.CreateArtistAlbumsList(msg.Albums, m.styles, m.listWidth(), m.listHeight())
		return m, nil

	case ErrMsg:
		m.fetching = false
		m.errMsg = msg.Err.Error()
		m.showError = true
		return m, m.scheduleErrorDismiss()

	case DismissErrorMsg:
		m.showError = false
		m.errMsg = ""
		return m, nil

	case VolumeChangedMsg:
		if m.playbackState != nil {
			m.playbackState.Device.Volume = spotify.Numeric(msg.Volume)
		}
		return m, nil

	case ShuffleToggledMsg:
		if m.playbackState != nil {
			m.playbackState.ShuffleState = msg.NewState
		}
		return m, nil

	case RepeatCycledMsg:
		if m.playbackState != nil {
			m.playbackState.RepeatState = msg.NewState
		}
		return m, nil

	case LyricsLoadedMsg:
		m.fetching = false
		m.fetchingLyrics = false
		m.lyricsTrackName = msg.TrackName
		m.lyricsArtistName = msg.ArtistName
		m.lyricsScrollOffset = 0

		if msg.SyncedLyrics != "" {
			m.lyricsSynced = parseLRC(msg.SyncedLyrics)
			m.lyricsIsSynced = len(m.lyricsSynced) > 0
			m.lyricsData = msg.Lyrics
		} else {
			m.lyricsSynced = nil
			m.lyricsIsSynced = false
			m.lyricsData = msg.Lyrics
		}

		m.view = ViewLyrics
		if m.lyricsIsSynced {
			return m, m.scheduleProgressTick()
		}
		return m, nil

	case LikeToggledMsg:
		action := "♥ Liked"
		if !msg.IsLiked {
			action = "♡ Unliked"
		}
		m.notifyMsg = action + ": " + msg.TrackName
		m.showNotify = true
		return m, m.scheduleNotifyDismiss()

	case DismissNotifyMsg:
		m.showNotify = false
		m.notifyMsg = ""
		return m, nil

	case TrackAddedToPlaylistMsg:
		m.notifyMsg = "Added to " + msg.PlaylistName
		m.showNotify = true
		return m, m.scheduleNotifyDismiss()
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
	if m.width < MinWidth || m.height < MinHeight {
		return m.renderTooSmall()
	}

	switch m.view {
	case ViewLoading:
		return m.renderLoading()
	case ViewPlaylists:
		return m.renderPlaylists()
	case ViewTracks:
		return m.renderTracks()
	case ViewAlbum:
		return m.renderAlbum()
	case ViewArtist:
		return m.renderArtist()
	case ViewDevices:
		return m.renderDevices()
	case ViewSearch:
		return m.renderSearch()
	case ViewHistory:
		return m.renderHistory()
	case ViewHelp:
		return m.renderHelp()
	case ViewLyrics:
		return m.renderLyrics()
	case ViewAddToPlaylist:
		return m.renderAddToPlaylist()
	default:
		return "Unknown view"
	}
}

func (m Model) listHeight() int {
	h := m.height - ListHeightSub
	if h < MinListHeight {
		return MinListHeight
	}
	return h
}

func (m Model) listWidth() int {
	return m.width - HorizontalPad
}
