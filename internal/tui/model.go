package tui

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
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
	ViewAlbum
	ViewArtist
	ViewDevices
	ViewSearch
	ViewHistory
	ViewHelp
	ViewLyrics
	ViewAddToPlaylist
)

// Model is the root application state
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

	// Styling and keybindings
	styles styles.Styles
	keys   KeyMap
}

type ProgressTickMsg struct{}
type SeekTickMsg struct{}

type SyncedLyricLine struct {
	TimeMs int
	Text   string
}

var lrcRegex = regexp.MustCompile(`\[(\d{2}):(\d{2})\.(\d{2,3})\](.*)`)

func parseLRC(lrc string) []SyncedLyricLine {
	var lines []SyncedLyricLine
	for _, line := range strings.Split(lrc, "\n") {
		matches := lrcRegex.FindStringSubmatch(line)
		if len(matches) == 5 {
			min, _ := strconv.Atoi(matches[1])
			sec, _ := strconv.Atoi(matches[2])
			msStr := matches[3]
			ms, _ := strconv.Atoi(msStr)
			if len(msStr) == 2 {
				ms *= 10
			}
			timeMs := min*60*1000 + sec*1000 + ms
			text := strings.TrimSpace(matches[4])
			lines = append(lines, SyncedLyricLine{TimeMs: timeMs, Text: text})
		}
	}
	return lines
}

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
		m.playlists = views.CreatePlaylistList(msg.Playlists, msg.LikedSongsTotal, m.styles, m.width-4, listHeight)
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
		} else {
			m.tracks.Title = "Liked Songs"
		}
		m.view = ViewTracks
		return m, nil

	case AlbumTracksLoadedMsg:
		m.albumTracksData = msg.Tracks
		listHeight := m.height - 12
		if listHeight < 5 {
			listHeight = 5
		}
		currentTrack := ""
		if m.playbackState != nil && m.playbackState.Item != nil {
			currentTrack = string(m.playbackState.Item.URI)
		}
		m.albumTracks = views.CreateAlbumTrackList(msg.Tracks, m.styles, currentTrack, m.width-4, listHeight)
		if m.selectedAlbum != nil {
			m.albumTracks.Title = m.selectedAlbum.Name
		}
		m.view = ViewAlbum
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

	case PollPlaybackMsg:
		return m, tea.Batch(m.pollPlaybackState(), m.schedulePlaybackPoll())

	case DevicesLoadedMsg:
		m.devicesData = msg.Devices
		listHeight := m.height - 12
		if listHeight < 5 {
			listHeight = 5
		}
		m.devices = views.CreateDeviceList(msg.Devices, m.styles, m.width-4, listHeight)
		m.view = ViewDevices
		return m, nil

	case SearchResultsMsg:
		m.searching = false
		m.searchTracksData = msg.Tracks
		m.searchAlbumsData = msg.Albums
		m.searchPlaylistsData = msg.Playlists
		m.searchArtistsData = msg.Artists
		listHeight := m.height - 14
		if listHeight < 5 {
			listHeight = 5
		}
		currentTrack := ""
		if m.playbackState != nil && m.playbackState.Item != nil {
			currentTrack = string(m.playbackState.Item.URI)
		}
		m.searchResults = views.CreateSearchResultsList(msg.Tracks, msg.Albums, msg.Playlists, msg.Artists, m.styles, currentTrack, m.width-4, listHeight)
		return m, nil

	case HistoryLoadedMsg:
		m.historyData = msg.Items
		listHeight := m.height - 12
		if listHeight < 5 {
			listHeight = 5
		}
		currentTrack := ""
		if m.playbackState != nil && m.playbackState.Item != nil {
			currentTrack = string(m.playbackState.Item.URI)
		}
		m.historyTracks = views.CreateHistoryList(msg.Items, m.styles, currentTrack, m.width-4, listHeight)
		m.view = ViewHistory
		return m, nil

	case ArtistLoadedMsg:
		m.selectedArtist = msg.Artist
		m.artistTopTracksData = msg.TopTracks
		m.artistAlbumsData = msg.Albums
		m.artistViewMode = "tracks"
		listHeight := m.height - 14
		if listHeight < 5 {
			listHeight = 5
		}
		currentTrack := ""
		if m.playbackState != nil && m.playbackState.Item != nil {
			currentTrack = string(m.playbackState.Item.URI)
		}
		m.artistTopTracks = views.CreateArtistTopTracksList(msg.TopTracks, m.styles, currentTrack, m.width-4, listHeight)
		m.artistAlbums = views.CreateArtistAlbumsList(msg.Albums, m.styles, m.width-4, listHeight)
		m.view = ViewArtist
		return m, nil

	case ErrMsg:
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

	case LyricsLoadedMsg:
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

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.view == ViewSearch && m.searchInput.Focused() {
		return m.handleSearchKeys(msg)
	}

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
		if m.view == ViewAlbum {
			m.view = m.prevView
			return m, nil
		}
		if m.view == ViewArtist {
			m.view = m.prevView
			return m, nil
		}
		if m.view == ViewLyrics {
			m.view = m.prevView
			return m, nil
		}
		if m.view == ViewDevices {
			m.view = m.prevView
			return m, nil
		}
		if m.view == ViewSearch {
			if !m.searchInput.Focused() {
				m.view = m.prevView
				return m, nil
			}
		}
		if m.view == ViewHistory {
			m.view = m.prevView
			return m, nil
		}
		if m.view == ViewAddToPlaylist {
			m.view = m.prevView
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

	case key.Matches(msg, m.keys.SeekBackward):
		if !m.seekPending {
			m.pendingSeek = m.localProgress
		}
		m.pendingSeek -= 5000
		if m.pendingSeek < 0 {
			m.pendingSeek = 0
		}
		m.localProgress = m.pendingSeek
		m.seekPending = true
		m.lastSeekRequest = time.Now()
		return m, m.scheduleSeekTick()

	case key.Matches(msg, m.keys.SeekForward):
		if !m.seekPending {
			m.pendingSeek = m.localProgress
		}
		m.pendingSeek += 5000
		m.localProgress = m.pendingSeek
		m.seekPending = true
		m.lastSeekRequest = time.Now()
		return m, m.scheduleSeekTick()

	case key.Matches(msg, m.keys.Refresh):
		return m, m.pollPlaybackState()

	case key.Matches(msg, m.keys.Devices):
		m.prevView = m.view
		return m, m.fetchDevices()

	case key.Matches(msg, m.keys.GlobalSearch):
		m.prevView = m.view
		m.view = ViewSearch
		m.searchInput.Focus()
		m.searchTracksData = nil
		return m, textinput.Blink

	case key.Matches(msg, m.keys.History):
		m.prevView = m.view
		return m, m.fetchRecentlyPlayed()

	case key.Matches(msg, m.keys.Lyrics):
		if m.playbackState != nil && m.playbackState.Item != nil {
			trackName := m.playbackState.Item.Name
			artistName := ""
			if len(m.playbackState.Item.Artists) > 0 {
				artistName = m.playbackState.Item.Artists[0].Name
			}
			m.prevView = m.view
			m.fetchingLyrics = true
			return m, m.fetchLyrics(trackName, artistName)
		}
		return m, nil
	}

	// View-specific keybindings
	switch m.view {
	case ViewPlaylists:
		return m.handlePlaylistKeys(msg)

	case ViewTracks:
		return m.handleTrackKeys(msg)

	case ViewAlbum:
		return m.handleAlbumKeys(msg)

	case ViewArtist:
		return m.handleArtistKeys(msg)

	case ViewDevices:
		return m.handleDeviceKeys(msg)

	case ViewSearch:
		return m.handleSearchKeys(msg)

	case ViewHistory:
		return m.handleHistoryKeys(msg)

	case ViewLyrics:
		return m.handleLyricsKeys(msg)

	case ViewAddToPlaylist:
		return m.handleAddToPlaylistKeys(msg)
	}

	return m, nil
}

