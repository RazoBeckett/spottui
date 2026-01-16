package tui

import (
	"testing"
	"time"

	"github.com/zmb3/spotify/v2"
)

func TestScheduleErrorDismiss(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	cmd := m.scheduleErrorDismiss()

	if cmd == nil {
		t.Error("scheduleErrorDismiss() should return a cmd")
	}
}

func TestScheduleNotifyDismiss(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	cmd := m.scheduleNotifyDismiss()

	if cmd == nil {
		t.Error("scheduleNotifyDismiss() should return a cmd")
	}
}

func TestSchedulePlaybackPoll(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	cmd := m.schedulePlaybackPoll()

	if cmd == nil {
		t.Error("schedulePlaybackPoll() should return a cmd")
	}
}

func TestScheduleProgressTick(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	cmd := m.scheduleProgressTick()

	if cmd == nil {
		t.Error("scheduleProgressTick() should return a cmd")
	}
}

func TestScheduleFetchingTick(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	cmd := m.scheduleFetchingTick()

	if cmd == nil {
		t.Error("scheduleFetchingTick() should return a cmd")
	}
}

func TestScheduleSeekTick(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	cmd := m.scheduleSeekTick()

	if cmd == nil {
		t.Error("scheduleSeekTick() should return a cmd")
	}
}

func TestStartFetching(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.fetching = false
	m.fetchingDots = 5

	cmd := m.startFetching()

	if !m.fetching {
		t.Error("startFetching() should set fetching to true")
	}
	if m.fetchingDots != 0 {
		t.Errorf("startFetching() should reset fetchingDots, got %d", m.fetchingDots)
	}
	if cmd == nil {
		t.Error("startFetching() should return a cmd")
	}
}

func TestStopFetching(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.fetching = true
	m.fetchingDots = 3

	m.stopFetching()

	if m.fetching {
		t.Error("stopFetching() should set fetching to false")
	}
	if m.fetchingDots != 0 {
		t.Errorf("stopFetching() should reset fetchingDots, got %d", m.fetchingDots)
	}
}

func TestUserDataMsgFields(t *testing.T) {
	user := &spotify.PrivateUser{}
	playlists := []spotify.SimplePlaylist{{Name: "Test"}}
	msg := UserDataMsg{
		User:            user,
		Playlists:       playlists,
		LikedSongsTotal: 50,
	}

	if msg.User != user {
		t.Error("UserDataMsg.User should be set")
	}
	if len(msg.Playlists) != 1 {
		t.Error("UserDataMsg.Playlists should have one item")
	}
	if msg.LikedSongsTotal != 50 {
		t.Errorf("UserDataMsg.LikedSongsTotal = %d, want 50", msg.LikedSongsTotal)
	}
}

func TestTracksLoadedMsgFields(t *testing.T) {
	tracks := []spotify.PlaylistTrack{
		{Track: spotify.FullTrack{SimpleTrack: spotify.SimpleTrack{Name: "Track 1"}}},
	}
	msg := TracksLoadedMsg{Tracks: tracks}

	if len(msg.Tracks) != 1 {
		t.Errorf("TracksLoadedMsg.Tracks len = %d, want 1", len(msg.Tracks))
	}
}

func TestPlaybackStateMsgFields(t *testing.T) {
	state := &spotify.PlayerState{}
	msg := PlaybackStateMsg{State: state}

	if msg.State != state {
		t.Error("PlaybackStateMsg.State should be set")
	}
}

func TestShuffleToggledMsgFields(t *testing.T) {
	msg := ShuffleToggledMsg{NewState: true}

	if !msg.NewState {
		t.Error("ShuffleToggledMsg.NewState should be true")
	}
}

func TestRepeatCycledMsgFields(t *testing.T) {
	msg := RepeatCycledMsg{NewState: "context"}

	if msg.NewState != "context" {
		t.Errorf("RepeatCycledMsg.NewState = %s, want 'context'", msg.NewState)
	}
}

func TestDevicesLoadedMsgFields(t *testing.T) {
	devices := []spotify.PlayerDevice{{Name: "Device 1"}}
	msg := DevicesLoadedMsg{Devices: devices}

	if len(msg.Devices) != 1 {
		t.Errorf("DevicesLoadedMsg.Devices len = %d, want 1", len(msg.Devices))
	}
}

func TestAlbumTracksLoadedMsgFields(t *testing.T) {
	tracks := []spotify.SimpleTrack{{Name: "Track 1"}}
	msg := AlbumTracksLoadedMsg{Tracks: tracks}

	if len(msg.Tracks) != 1 {
		t.Errorf("AlbumTracksLoadedMsg.Tracks len = %d, want 1", len(msg.Tracks))
	}
}

func TestHistoryLoadedMsgFields(t *testing.T) {
	items := []spotify.RecentlyPlayedItem{
		{Track: spotify.SimpleTrack{Name: "Recent Track"}, PlayedAt: time.Now()},
	}
	msg := HistoryLoadedMsg{Items: items}

	if len(msg.Items) != 1 {
		t.Errorf("HistoryLoadedMsg.Items len = %d, want 1", len(msg.Items))
	}
}

