package tui

import (
	"context"
	"time"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/zmb3/spotify/v2"

	"github.com/razobeckett/spottui/internal/config"
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
	ListHeightSub = 17
	MinListHeight = 5
)

type Model struct {
	// Grouped sub-state
	Nav      NavState      // active view + back-navigation stack
	UI       UIState       // transient presentation state
	Playback PlaybackState // now-playing snapshot + local progress/seek

	// Spotify client
	client *spotify.Client
	ctx    context.Context

	// Configuration
	cfg *config.Config

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
	devicesData         []spotify.PlayerDevice
	searchTracksData    []spotify.FullTrack
	searchAlbumsData    []spotify.SimpleAlbum
	searchPlaylistsData []spotify.SimplePlaylist
	searchArtistsData   []spotify.FullArtist
	historyData         []spotify.RecentlyPlayedItem
	artistViewMode      string
	lyricsData          string
	lyricsSynced        []SyncedLyricLine
	lyricsIsSynced      bool
	lyricsTrackName     string
	lyricsArtistName    string
	lyricsScrollOffset  int
	fetchingLyrics      bool

	addToPlaylistTrack spotify.ID
	addToPlaylistList  list.Model

	// Styling and keybindings
	styles styles.Styles
	keys   KeyMap
	cache  *Cache
}

type ProgressTickMsg struct{}
type SeekTickMsg struct{}
type FetchingTickMsg struct{}

