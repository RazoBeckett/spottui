package tui

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zmb3/spotify/v2"

	"github.com/razobeckett/spottui/internal/tui/views"
)

// Message types for async results

// UserDataMsg contains initial user data loaded at startup
type UserDataMsg struct {
	User            *spotify.PrivateUser
	Playlists       []spotify.SimplePlaylist
	LikedSongsTotal int
}

// TracksLoadedMsg contains tracks for a playlist
type TracksLoadedMsg struct {
	Tracks []spotify.PlaylistTrack
}

// PlaybackStateMsg contains current playback state
type PlaybackStateMsg struct {
	State *spotify.PlayerState
}

// PollPlaybackMsg triggers a playback state refresh
type PollPlaybackMsg struct{}

type DevicesLoadedMsg struct {
	Devices []spotify.PlayerDevice
}

type AlbumTracksLoadedMsg struct {
	Tracks []spotify.SimpleTrack
}

type HistoryLoadedMsg struct {
	Items []spotify.RecentlyPlayedItem
}

type ArtistLoadedMsg struct {
	Artist    *spotify.FullArtist
	TopTracks []spotify.FullTrack
	Albums    []spotify.SimpleAlbum
}

// ErrMsg contains an error from async operations
type ErrMsg struct {
	Err error
}

type DismissErrorMsg struct{}

type DismissNotifyMsg struct{}

type VolumeChangedMsg struct {
	Volume int
}

type TrackAddedToPlaylistMsg struct {
	PlaylistName string
	TrackName    string
}

func (m Model) scheduleErrorDismiss() tea.Cmd {
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return DismissErrorMsg{}
	})
}

func (m Model) scheduleNotifyDismiss() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return DismissNotifyMsg{}
	})
}

func (m Model) fetchInitialData() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(m.ctx, 15*time.Second)
		defer cancel()

		user, err := m.client.CurrentUser(ctx)
		if err != nil {
			return ErrMsg{Err: err}
		}

		var allPlaylists []spotify.SimplePlaylist
		limit := 50
		offset := 0

		for {
			playlists, err := m.client.CurrentUsersPlaylists(ctx,
				spotify.Limit(limit), spotify.Offset(offset))
			if err != nil {
				return ErrMsg{Err: err}
			}

			allPlaylists = append(allPlaylists, playlists.Playlists...)

			if len(playlists.Playlists) < limit {
				break
			}
			offset += limit
		}

		likedSongs, err := m.client.CurrentUsersTracks(ctx, spotify.Limit(1))
		likedTotal := 0
		if err == nil {
			likedTotal = int(likedSongs.Total)
		}

		return UserDataMsg{
			User:            user,
			Playlists:       allPlaylists,
			LikedSongsTotal: likedTotal,
		}
	}
}

func (m Model) fetchTracks(playlistID spotify.ID) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(m.ctx, 15*time.Second)
		defer cancel()

		var allTracks []spotify.PlaylistTrack
		limit := 100
		offset := 0

		for {
			tracks, err := m.client.GetPlaylistTracks(ctx, playlistID,
				spotify.Limit(limit), spotify.Offset(offset))
			if err != nil {
				return ErrMsg{Err: err}
			}

			allTracks = append(allTracks, tracks.Tracks...)

			if len(tracks.Tracks) < limit {
				break
			}
			offset += limit
		}

		return TracksLoadedMsg{Tracks: allTracks}
	}
}

func (m Model) fetchLikedTracks() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(m.ctx, 15*time.Second)
		defer cancel()

		var allTracks []spotify.PlaylistTrack
		limit := 50
		offset := 0

		for {
			saved, err := m.client.CurrentUsersTracks(ctx,
				spotify.Limit(limit), spotify.Offset(offset))
			if err != nil {
				return ErrMsg{Err: err}
			}

			for _, s := range saved.Tracks {
				allTracks = append(allTracks, spotify.PlaylistTrack{
					Track: s.FullTrack,
				})
			}

			if len(saved.Tracks) < limit {
				break
			}
			offset += limit
		}

		return TracksLoadedMsg{Tracks: allTracks}
	}
}