func TestArtistLoadedMsgFields(t *testing.T) {
	artist := &spotify.FullArtist{}
	topTracks := []spotify.FullTrack{{SimpleTrack: spotify.SimpleTrack{Name: "Top Track"}}}
	albums := []spotify.SimpleAlbum{{Name: "Album 1"}}

	msg := ArtistLoadedMsg{
		Artist:    artist,
		TopTracks: topTracks,
		Albums:    albums,
	}

	if msg.Artist != artist {
		t.Error("ArtistLoadedMsg.Artist should be set")
	}
	if len(msg.TopTracks) != 1 {
		t.Errorf("ArtistLoadedMsg.TopTracks len = %d, want 1", len(msg.TopTracks))
	}
	if len(msg.Albums) != 1 {
		t.Errorf("ArtistLoadedMsg.Albums len = %d, want 1", len(msg.Albums))
	}
}

func TestErrMsgFields(t *testing.T) {
	err := &testCmdError{msg: "test error"}
	msg := ErrMsg{Err: err}

	if msg.Err == nil {
		t.Error("ErrMsg.Err should be set")
	}
	if msg.Err.Error() != "test error" {
		t.Errorf("ErrMsg.Err.Error() = %s, want 'test error'", msg.Err.Error())
	}
}

type testCmdError struct {
	msg string
}

func (e *testCmdError) Error() string {
	return e.msg
}

func TestVolumeChangedMsgFields(t *testing.T) {
	msg := VolumeChangedMsg{Volume: 75}

	if msg.Volume != 75 {
		t.Errorf("VolumeChangedMsg.Volume = %d, want 75", msg.Volume)
	}
}

func TestTrackAddedToPlaylistMsgFields(t *testing.T) {
	msg := TrackAddedToPlaylistMsg{
		PlaylistName: "My Playlist",
		TrackName:    "Test Track",
	}

	if msg.PlaylistName != "My Playlist" {
		t.Errorf("TrackAddedToPlaylistMsg.PlaylistName = %s, want 'My Playlist'", msg.PlaylistName)
	}
	if msg.TrackName != "Test Track" {
		t.Errorf("TrackAddedToPlaylistMsg.TrackName = %s, want 'Test Track'", msg.TrackName)
	}
}

func TestSearchResultsMsgFields(t *testing.T) {
	tracks := []spotify.FullTrack{{SimpleTrack: spotify.SimpleTrack{Name: "Track 1"}}}
	albums := []spotify.SimpleAlbum{{Name: "Album 1"}}
	playlists := []spotify.SimplePlaylist{{Name: "Playlist 1"}}
	artists := []spotify.FullArtist{{SimpleArtist: spotify.SimpleArtist{Name: "Artist 1"}}}

	msg := SearchResultsMsg{
		Tracks:    tracks,
		Albums:    albums,
		Playlists: playlists,
		Artists:   artists,
	}

	if len(msg.Tracks) != 1 {
		t.Errorf("SearchResultsMsg.Tracks len = %d, want 1", len(msg.Tracks))
	}
	if len(msg.Albums) != 1 {
		t.Errorf("SearchResultsMsg.Albums len = %d, want 1", len(msg.Albums))
	}
	if len(msg.Playlists) != 1 {
		t.Errorf("SearchResultsMsg.Playlists len = %d, want 1", len(msg.Playlists))
	}
	if len(msg.Artists) != 1 {
		t.Errorf("SearchResultsMsg.Artists len = %d, want 1", len(msg.Artists))
	}
}

func TestLyricsLoadedMsgFields(t *testing.T) {
	msg := LyricsLoadedMsg{
		Lyrics:       "Plain lyrics",
		SyncedLyrics: "[00:01.00]Line one",
		TrackName:    "Test Track",
		ArtistName:   "Test Artist",
	}

	if msg.Lyrics != "Plain lyrics" {
		t.Errorf("LyricsLoadedMsg.Lyrics = %s, want 'Plain lyrics'", msg.Lyrics)
	}
	if msg.SyncedLyrics != "[00:01.00]Line one" {
		t.Errorf("LyricsLoadedMsg.SyncedLyrics = %s, want '[00:01.00]Line one'", msg.SyncedLyrics)
	}
	if msg.TrackName != "Test Track" {
		t.Errorf("LyricsLoadedMsg.TrackName = %s, want 'Test Track'", msg.TrackName)
	}
	if msg.ArtistName != "Test Artist" {
		t.Errorf("LyricsLoadedMsg.ArtistName = %s, want 'Test Artist'", msg.ArtistName)
	}
}

func TestLikeToggledMsgFields(t *testing.T) {
	msg := LikeToggledMsg{
		TrackID:   "track123",
		IsLiked:   true,
		TrackName: "Test Track",
	}

	if msg.TrackID != "track123" {
		t.Errorf("LikeToggledMsg.TrackID = %s, want 'track123'", msg.TrackID)
	}
	if !msg.IsLiked {
		t.Error("LikeToggledMsg.IsLiked should be true")
	}
	if msg.TrackName != "Test Track" {
		t.Errorf("LikeToggledMsg.TrackName = %s, want 'Test Track'", msg.TrackName)
	}
}

func TestTickMessageTypes(t *testing.T) {
	_ = ProgressTickMsg{}
	_ = SeekTickMsg{}
	_ = FetchingTickMsg{}
	_ = PollPlaybackMsg{}
	_ = DismissErrorMsg{}
	_ = DismissNotifyMsg{}
}
