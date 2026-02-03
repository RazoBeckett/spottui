package tui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/razobeckett/spottui/internal/tui/views"
)

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
		m.fetching = true
		m.fetchingDots = 0
		m.view = ViewDevices
		m.devicesData = nil
		return m, tea.Batch(m.fetchDevices(), m.scheduleFetchingTick())

	case key.Matches(msg, m.keys.GlobalSearch):
		m.prevView = m.view
		m.view = ViewSearch
		m.searchInput.Focus()
		m.searchTracksData = nil
		return m, textinput.Blink

	case key.Matches(msg, m.keys.History):
		m.prevView = m.view
		m.fetching = true
		m.fetchingDots = 0
		m.view = ViewHistory
		m.historyData = nil
		return m, tea.Batch(m.fetchRecentlyPlayed(), m.scheduleFetchingTick())

	case key.Matches(msg, m.keys.Lyrics):
		if m.playbackState != nil && m.playbackState.Item != nil {
			trackName := m.playbackState.Item.Name
			artistName := ""
			if len(m.playbackState.Item.Artists) > 0 {
				artistName = m.playbackState.Item.Artists[0].Name
			}
			m.prevView = m.view
			m.fetchingLyrics = true
			m.fetching = true
			m.fetchingDots = 0
			return m, tea.Batch(m.fetchLyrics(trackName, artistName), m.scheduleFetchingTick())
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
		if len(m.playlists.Items()) == 0 {
			return m, nil
		}
		switch item := m.playlists.SelectedItem().(type) {
		case views.PlaylistItem:
			m.selectedPlaylist = &item.Playlist
			m.fetching = true
			m.fetchingDots = 0
			m.view = ViewTracks
			m.tracksData = nil
			return m, tea.Batch(m.fetchTracks(item.Playlist.ID), m.scheduleFetchingTick())
		case views.LikedSongsItem:
			m.selectedPlaylist = nil
			m.fetching = true
			m.fetchingDots = 0
			m.view = ViewTracks
			m.tracksData = nil
			return m, tea.Batch(m.fetchLikedTracks(), m.scheduleFetchingTick())
		}
	}

	var cmd tea.Cmd
	m.playlists, cmd = m.playlists.Update(msg)
	return m, cmd
}

func (m Model) handleTrackKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Enter) {
		if len(m.tracks.Items()) > 0 {
			if item, ok := m.tracks.SelectedItem().(views.TrackItem); ok {
				return m, m.playTrack(item.Track)
			}
		}
	}

	if key.Matches(msg, m.keys.Artist) {
		if len(m.tracks.Items()) > 0 {
			if item, ok := m.tracks.SelectedItem().(views.TrackItem); ok {
				if len(item.Track.Track.Artists) > 0 {
					m.prevView = m.view
					m.fetching = true
					m.fetchingDots = 0
					m.view = ViewArtist
					m.artistTopTracksData = nil
					m.artistAlbumsData = nil
					return m, tea.Batch(m.fetchArtist(item.Track.Track.Artists[0].ID), m.scheduleFetchingTick())
				}
			}
		}
	}

	if key.Matches(msg, m.keys.Like) {
		if len(m.tracks.Items()) > 0 {
			if item, ok := m.tracks.SelectedItem().(views.TrackItem); ok {
				return m, m.toggleLikeTrack(item.Track.Track.ID, item.Track.Track.Name)
			}
		}
	}

	if key.Matches(msg, m.keys.AddToPlaylist) {
		if len(m.tracks.Items()) > 0 {
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
	}

	var cmd tea.Cmd
	m.tracks, cmd = m.tracks.Update(msg)
	return m, cmd
}

func (m Model) handleAlbumKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Enter) {
		if len(m.albumTracks.Items()) > 0 {
			if item, ok := m.albumTracks.SelectedItem().(views.AlbumTrackItem); ok {
				return m, m.playAlbumTrack(item.Track)
			}
		}
	}

	if key.Matches(msg, m.keys.Artist) {
		if len(m.albumTracks.Items()) > 0 {
			if item, ok := m.albumTracks.SelectedItem().(views.AlbumTrackItem); ok {
				if len(item.Track.Artists) > 0 {
					m.prevView = m.view
					m.fetching = true
					m.fetchingDots = 0
					m.view = ViewArtist
					m.artistTopTracksData = nil
					m.artistAlbumsData = nil
					return m, tea.Batch(m.fetchArtist(item.Track.Artists[0].ID), m.scheduleFetchingTick())
				}
			}
		}
	}

	if key.Matches(msg, m.keys.Like) {
		if len(m.albumTracks.Items()) > 0 {
			if item, ok := m.albumTracks.SelectedItem().(views.AlbumTrackItem); ok {
				return m, m.toggleLikeTrack(item.Track.ID, item.Track.Name)
			}
		}
	}

	if key.Matches(msg, m.keys.AddToPlaylist) {
		if len(m.albumTracks.Items()) > 0 {
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
	}

	var cmd tea.Cmd
	m.albumTracks, cmd = m.albumTracks.Update(msg)
	return m, cmd
}