func (m Model) handlePlaylistKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Enter) {
		switch item := m.playlists.SelectedItem().(type) {
		case views.PlaylistItem:
			m.selectedPlaylist = &item.Playlist
			return m, m.fetchTracks(item.Playlist.ID)
		case views.LikedSongsItem:
			m.selectedPlaylist = nil
			return m, m.fetchLikedTracks()
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

	if key.Matches(msg, m.keys.Artist) {
		if item, ok := m.tracks.SelectedItem().(views.TrackItem); ok {
			if len(item.Track.Track.Artists) > 0 {
				m.prevView = m.view
				return m, m.fetchArtist(item.Track.Track.Artists[0].ID)
			}
		}
	}

	if key.Matches(msg, m.keys.Like) {
		if item, ok := m.tracks.SelectedItem().(views.TrackItem); ok {
			return m, m.toggleLikeTrack(item.Track.Track.ID, item.Track.Track.Name)
		}
	}

	if key.Matches(msg, m.keys.AddToPlaylist) {
		if item, ok := m.tracks.SelectedItem().(views.TrackItem); ok {
			m.addToPlaylistTrack = item.Track.Track.ID
			listHeight := m.height - 12
			if listHeight < 5 {
				listHeight = 5
			}
			m.addToPlaylistList = views.CreateAddToPlaylistList(m.playlistsData, m.styles, m.width-4, listHeight)
			m.prevView = m.view
			m.view = ViewAddToPlaylist
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.tracks, cmd = m.tracks.Update(msg)
	return m, cmd
}

func (m Model) handleAlbumKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Enter) {
		if item, ok := m.albumTracks.SelectedItem().(views.AlbumTrackItem); ok {
			return m, m.playAlbumTrack(item.Track)
		}
	}

	if key.Matches(msg, m.keys.Artist) {
		if item, ok := m.albumTracks.SelectedItem().(views.AlbumTrackItem); ok {
			if len(item.Track.Artists) > 0 {
				m.prevView = m.view
				return m, m.fetchArtist(item.Track.Artists[0].ID)
			}
		}
	}

	if key.Matches(msg, m.keys.Like) {
		if item, ok := m.albumTracks.SelectedItem().(views.AlbumTrackItem); ok {
			return m, m.toggleLikeTrack(item.Track.ID, item.Track.Name)
		}
	}

	if key.Matches(msg, m.keys.AddToPlaylist) {
		if item, ok := m.albumTracks.SelectedItem().(views.AlbumTrackItem); ok {
			m.addToPlaylistTrack = item.Track.ID
			listHeight := m.height - 12
			if listHeight < 5 {
				listHeight = 5
			}
			m.addToPlaylistList = views.CreateAddToPlaylistList(m.playlistsData, m.styles, m.width-4, listHeight)
			m.prevView = m.view
			m.view = ViewAddToPlaylist
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.albumTracks, cmd = m.albumTracks.Update(msg)
	return m, cmd
}

func (m Model) handleDeviceKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Enter) {
		if item, ok := m.devices.SelectedItem().(views.DeviceItem); ok {
			m.view = m.prevView
			return m, m.transferPlayback(item.Device.ID)
		}
	}

	var cmd tea.Cmd
	m.devices, cmd = m.devices.Update(msg)
	return m, cmd
}

func (m Model) handleSearchKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.searchInput.Focused() {
		switch msg.String() {
		case "esc", "q":
			m.searchInput.Blur()
			return m, nil
		case "enter":
			if m.searchInput.Value() != "" {
				m.searchInput.Blur()
				m.searching = true
				return m, m.searchTracks(m.searchInput.Value())
			}
			return m, nil
		}
		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)
		return m, cmd
	}

	if key.Matches(msg, m.keys.Enter) {
		if m.hasSearchResults() {
			if item, ok := m.searchResults.SelectedItem().(views.SearchItem); ok {
				switch item.Type {
				case views.SearchResultTrack:
					return m, m.playSearchItem(item)
				case views.SearchResultAlbum:
					m.selectedAlbum = item.Album
					m.prevView = m.view
					return m, m.fetchAlbumTracks(item.Album.ID)
				case views.SearchResultPlaylist:
					m.selectedPlaylist = item.Playlist
					m.prevView = m.view
					return m, m.fetchTracks(item.Playlist.ID)
				case views.SearchResultArtist:
					m.prevView = m.view
					return m, m.fetchArtist(item.Artist.ID)
				}
			}
		}
	}

	if key.Matches(msg, m.keys.Artist) {
		if m.hasSearchResults() {
			if item, ok := m.searchResults.SelectedItem().(views.SearchItem); ok {
				if item.Type == views.SearchResultTrack && item.Track != nil && len(item.Track.Artists) > 0 {
					m.prevView = m.view
					return m, m.fetchArtist(item.Track.Artists[0].ID)
				}
			}
		}
	}

	if key.Matches(msg, m.keys.Like) {
		if m.hasSearchResults() {
			if item, ok := m.searchResults.SelectedItem().(views.SearchItem); ok {
				if item.Type == views.SearchResultTrack && item.Track != nil {
					return m, m.toggleLikeTrack(item.Track.ID, item.Track.Name)
				}
			}
		}
	}

	if key.Matches(msg, m.keys.AddToPlaylist) {
		if m.hasSearchResults() {
			if item, ok := m.searchResults.SelectedItem().(views.SearchItem); ok {
				if item.Type == views.SearchResultTrack && item.Track != nil {
					m.addToPlaylistTrack = item.Track.ID
					listHeight := m.height - 12
					if listHeight < 5 {
						listHeight = 5
					}
					m.addToPlaylistList = views.CreateAddToPlaylistList(m.playlistsData, m.styles, m.width-4, listHeight)
					m.prevView = m.view
					m.view = ViewAddToPlaylist
					return m, nil
				}
			}
		}
	}

	switch msg.String() {
	case "/", "i":
		m.searchInput.Focus()
		return m, textinput.Blink
	}

	if m.hasSearchResults() {
		var cmd tea.Cmd
		m.searchResults, cmd = m.searchResults.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) hasSearchResults() bool {
	return len(m.searchTracksData) > 0 || len(m.searchAlbumsData) > 0 || len(m.searchPlaylistsData) > 0 || len(m.searchArtistsData) > 0
}