func (m Model) fetchAlbumTracks(albumID spotify.ID) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(m.ctx, 15*time.Second)
		defer cancel()

		var allTracks []spotify.SimpleTrack
		limit := 50
		offset := 0

		for {
			tracks, err := m.client.GetAlbumTracks(ctx, albumID,
				spotify.Limit(limit), spotify.Offset(offset))
			if err != nil {
				return ErrMsg{Err: err}
			}

			allTracks = append(allTracks, tracks.Tracks...)

			if len(tracks.Tracks) < limit {
				break
			}
			offset += limit
		}

		return AlbumTracksLoadedMsg{Tracks: allTracks}
	}
}

func (m Model) pollPlaybackState() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(m.ctx, 5*time.Second)
		defer cancel()

		state, err := m.client.PlayerState(ctx)
		if err != nil {
			// Don't treat as error - player might just be inactive
			return PlaybackStateMsg{State: nil}
		}

		return PlaybackStateMsg{State: state}
	}
}

func (m Model) schedulePlaybackPoll() tea.Cmd {
	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return PollPlaybackMsg{}
	})
}

// Playback control commands

func (m Model) togglePlayback() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		state, err := m.client.PlayerState(ctx)
		if err != nil {
			return ErrMsg{Err: err}
		}

		if state != nil && state.Playing {
			err = m.client.Pause(ctx)
		} else {
			err = m.client.Play(ctx)
		}

		if err != nil {
			return ErrMsg{Err: err}
		}

		// Refresh state after a short delay
		time.Sleep(200 * time.Millisecond)
		return PollPlaybackMsg{}
	}
}

func (m Model) nextTrack() tea.Cmd {
	return func() tea.Msg {
		if err := m.client.Next(context.Background()); err != nil {
			return ErrMsg{Err: err}
		}
		time.Sleep(200 * time.Millisecond)
		return PollPlaybackMsg{}
	}
}

func (m Model) prevTrack() tea.Cmd {
	return func() tea.Msg {
		if err := m.client.Previous(context.Background()); err != nil {
			return ErrMsg{Err: err}
		}
		time.Sleep(200 * time.Millisecond)
		return PollPlaybackMsg{}
	}
}

func (m Model) volumeUp() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		state, err := m.client.PlayerState(ctx)
		if err != nil || state == nil {
			return nil
		}

		newVol := int(state.Device.Volume) + 10
		if newVol > 100 {
			newVol = 100
		}

		m.client.Volume(ctx, newVol)
		return VolumeChangedMsg{Volume: newVol}
	}
}

func (m Model) volumeDown() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		state, err := m.client.PlayerState(ctx)
		if err != nil || state == nil {
			return nil
		}

		newVol := int(state.Device.Volume) - 10
		if newVol < 0 {
			newVol = 0
		}

		m.client.Volume(ctx, newVol)
		return VolumeChangedMsg{Volume: newVol}
	}
}

func (m Model) toggleShuffle() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		state, err := m.client.PlayerState(ctx)
		if err != nil || state == nil {
			return nil
		}

		m.client.Shuffle(ctx, !state.ShuffleState)
		return PollPlaybackMsg{}
	}
}

func (m Model) cycleRepeat() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		state, err := m.client.PlayerState(ctx)
		if err != nil || state == nil {
			return nil
		}

		// Cycle: off -> context -> track -> off
		var newState string
		switch state.RepeatState {
		case "off":
			newState = "context"
		case "context":
			newState = "track"
		default:
			newState = "off"
		}

		m.client.Repeat(ctx, newState)
		return PollPlaybackMsg{}
	}
}