func NewModel(client *spotify.Client, cfg *config.Config) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	ti := textinput.New()
	ti.Placeholder = "Search tracks..."
	ti.CharLimit = 100
	ti.SetWidth(40)

	return Model{
		Nav:         NavState{Current: ViewLoading},
		client:      client,
		ctx:         context.Background(),
		cfg:         cfg,
		spinner:     s,
		searchInput: ti,
		styles:      styles.DefaultStyles(),
		keys:        DefaultKeyMap(),
		cache:       NewCache(cfg),
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
	case tea.KeyPressMsg:
		if m.Nav.Current == ViewPlaylists && m.playlists.FilterState() == list.Filtering {
			var cmd tea.Cmd
			m.playlists, cmd = m.playlists.Update(msg)
			return m, cmd
		}
		if m.Nav.Current == ViewTracks && m.tracks.FilterState() == list.Filtering {
			var cmd tea.Cmd
			m.tracks, cmd = m.tracks.Update(msg)
			return m, cmd
		}
		if m.Nav.Current == ViewAlbum && m.albumTracks.FilterState() == list.Filtering {
			var cmd tea.Cmd
			m.albumTracks, cmd = m.albumTracks.Update(msg)
			return m, cmd
		}
		if m.Nav.Current == ViewHistory && m.historyTracks.FilterState() == list.Filtering {
			var cmd tea.Cmd
			m.historyTracks, cmd = m.historyTracks.Update(msg)
			return m, cmd
		}
		if m.Nav.Current == ViewArtist && m.artistAlbums.FilterState() == list.Filtering {
			var cmd tea.Cmd
			m.artistAlbums, cmd = m.artistAlbums.Update(msg)
			return m, cmd
		}
		if m.Nav.Current == ViewAddToPlaylist && m.addToPlaylistList.FilterState() == list.Filtering {
			var cmd tea.Cmd
			m.addToPlaylistList, cmd = m.addToPlaylistList.Update(msg)
			return m, cmd
		}
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.UI.Width = msg.Width
		m.UI.Height = msg.Height
		listWidth := m.UI.Width - HorizontalPad
		listHeight := m.UI.Height - ListHeightSub
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
		m.Nav.Current = ViewPlaylists
		return m, tea.Batch(m.pollPlaybackState(), m.schedulePlaybackPoll(), m.scheduleProgressTick())

	case TracksLoadedMsg:
		m.UI.Fetching = false
		m.tracksData = msg.Tracks
		currentTrack := ""
		if m.Playback.State != nil && m.Playback.State.Item != nil {
			currentTrack = string(m.Playback.State.Item.URI)
		}
		m.tracks = views.CreateTrackList(msg.Tracks, m.styles, currentTrack, m.listWidth(), m.listHeight())
		if m.selectedPlaylist != nil {
			m.tracks.Title = m.selectedPlaylist.Name
		} else {
			m.tracks.Title = "Liked Songs"
		}
		return m, nil

	case AlbumTracksLoadedMsg:
		m.UI.Fetching = false
		m.albumTracksData = msg.Tracks
		currentTrack := ""
		if m.Playback.State != nil && m.Playback.State.Item != nil {
			currentTrack = string(m.Playback.State.Item.URI)
		}
		m.albumTracks = views.CreateAlbumTrackList(msg.Tracks, m.styles, currentTrack, m.listWidth(), m.listHeight())
		if m.selectedAlbum != nil {
			m.albumTracks.Title = m.selectedAlbum.Name
		}
		return m, nil

	case PlaybackStateMsg:
		m.Playback.State = msg.State
		if msg.State != nil {
			// Don't clobber the optimistic position while a seek is in flight.
			if !m.Playback.SeekPending {
				m.Playback.LocalProgress = int(msg.State.Progress)
				m.Playback.LastProgressAt = time.Now()
			}
			m.Playback.IsPlaying = msg.State.Playing
		}
		if m.Nav.Current == ViewTracks && m.Playback.State != nil && m.Playback.State.Item != nil {
			currentTrack := string(m.Playback.State.Item.URI)
			delegate := views.TrackDelegate{Styles: m.styles, CurrentTrack: currentTrack}
			m.tracks.SetDelegate(delegate)
		}
		if m.Nav.Current == ViewLyrics && m.Playback.State != nil && m.Playback.State.Item != nil && !m.fetchingLyrics {
			newTrackName := m.Playback.State.Item.Name
			newArtistName := ""
			if len(m.Playback.State.Item.Artists) > 0 {
				newArtistName = m.Playback.State.Item.Artists[0].Name
			}
			if newTrackName != m.lyricsTrackName || newArtistName != m.lyricsArtistName {
				m.fetchingLyrics = true
				m.UI.Fetching = true
				return m, tea.Batch(m.fetchLyrics(newTrackName, newArtistName), m.scheduleFetchingTick())
			}
		}
		return m, nil

	case ProgressTickMsg:
		if m.Playback.IsPlaying {
			elapsed := time.Since(m.Playback.LastProgressAt)
			m.Playback.LocalProgress += int(elapsed.Milliseconds())
			m.Playback.LastProgressAt = time.Now()
		}
		return m, m.scheduleProgressTick()

	case SeekTickMsg:
		if m.Playback.SeekPending && time.Since(m.Playback.LastSeekRequest) >= 100*time.Millisecond {
			m.Playback.SeekPending = false
			return m, m.executeSeek(m.Playback.PendingSeek)
		}
		if m.Playback.SeekPending {
			return m, m.scheduleSeekTick()
		}
		return m, nil

	case FetchingTickMsg:
		if m.UI.Fetching {
			m.UI.FetchingDots = (m.UI.FetchingDots + 1) % 4
			return m, m.scheduleFetchingTick()
		}
		return m, nil

	case PollPlaybackMsg:
		return m, tea.Batch(m.pollPlaybackState(), m.schedulePlaybackPoll())

	case DevicesLoadedMsg:
		m.UI.Fetching = false
		m.devicesData = msg.Devices
		m.devices = views.CreateDeviceList(msg.Devices, m.styles, m.listWidth(), m.listHeight())
		return m, nil

	case SearchResultsMsg:
		m.UI.Fetching = false
		m.UI.Searching = false
		m.searchTracksData = msg.Tracks
		m.searchAlbumsData = msg.Albums
		m.searchPlaylistsData = msg.Playlists
		m.searchArtistsData = msg.Artists
		currentTrack := ""
		if m.Playback.State != nil && m.Playback.State.Item != nil {
			currentTrack = string(m.Playback.State.Item.URI)
		}
		m.searchResults = views.CreateSearchResultsList(msg.Tracks, msg.Albums, msg.Playlists, msg.Artists, m.styles, currentTrack, m.listWidth(), m.listHeight())
		return m, nil

	case HistoryLoadedMsg:
		m.UI.Fetching = false
		m.historyData = msg.Items
		currentTrack := ""
		if m.Playback.State != nil && m.Playback.State.Item != nil {
			currentTrack = string(m.Playback.State.Item.URI)
		}
		m.historyTracks = views.CreateHistoryList(msg.Items, m.styles, currentTrack, m.listWidth(), m.listHeight())
		return m, nil

	case ArtistLoadedMsg:
		m.UI.Fetching = false
		m.selectedArtist = msg.Artist
		m.artistTopTracksData = msg.TopTracks
		m.artistAlbumsData = msg.Albums
		m.artistViewMode = "tracks"
		currentTrack := ""
		if m.Playback.State != nil && m.Playback.State.Item != nil {
			currentTrack = string(m.Playback.State.Item.URI)
		}
		m.artistTopTracks = views.CreateArtistTopTracksList(msg.TopTracks, m.styles, currentTrack, m.listWidth(), m.listHeight())
		m.artistAlbums = views.CreateArtistAlbumsList(msg.Albums, m.styles, m.listWidth(), m.listHeight())
		return m, nil

	case ErrMsg:
		m.UI.Fetching = false
		m.UI.ErrMsg = msg.Err.Error()
		m.UI.ShowError = true
		return m, m.scheduleErrorDismiss()

	case DismissErrorMsg:
		m.UI.ShowError = false
		m.UI.ErrMsg = ""
		return m, nil

	case VolumeChangedMsg:
		if m.Playback.State != nil && m.Playback.State.Device.ID != "" {
			m.Playback.State.Device.Volume = spotify.Numeric(msg.Volume)
		}
		return m, nil

	case ShuffleToggledMsg:
		if m.Playback.State != nil {
			m.Playback.State.ShuffleState = msg.NewState
		}
		return m, nil

	case RepeatCycledMsg:
		if m.Playback.State != nil {
			m.Playback.State.RepeatState = msg.NewState
		}
		return m, nil

	case LyricsLoadedMsg:
		m.UI.Fetching = false
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

		if m.Nav.Current != ViewLyrics {
			m.Nav.Push(ViewLyrics)
		}
		return m, nil

	case LikeToggledMsg:
		action := "♥ Liked"
		if !msg.IsLiked {
			action = "♡ Unliked"
		}
		m.UI.NotifyMsg = action + ": " + msg.TrackName
		m.UI.ShowNotify = true
		return m, m.scheduleNotifyDismiss()

	case DismissNotifyMsg:
		m.UI.ShowNotify = false
		m.UI.NotifyMsg = ""
		return m, nil

	case TrackAddedToPlaylistMsg:
		m.UI.NotifyMsg = "Added to " + msg.PlaylistName
		m.UI.ShowNotify = true
		return m, m.scheduleNotifyDismiss()
	}

	// Delegate to active view's sub-model
	switch m.Nav.Current {
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
func (m Model) View() tea.View {
	v := tea.NewView(m.renderContent())
	v.AltScreen = true
	return v
}

func (m Model) renderContent() string {
	if m.UI.Width < MinWidth || m.UI.Height < MinHeight {
		return m.renderTooSmall()
	}

	switch m.Nav.Current {
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
	h := m.UI.Height - ListHeightSub
	if h < MinListHeight {
		return MinListHeight
	}
	return h
}

func (m Model) listWidth() int {
	return m.UI.Width - HorizontalPad
}