func (m Model) handleHistoryKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Enter) {
		if item, ok := m.historyTracks.SelectedItem().(views.AlbumTrackItem); ok {
			return m, m.playHistoryTrack(item.Track)
		}
	}

	if key.Matches(msg, m.keys.Artist) {
		if item, ok := m.historyTracks.SelectedItem().(views.AlbumTrackItem); ok {
			if len(item.Track.Artists) > 0 {
				m.prevView = m.view
				return m, m.fetchArtist(item.Track.Artists[0].ID)
			}
		}
	}

	if key.Matches(msg, m.keys.Like) {
		if item, ok := m.historyTracks.SelectedItem().(views.AlbumTrackItem); ok {
			return m, m.toggleLikeTrack(item.Track.ID, item.Track.Name)
		}
	}

	if key.Matches(msg, m.keys.AddToPlaylist) {
		if item, ok := m.historyTracks.SelectedItem().(views.AlbumTrackItem); ok {
			m.addToPlaylistTrack = item.Track.ID
			listHeight := m.height - 12
			if listHeight < 5 {
				listHeight = 5
			}
			m.addToPlaylistList = views.CreateAddToPlaylistList(m.playlistsData, m.styles, m.width-4, listHeight)
			m.prevView = m.view
			m.view = ViewAddToPlaylist
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.historyTracks, cmd = m.historyTracks.Update(msg)
	return m, cmd
}

func (m Model) handleLyricsKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	lines := strings.Split(m.lyricsData, "\n")
	visibleHeight := m.height - 6

	switch {
	case key.Matches(msg, m.keys.Up):
		if m.lyricsScrollOffset > 0 {
			m.lyricsScrollOffset--
		}
	case key.Matches(msg, m.keys.Down):
		if m.lyricsScrollOffset < len(lines)-visibleHeight {
			m.lyricsScrollOffset++
		}
	}

	return m, nil
}