func (m Model) playTrack(track spotify.PlaylistTrack) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		var opts *spotify.PlayOptions

		if m.selectedPlaylist != nil {
			opts = &spotify.PlayOptions{
				PlaybackContext: &m.selectedPlaylist.URI,
				PlaybackOffset: &spotify.PlaybackOffset{
					URI: track.Track.URI,
				},
			}
		} else {
			opts = &spotify.PlayOptions{
				URIs: []spotify.URI{track.Track.URI},
			}
		}

		if err := m.client.PlayOpt(ctx, opts); err != nil {
			return ErrMsg{Err: err}
		}

		time.Sleep(200 * time.Millisecond)
		return PollPlaybackMsg{}
	}
}

func (m Model) playAlbumTrack(track spotify.SimpleTrack) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		if m.selectedAlbum == nil {
			return nil
		}

		opts := &spotify.PlayOptions{
			PlaybackContext: &m.selectedAlbum.URI,
			PlaybackOffset: &spotify.PlaybackOffset{
				URI: track.URI,
			},
		}

		if err := m.client.PlayOpt(ctx, opts); err != nil {
			return ErrMsg{Err: err}
		}

		time.Sleep(200 * time.Millisecond)
		return PollPlaybackMsg{}
	}
}

func (m Model) fetchDevices() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(m.ctx, 10*time.Second)
		defer cancel()

		devices, err := m.client.PlayerDevices(ctx)
		if err != nil {
			return ErrMsg{Err: err}
		}

		return DevicesLoadedMsg{Devices: devices}
	}
}

func (m Model) transferPlayback(deviceID spotify.ID) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		err := m.client.TransferPlayback(ctx, deviceID, true)
		if err != nil {
			return ErrMsg{Err: err}
		}

		time.Sleep(300 * time.Millisecond)
		return PollPlaybackMsg{}
	}
}

type SearchResultsMsg struct {
	Tracks    []spotify.FullTrack
	Albums    []spotify.SimpleAlbum
	Playlists []spotify.SimplePlaylist
	Artists   []spotify.FullArtist
}

func (m Model) searchTracks(query string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(m.ctx, 15*time.Second)
		defer cancel()

		searchTypes := spotify.SearchTypeTrack | spotify.SearchTypeAlbum | spotify.SearchTypePlaylist | spotify.SearchTypeArtist
		result, err := m.client.Search(ctx, query, searchTypes, spotify.Limit(20))
		if err != nil {
			return ErrMsg{Err: err}
		}

		msg := SearchResultsMsg{}
		if result.Tracks != nil {
			msg.Tracks = result.Tracks.Tracks
		}
		if result.Albums != nil {
			msg.Albums = result.Albums.Albums
		}
		if result.Playlists != nil {
			msg.Playlists = result.Playlists.Playlists
		}
		if result.Artists != nil {
			msg.Artists = result.Artists.Artists
		}

		return msg
	}
}

func (m Model) playSearchTrack(track spotify.FullTrack) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		devices, err := m.client.PlayerDevices(ctx)
		if err != nil {
			return ErrMsg{Err: err}
		}

		var activeDeviceID *spotify.ID
		for _, d := range devices {
			if d.Active {
				activeDeviceID = &d.ID
				break
			}
		}

		if activeDeviceID == nil && len(devices) > 0 {
			activeDeviceID = &devices[0].ID
			if err := m.client.TransferPlayback(ctx, *activeDeviceID, false); err != nil {
				return ErrMsg{Err: err}
			}
			time.Sleep(300 * time.Millisecond)
		}

		if activeDeviceID == nil {
			return ErrMsg{Err: fmt.Errorf("no Spotify devices available - open Spotify on a device first")}
		}

		opts := &spotify.PlayOptions{
			URIs:     []spotify.URI{track.URI},
			DeviceID: activeDeviceID,
		}

		if err := m.client.PlayOpt(ctx, opts); err != nil {
			return ErrMsg{Err: err}
		}

		time.Sleep(200 * time.Millisecond)
		return PollPlaybackMsg{}
	}
}