func (m Model) handleDeviceKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Enter) {
		if len(m.devices.Items()) > 0 {
			if item, ok := m.devices.SelectedItem().(views.DeviceItem); ok {
				m.view = m.prevView
				return m, m.transferPlayback(item.Device.ID)
			}
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
			if len(m.searchResults.Items()) > 0 {
				if item, ok := m.searchResults.SelectedItem().(views.SearchItem); ok {
					switch item.Type {
					case views.SearchResultTrack:
						return m, m.playSearchItem(item)
					case views.SearchResultAlbum:
						m.selectedAlbum = item.Album
						m.prevView = m.view
						m.fetching = true
						m.fetchingDots = 0
						m.view = ViewAlbum
						m.albumTracksData = nil
						return m, tea.Batch(m.fetchAlbumTracks(item.Album.ID), m.scheduleFetchingTick())
					case views.SearchResultPlaylist:
						m.selectedPlaylist = item.Playlist
						m.prevView = m.view
						m.fetching = true
						m.fetchingDots = 0
						m.view = ViewTracks
						m.tracksData = nil
						return m, tea.Batch(m.fetchTracks(item.Playlist.ID), m.scheduleFetchingTick())
					case views.SearchResultArtist:
						m.prevView = m.view
						m.fetching = true
						m.fetchingDots = 0
						m.view = ViewArtist
						m.artistTopTracksData = nil
						m.artistAlbumsData = nil
						return m, tea.Batch(m.fetchArtist(item.Artist.ID), m.scheduleFetchingTick())
					}
				}
			}
		}
	}

	if key.Matches(msg, m.keys.Artist) {
		if m.hasSearchResults() {
			if len(m.searchResults.Items()) > 0 {
				if item, ok := m.searchResults.SelectedItem().(views.SearchItem); ok {
					if item.Type == views.SearchResultTrack && item.Track != nil && len(item.Track.Artists) > 0 {
						m.prevView = m.view
						m.fetching = true
						m.fetchingDots = 0
						m.view = ViewArtist
						m.artistTopTracksData = nil
						m.artistAlbumsData = nil
						return m, tea.Batch(m.fetchArtist(item.Track.Artists[0].ID), m.scheduleFetchingTick())
					}
				}
			}
		}
	}

	if key.Matches(msg, m.keys.Like) {
		if m.hasSearchResults() {
			if len(m.searchResults.Items()) > 0 {
				if item, ok := m.searchResults.SelectedItem().(views.SearchItem); ok {
					if item.Type == views.SearchResultTrack && item.Track != nil {
						return m, m.toggleLikeTrack(item.Track.ID, item.Track.Name)
					}
				}
			}
		}
	}

	if key.Matches(msg, m.keys.AddToPlaylist) {
		if m.hasSearchResults() {
			if len(m.searchResults.Items()) > 0 {
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
	}

	switch msg.String() {
	case "/":
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
		if len(m.historyTracks.Items()) > 0 {
			if item, ok := m.historyTracks.SelectedItem().(views.AlbumTrackItem); ok {
				return m, m.playHistoryTrack(item.Track)
			}
		}
	}

	if key.Matches(msg, m.keys.Artist) {
		if len(m.historyTracks.Items()) > 0 {
			if item, ok := m.historyTracks.SelectedItem().(views.AlbumTrackItem); ok {
				if len(item.Track.Artists) > 0 {
					m.prevView = m.view
					m.fetching = true
					m.fetchingDots = 0
					m.view = ViewArtist
					m.artistTopTracksData = nil
					m.artistAlbumsData = nil
					return m, tea.Batch(m.fetchArtist(item.Track.Artists[0].ID), m.scheduleFetchingTick())
				}
			}
		}
	}

	if key.Matches(msg, m.keys.Like) {
		if len(m.historyTracks.Items()) > 0 {
			if item, ok := m.historyTracks.SelectedItem().(views.AlbumTrackItem); ok {
				return m, m.toggleLikeTrack(item.Track.ID, item.Track.Name)
			}
		}
	}

	if key.Matches(msg, m.keys.AddToPlaylist) {
		if len(m.historyTracks.Items()) > 0 {
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
			if len(m.artistTopTracks.Items()) > 0 {
				if item, ok := m.artistTopTracks.SelectedItem().(views.ArtistTopTrackItem); ok {
					return m, m.playArtistTopTrack(item.Track)
				}
			}
		} else {
			if len(m.artistAlbums.Items()) > 0 {
				if item, ok := m.artistAlbums.SelectedItem().(views.ArtistAlbumItem); ok {
					m.selectedAlbum = &item.Album
					m.prevView = m.view
					m.fetching = true
					m.fetchingDots = 0
					m.view = ViewAlbum
					m.albumTracksData = nil
					return m, tea.Batch(m.fetchAlbumTracks(item.Album.ID), m.scheduleFetchingTick())
				}
			}
		}
	}

	if key.Matches(msg, m.keys.Like) {
		if m.artistViewMode == "tracks" {
			if len(m.artistTopTracks.Items()) > 0 {
				if item, ok := m.artistTopTracks.SelectedItem().(views.ArtistTopTrackItem); ok {
					return m, m.toggleLikeTrack(item.Track.ID, item.Track.Name)
				}
			}
		}
	}

	if key.Matches(msg, m.keys.AddToPlaylist) {
		if m.artistViewMode == "tracks" {
			if len(m.artistTopTracks.Items()) > 0 {
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
		if len(m.addToPlaylistList.Items()) > 0 {
			if item, ok := m.addToPlaylistList.SelectedItem().(views.PlaylistItem); ok {
				m.view = m.prevView
				return m, m.addTrackToPlaylist(item.Playlist.ID, item.Playlist.Name, m.addToPlaylistTrack, "")
			}
		}
	}

	var cmd tea.Cmd
	m.addToPlaylistList, cmd = m.addToPlaylistList.Update(msg)
	return m, cmd
}