func (m Model) handleArtistKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Back) {
		m.view = m.prevView
		return m, nil
	}

	if msg.String() == "tab" {
		if m.artistViewMode == "tracks" {
			m.artistViewMode = "albums"
		} else {
			m.artistViewMode = "tracks"
		}
		return m, nil
	}

	if key.Matches(msg, m.keys.Enter) {
		if m.artistViewMode == "tracks" {
			if item, ok := m.artistTopTracks.SelectedItem().(views.ArtistTopTrackItem); ok {
				return m, m.playArtistTopTrack(item.Track)
			}
		} else {
			if item, ok := m.artistAlbums.SelectedItem().(views.ArtistAlbumItem); ok {
				m.selectedAlbum = &item.Album
				m.prevView = m.view
				return m, m.fetchAlbumTracks(item.Album.ID)
			}
		}
	}

	if key.Matches(msg, m.keys.Like) {
		if m.artistViewMode == "tracks" {
			if item, ok := m.artistTopTracks.SelectedItem().(views.ArtistTopTrackItem); ok {
				return m, m.toggleLikeTrack(item.Track.ID, item.Track.Name)
			}
		}
	}

	if key.Matches(msg, m.keys.AddToPlaylist) {
		if m.artistViewMode == "tracks" {
			if item, ok := m.artistTopTracks.SelectedItem().(views.ArtistTopTrackItem); ok {
				m.addToPlaylistTrack = item.Track.ID
				listHeight := m.height - 12
				if listHeight < 5 {
					listHeight = 5
				}
				m.addToPlaylistList = views.CreateAddToPlaylistList(m.playlistsData, m.styles, m.width-4, listHeight)
				m.prevView = m.view
				m.view = ViewAddToPlaylist
				return m, nil
			}
		}
	}

	if m.artistViewMode == "tracks" {
		var cmd tea.Cmd
		m.artistTopTracks, cmd = m.artistTopTracks.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.artistAlbums, cmd = m.artistAlbums.Update(msg)
	return m, cmd
}