func (m Model) playSearchItem(item views.SearchItem) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		devices, err := m.client.PlayerDevices(ctx)
		if err != nil {
			return ErrMsg{Err: err}
		}

		var activeDeviceID *spotify.ID
		for _, d := range devices {
			if d.Active {
				activeDeviceID = &d.ID
				break
			}
		}

		if activeDeviceID == nil && len(devices) > 0 {
			activeDeviceID = &devices[0].ID
			if err := m.client.TransferPlayback(ctx, *activeDeviceID, false); err != nil {
				return ErrMsg{Err: err}
			}
			time.Sleep(300 * time.Millisecond)
		}

		if activeDeviceID == nil {
			return ErrMsg{Err: fmt.Errorf("no Spotify devices available - open Spotify on a device first")}
		}

		var opts *spotify.PlayOptions
		switch item.Type {
		case views.SearchResultTrack:
			opts = &spotify.PlayOptions{
				URIs:     []spotify.URI{item.Track.URI},
				DeviceID: activeDeviceID,
			}
		case views.SearchResultAlbum:
			opts = &spotify.PlayOptions{
				PlaybackContext: &item.Album.URI,
				DeviceID:        activeDeviceID,
			}
		case views.SearchResultPlaylist:
			opts = &spotify.PlayOptions{
				PlaybackContext: &item.Playlist.URI,
				DeviceID:        activeDeviceID,
			}
		}

		if err := m.client.PlayOpt(ctx, opts); err != nil {
			return ErrMsg{Err: err}
		}

		time.Sleep(200 * time.Millisecond)
		return PollPlaybackMsg{}
	}
}

func (m Model) fetchRecentlyPlayed() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(m.ctx, 10*time.Second)
		defer cancel()

		items, err := m.client.PlayerRecentlyPlayed(ctx)
		if err != nil {
			return ErrMsg{Err: err}
		}

		return HistoryLoadedMsg{Items: items}
	}
}

func (m Model) playHistoryTrack(track spotify.SimpleTrack) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		devices, err := m.client.PlayerDevices(ctx)
		if err != nil {
			return ErrMsg{Err: err}
		}

		var activeDeviceID *spotify.ID
		for _, d := range devices {
			if d.Active {
				activeDeviceID = &d.ID
				break
			}
		}

		if activeDeviceID == nil && len(devices) > 0 {
			activeDeviceID = &devices[0].ID
			if err := m.client.TransferPlayback(ctx, *activeDeviceID, false); err != nil {
				return ErrMsg{Err: err}
			}
			time.Sleep(300 * time.Millisecond)
		}

		if activeDeviceID == nil {
			return ErrMsg{Err: fmt.Errorf("no Spotify devices available - open Spotify on a device first")}
		}

		opts := &spotify.PlayOptions{
			URIs:     []spotify.URI{track.URI},
			DeviceID: activeDeviceID,
		}

		if err := m.client.PlayOpt(ctx, opts); err != nil {
			return ErrMsg{Err: err}
		}

		time.Sleep(200 * time.Millisecond)
		return PollPlaybackMsg{}
	}
}

func (m Model) fetchArtist(artistID spotify.ID) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(m.ctx, 15*time.Second)
		defer cancel()

		artist, err := m.client.GetArtist(ctx, artistID)
		if err != nil {
			return ErrMsg{Err: err}
		}

		topTracks, err := m.client.GetArtistsTopTracks(ctx, artistID, "US")
		if err != nil {
			return ErrMsg{Err: err}
		}

		var allAlbums []spotify.SimpleAlbum
		limit := 50
		offset := 0

		for {
			albums, err := m.client.GetArtistAlbums(ctx, artistID,
				[]spotify.AlbumType{spotify.AlbumTypeAlbum, spotify.AlbumTypeSingle},
				spotify.Limit(limit), spotify.Offset(offset))
			if err != nil {
				return ErrMsg{Err: err}
			}

			allAlbums = append(allAlbums, albums.Albums...)

			if len(albums.Albums) < limit {
				break
			}
			offset += limit
		}

		return ArtistLoadedMsg{
			Artist:    artist,
			TopTracks: topTracks,
			Albums:    allAlbums,
		}
	}
}