func (m Model) handleAddToPlaylistKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Enter) {
		if item, ok := m.addToPlaylistList.SelectedItem().(views.PlaylistItem); ok {
			m.view = m.prevView
			return m, m.addTrackToPlaylist(item.Playlist.ID, item.Playlist.Name, m.addToPlaylistTrack, "")
		}
	}

	var cmd tea.Cmd
	m.addToPlaylistList, cmd = m.addToPlaylistList.Update(msg)
	return m, cmd
}

func (m Model) renderLoading() string {
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		m.styles.Spinner.Render(m.spinner.View()+" Loading your Spotify data..."),
	)
}

func (m Model) renderPlaylists() string {
	header := m.renderHeader()
	notification := m.renderNotification()
	player := views.RenderNowPlaying(m.playbackState, m.styles, m.width)
	help := m.renderHelpBar()

	headerHeight := lipgloss.Height(header)
	notificationHeight := lipgloss.Height(notification)
	playerHeight := lipgloss.Height(player)
	helpHeight := lipgloss.Height(help)
	contentHeight := m.height - headerHeight - notificationHeight - playerHeight - helpHeight

	if contentHeight < 1 {
		contentHeight = 1
	}

	content := lipgloss.NewStyle().Height(contentHeight).Render(m.playlists.View())

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		content,
		notification,
		player,
		help,
	)
}

func (m Model) renderTracks() string {
	header := m.renderHeader()
	notification := m.renderNotification()
	player := views.RenderNowPlaying(m.playbackState, m.styles, m.width)
	help := m.renderHelpBar()

	headerHeight := lipgloss.Height(header)
	notificationHeight := lipgloss.Height(notification)
	playerHeight := lipgloss.Height(player)
	helpHeight := lipgloss.Height(help)
	contentHeight := m.height - headerHeight - notificationHeight - playerHeight - helpHeight

	if contentHeight < 1 {
		contentHeight = 1
	}

	content := lipgloss.NewStyle().Height(contentHeight).Render(m.tracks.View())

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		content,
		notification,
		player,
		help,
	)
}

func (m Model) renderAlbum() string {
	header := m.renderHeader()
	notification := m.renderNotification()
	player := views.RenderNowPlaying(m.playbackState, m.styles, m.width)
	help := m.renderHelpBar()

	headerHeight := lipgloss.Height(header)
	notificationHeight := lipgloss.Height(notification)
	playerHeight := lipgloss.Height(player)
	helpHeight := lipgloss.Height(help)
	contentHeight := m.height - headerHeight - notificationHeight - playerHeight - helpHeight

	if contentHeight < 1 {
		contentHeight = 1
	}

	content := lipgloss.NewStyle().Height(contentHeight).Render(m.albumTracks.View())

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		content,
		notification,
		player,
		help,
	)
}

func (m Model) renderDevices() string {
	header := m.renderHeader()
	help := m.styles.HelpBar.Render("↑/↓ navigate • enter select • esc back")

	headerHeight := lipgloss.Height(header)
	helpHeight := lipgloss.Height(help)
	contentHeight := m.height - headerHeight - helpHeight

	if contentHeight < 1 {
		contentHeight = 1
	}

	content := lipgloss.NewStyle().Height(contentHeight).Render(m.devices.View())

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		content,
		help,
	)
}

func (m Model) renderSearch() string {
	header := m.renderHeader()

	inputStyle := m.styles.Header.Copy().Padding(0, 1)
	searchBox := inputStyle.Render("🔍 " + m.searchInput.View())

	notification := m.renderNotification()
	player := views.RenderNowPlaying(m.playbackState, m.styles, m.width)
	help := m.styles.HelpBar.Render("enter search • ↑/↓ navigate results • esc back")

	headerHeight := lipgloss.Height(header)
	searchBoxHeight := lipgloss.Height(searchBox)
	notificationHeight := lipgloss.Height(notification)
	playerHeight := lipgloss.Height(player)
	helpHeight := lipgloss.Height(help)
	contentHeight := m.height - headerHeight - searchBoxHeight - notificationHeight - playerHeight - helpHeight

	if contentHeight < 1 {
		contentHeight = 1
	}

	var contentView string
	if m.hasSearchResults() {
		contentView = m.searchResults.View()
	} else if m.searching {
		contentView = m.styles.Muted.Render("\n  Searching...")
	} else if m.searchInput.Value() != "" && !m.searchInput.Focused() {
		contentView = m.styles.Muted.Render("\n  No results found")
	} else {
		contentView = m.styles.Muted.Render("\n  Type your query and press Enter to search...")
	}

	content := lipgloss.NewStyle().Height(contentHeight).Render(contentView)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		searchBox,
		content,
		notification,
		player,
		help,
	)
}

func (m Model) renderHistory() string {
	header := m.renderHeader()
	notification := m.renderNotification()
	player := views.RenderNowPlaying(m.playbackState, m.styles, m.width)
	help := m.renderHelpBar()

	headerHeight := lipgloss.Height(header)
	notificationHeight := lipgloss.Height(notification)
	playerHeight := lipgloss.Height(player)
	helpHeight := lipgloss.Height(help)
	contentHeight := m.height - headerHeight - notificationHeight - playerHeight - helpHeight

	if contentHeight < 1 {
		contentHeight = 1
	}

	content := lipgloss.NewStyle().Height(contentHeight).Render(m.historyTracks.View())

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		content,
		notification,
		player,
		help,
	)
}