func (m Model) playArtistTopTrack(track spotify.FullTrack) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		devices, err := m.client.PlayerDevices(ctx)
		if err != nil {
			return ErrMsg{Err: err}
		}

		var activeDeviceID *spotify.ID
		for _, d := range devices {
			if d.Active {
				activeDeviceID = &d.ID
				break
			}
		}

		if activeDeviceID == nil && len(devices) > 0 {
			activeDeviceID = &devices[0].ID
			if err := m.client.TransferPlayback(ctx, *activeDeviceID, false); err != nil {
				return ErrMsg{Err: err}
			}
			time.Sleep(300 * time.Millisecond)
		}

		if activeDeviceID == nil {
			return ErrMsg{Err: fmt.Errorf("no Spotify devices available - open Spotify on a device first")}
		}

		opts := &spotify.PlayOptions{
			URIs:     []spotify.URI{track.URI},
			DeviceID: activeDeviceID,
		}

		if err := m.client.PlayOpt(ctx, opts); err != nil {
			return ErrMsg{Err: err}
		}

		time.Sleep(200 * time.Millisecond)
		return PollPlaybackMsg{}
	}
}

type LyricsLoadedMsg struct {
	Lyrics       string
	SyncedLyrics string
	TrackName    string
	ArtistName   string
}

func (m Model) fetchLyrics(trackName, artistName string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		apiURL := fmt.Sprintf("https://lrclib.net/api/get?track_name=%s&artist_name=%s",
			url.QueryEscape(trackName),
			url.QueryEscape(artistName),
		)

		req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
		if err != nil {
			return LyricsLoadedMsg{Lyrics: "", SyncedLyrics: "", TrackName: trackName, ArtistName: artistName}
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return LyricsLoadedMsg{Lyrics: "", SyncedLyrics: "", TrackName: trackName, ArtistName: artistName}
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return LyricsLoadedMsg{Lyrics: "", SyncedLyrics: "", TrackName: trackName, ArtistName: artistName}
		}

		var result struct {
			PlainLyrics  string `json:"plainLyrics"`
			SyncedLyrics string `json:"syncedLyrics"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return LyricsLoadedMsg{Lyrics: "", SyncedLyrics: "", TrackName: trackName, ArtistName: artistName}
		}

		return LyricsLoadedMsg{
			Lyrics:       result.PlainLyrics,
			SyncedLyrics: result.SyncedLyrics,
			TrackName:    trackName,
			ArtistName:   artistName,
		}
	}
}

type LikeToggledMsg struct {
	TrackID   spotify.ID
	IsLiked   bool
	TrackName string
}

func (m Model) toggleLikeTrack(trackID spotify.ID, trackName string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		saved, err := m.client.UserHasTracks(ctx, trackID)
		if err != nil {
			return ErrMsg{Err: err}
		}

		if len(saved) == 0 {
			return ErrMsg{Err: fmt.Errorf("could not check track status")}
		}

		isCurrentlyLiked := saved[0]

		if isCurrentlyLiked {
			err = m.client.RemoveTracksFromLibrary(ctx, trackID)
		} else {
			err = m.client.AddTracksToLibrary(ctx, trackID)
		}

		if err != nil {
			return ErrMsg{Err: err}
		}

		return LikeToggledMsg{
			TrackID:   trackID,
			IsLiked:   !isCurrentlyLiked,
			TrackName: trackName,
		}
	}
}

func (m Model) addTrackToPlaylist(playlistID spotify.ID, playlistName string, trackID spotify.ID, trackName string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(m.ctx, 10*time.Second)
		defer cancel()

		_, err := m.client.AddTracksToPlaylist(ctx, playlistID, trackID)
		if err != nil {
			return ErrMsg{Err: err}
		}

		return TrackAddedToPlaylistMsg{
			PlaylistName: playlistName,
			TrackName:    trackName,
		}
	}
}