func (m Model) renderArtist() string {
	header := m.renderHeader()

	artistName := "Unknown Artist"
	artistInfo := ""
	if m.selectedArtist != nil {
		artistName = m.selectedArtist.Name
		if len(m.selectedArtist.Genres) > 0 {
			genres := m.selectedArtist.Genres
			if len(genres) > 3 {
				genres = genres[:3]
			}
			artistInfo = " • " + strings.Join(genres, ", ")
		}
	}

	artistHeader := m.styles.ListTitle.Render("󰠃 " + artistName + artistInfo)

	tabStyle := m.styles.Muted
	activeTabStyle := m.styles.ListItemActive
	tracksTab := "Top Tracks"
	albumsTab := "Albums"
	if m.artistViewMode == "tracks" {
		tracksTab = activeTabStyle.Render("[" + tracksTab + "]")
		albumsTab = tabStyle.Render(" " + albumsTab + " ")
	} else {
		tracksTab = tabStyle.Render(" " + tracksTab + " ")
		albumsTab = activeTabStyle.Render("[" + albumsTab + "]")
	}
	tabs := "  " + tracksTab + "  " + albumsTab

	notification := m.renderNotification()
	player := views.RenderNowPlaying(m.playbackState, m.styles, m.width)
	help := m.styles.HelpBar.Render("tab switch view • enter play/select • a view artist • esc back")

	headerHeight := lipgloss.Height(header)
	artistHeaderHeight := lipgloss.Height(artistHeader)
	tabsHeight := 1
	notificationHeight := lipgloss.Height(notification)
	playerHeight := lipgloss.Height(player)
	helpHeight := lipgloss.Height(help)
	contentHeight := m.height - headerHeight - artistHeaderHeight - tabsHeight - notificationHeight - playerHeight - helpHeight

	if contentHeight < 1 {
		contentHeight = 1
	}

	var listView string
	if m.artistViewMode == "tracks" {
		listView = m.artistTopTracks.View()
	} else {
		listView = m.artistAlbums.View()
	}

	content := lipgloss.NewStyle().Height(contentHeight).Render(listView)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		artistHeader,
		tabs,
		content,
		notification,
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
		"↑/↓ navigate • enter select • esc back • space play/pause • ? help",
	)
}

func (m Model) renderNotification() string {
	if m.showError {
		return m.styles.Error.Render("⚠ " + m.errMsg)
	}
	if m.showNotify {
		return m.styles.Success.Render(m.notifyMsg)
	}
	return ""
}

func (m Model) renderHelp() string {
	title := m.styles.DialogTitle.Render("SpotTUI Help")

	navSection := m.styles.ListItemActive.Render("Navigation") + "\n" +
		"  ↑/k       Move up\n" +
		"  ↓/j       Move down\n" +
		"  enter     Select item\n" +
		"  esc       Go back\n" +
		"  /         Filter list\n" +
		"  S         Global search\n" +
		"  A         View artist"

	playbackSection := m.styles.ListItemActive.Render("Playback") + "\n" +
		"  space     Play/Pause\n" +
		"  n/>       Next track\n" +
		"  p/<       Previous track\n" +
		"  [/]       Seek -/+5s\n" +
		"  +/=       Volume up\n" +
		"  -         Volume down\n" +
		"  s         Toggle shuffle\n" +
		"  r         Cycle repeat mode\n" +
		"  l         Like/unlike track\n" +
		"  L         Show lyrics\n" +
		"  a         Add to playlist"

	generalSection := m.styles.ListItemActive.Render("General") + "\n" +
		"  H         Recently played\n" +
		"  d         Device selector\n" +
		"  ?         Toggle help\n" +
		"  ctrl+r    Refresh\n" +
		"  q         Quit"

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		navSection,
		"",
		playbackSection,
		"",
		generalSection,
	)

	footer := m.styles.Muted.Render("Press ? or esc to close")

	dialog := m.styles.Dialog.Render(content)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, dialog, "", footer),
	)
}

func (m Model) renderLyrics() string {
	player := views.RenderNowPlaying(m.playbackState, m.styles, m.width)
	help := m.styles.HelpBar.Render("↑/↓ scroll • esc back")

	playerHeight := lipgloss.Height(player)
	helpHeight := lipgloss.Height(help)

	noLyricsAvailable := m.lyricsData == "" && len(m.lyricsSynced) == 0
	if noLyricsAvailable {
		noLyrics := m.styles.Muted.Render("No lyrics available for this track")
		trackInfo := ""
		if m.lyricsTrackName != "" {
			trackInfo = m.styles.ListTitle.Render(m.lyricsTrackName+" - "+m.lyricsArtistName) + "\n\n"
		}
		content := trackInfo + noLyrics
		contentHeight := m.height - playerHeight - helpHeight
		centeredContent := lipgloss.Place(
			m.width, contentHeight,
			lipgloss.Center, lipgloss.Center,
			content,
		)
		return lipgloss.JoinVertical(lipgloss.Left,
			centeredContent,
			player,
			help,
		)
	}

	header := m.styles.ListTitle.Render("♫ " + m.lyricsTrackName + " - " + m.lyricsArtistName)
	header = lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(header)

	headerHeight := lipgloss.Height(header)
	lyricsAreaHeight := m.height - headerHeight - playerHeight - helpHeight - 2

	if lyricsAreaHeight < 1 {
		lyricsAreaHeight = 1
	}

	if m.lyricsIsSynced && len(m.lyricsSynced) > 0 {
		return m.renderSyncedLyrics(header, player, help, lyricsAreaHeight)
	}

	return m.renderPlainLyrics(header, player, help, lyricsAreaHeight)
}

func (m Model) renderPlainLyrics(header, player, help string, lyricsAreaHeight int) string {
	lines := strings.Split(m.lyricsData, "\n")

	startLine := m.lyricsScrollOffset
	if startLine > len(lines)-lyricsAreaHeight {
		startLine = len(lines) - lyricsAreaHeight
	}
	if startLine < 0 {
		startLine = 0
	}

	endLine := startLine + lyricsAreaHeight
	if endLine > len(lines) {
		endLine = len(lines)
	}

	visibleLines := lines[startLine:endLine]

	var styledLines []string
	for _, line := range visibleLines {
		centered := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(line)
		styledLines = append(styledLines, centered)
	}

	lyricsContent := strings.Join(styledLines, "\n")
	lyricsBox := lipgloss.NewStyle().Height(lyricsAreaHeight).Render(lyricsContent)

	return lipgloss.JoinVertical(lipgloss.Left,
		"",
		header,
		"",
		lyricsBox,
		player,
		help,
	)
}

func (m Model) renderSyncedLyrics(header, player, help string, lyricsAreaHeight int) string {
	currentTimeMs := m.localProgress

	currentLineIdx := 0
	for i, line := range m.lyricsSynced {
		if line.TimeMs <= currentTimeMs {
			currentLineIdx = i
		} else {
			break
		}
	}

	centerOffset := lyricsAreaHeight / 2
	startLine := currentLineIdx - centerOffset
	if startLine < 0 {
		startLine = 0
	}

	endLine := startLine + lyricsAreaHeight
	if endLine > len(m.lyricsSynced) {
		endLine = len(m.lyricsSynced)
		startLine = endLine - lyricsAreaHeight
		if startLine < 0 {
			startLine = 0
		}
	}

	activeStyle := lipgloss.NewStyle().
		Foreground(m.styles.ListItemActive.GetForeground()).
		Bold(true)
	mutedStyle := m.styles.Muted

	var styledLines []string
	for i := startLine; i < endLine; i++ {
		line := m.lyricsSynced[i]
		var styled string
		if i == currentLineIdx {
			styled = activeStyle.Render(line.Text)
		} else {
			styled = mutedStyle.Render(line.Text)
		}
		centered := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(styled)
		styledLines = append(styledLines, centered)
	}

	for len(styledLines) < lyricsAreaHeight {
		styledLines = append(styledLines, "")
	}

	lyricsContent := strings.Join(styledLines, "\n")
	lyricsBox := lipgloss.NewStyle().Height(lyricsAreaHeight).Render(lyricsContent)

	return lipgloss.JoinVertical(lipgloss.Left,
		"",
		header,
		"",
		lyricsBox,
		player,
		help,
	)
}

func (m Model) renderAddToPlaylist() string {
	header := m.styles.Header.Render("Add to Playlist")
	notification := m.renderNotification()
	player := views.RenderNowPlaying(m.playbackState, m.styles, m.width)
	help := m.renderHelpBar()

	headerHeight := lipgloss.Height(header)
	notificationHeight := lipgloss.Height(notification)
	playerHeight := lipgloss.Height(player)
	helpHeight := lipgloss.Height(help)
	availableHeight := m.height - headerHeight - notificationHeight - playerHeight - helpHeight - 2

	if availableHeight < 5 {
		availableHeight = 5
	}

	m.addToPlaylistList.SetSize(m.width-4, availableHeight)

	var content strings.Builder
	content.WriteString(header)
	content.WriteString("\n")
	content.WriteString(m.addToPlaylistList.View())

	if notification != "" {
		return lipgloss.JoinVertical(lipgloss.Left,
			content.String(),
			notification,
			player,
			help,
		)
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		content.String(),
		player,
		help,
	)
}
