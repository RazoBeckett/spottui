package tui

import (
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zmb3/spotify/v2"

	"github.com/razobeckett/spottui/internal/tui/views"
)

func TestNewModel(t *testing.T) {
	m := NewModel(nil)

	if m.view != ViewLoading {
		t.Errorf("NewModel() view = %v, want %v", m.view, ViewLoading)
	}

	if m.cache == nil {
		t.Error("NewModel() cache should not be nil")
	}

	if m.searchInput.Placeholder != "Search tracks..." {
		t.Errorf("NewModel() searchInput.Placeholder = %v, want 'Search tracks...'", m.searchInput.Placeholder)
	}
}

func TestViewConstants(t *testing.T) {
	tests := []struct {
		view View
		want int
	}{
		{ViewLoading, 0},
		{ViewPlaylists, 1},
		{ViewTracks, 2},
		{ViewAlbum, 3},
		{ViewArtist, 4},
		{ViewDevices, 5},
		{ViewSearch, 6},
		{ViewHistory, 7},
		{ViewHelp, 8},
		{ViewLyrics, 9},
		{ViewAddToPlaylist, 10},
	}

	for _, tt := range tests {
		if int(tt.view) != tt.want {
			t.Errorf("View constant %d = %d, want %d", tt.view, int(tt.view), tt.want)
		}
	}
}

func TestModelListDimensions(t *testing.T) {
	m := NewModel(nil)
	m.width = 100
	m.height = 50

	if m.listWidth() != 96 {
		t.Errorf("listWidth() = %d, want 96", m.listWidth())
	}

	expectedHeight := 50 - ListHeightSub
	if m.listHeight() != expectedHeight {
		t.Errorf("listHeight() = %d, want %d", m.listHeight(), expectedHeight)
	}
}

func TestModelListHeightMinimum(t *testing.T) {
	m := NewModel(nil)
	m.height = 10

	if m.listHeight() < MinListHeight {
		t.Errorf("listHeight() = %d, should be at least MinListHeight(%d)", m.listHeight(), MinListHeight)
	}
}

func TestUpdateWindowSizeMsg(t *testing.T) {
	m := NewModel(nil)
	msg := tea.WindowSizeMsg{Width: 120, Height: 40}

	newModel, cmd := m.Update(msg)
	updated := newModel.(Model)

	if updated.width != 120 {
		t.Errorf("Update(WindowSizeMsg) width = %d, want 120", updated.width)
	}
	if updated.height != 40 {
		t.Errorf("Update(WindowSizeMsg) height = %d, want 40", updated.height)
	}
	if cmd != nil {
		t.Error("Update(WindowSizeMsg) should return nil cmd")
	}
}

func TestUpdateSpinnerTickMsg(t *testing.T) {
	m := NewModel(nil)
	msg := spinner.TickMsg{}

	newModel, cmd := m.Update(msg)
	_ = newModel.(Model)

	if cmd == nil {
		t.Error("Update(spinner.TickMsg) should return a cmd for next tick")
	}
}

func TestUpdateUserDataMsg(t *testing.T) {
	m := NewModel(nil)
	m.width = 100
	m.height = 50

	user := &spotify.PrivateUser{
		User: spotify.User{
			DisplayName: "Test User",
		},
	}
	playlists := []spotify.SimplePlaylist{
		{Name: "Playlist 1"},
		{Name: "Playlist 2"},
	}

	msg := UserDataMsg{
		User:            user,
		Playlists:       playlists,
		LikedSongsTotal: 100,
	}

	newModel, cmd := m.Update(msg)
	updated := newModel.(Model)

	if updated.currentUser != user {
		t.Error("Update(UserDataMsg) should set currentUser")
	}
	if len(updated.playlistsData) != 2 {
		t.Errorf("Update(UserDataMsg) playlistsData len = %d, want 2", len(updated.playlistsData))
	}
	if updated.view != ViewPlaylists {
		t.Errorf("Update(UserDataMsg) view = %v, want ViewPlaylists", updated.view)
	}
	if cmd == nil {
		t.Error("Update(UserDataMsg) should return cmd for playback polling")
	}
}

func TestUpdateTracksLoadedMsg(t *testing.T) {
	m := NewModel(nil)
	m.width = 100
	m.height = 50
	m.fetching = true
	m.selectedPlaylist = &spotify.SimplePlaylist{Name: "Test Playlist"}

	tracks := []spotify.PlaylistTrack{
		{Track: spotify.FullTrack{SimpleTrack: spotify.SimpleTrack{Name: "Track 1"}}},
		{Track: spotify.FullTrack{SimpleTrack: spotify.SimpleTrack{Name: "Track 2"}}},
	}

	msg := TracksLoadedMsg{Tracks: tracks}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.fetching {
		t.Error("Update(TracksLoadedMsg) should set fetching to false")
	}
	if len(updated.tracksData) != 2 {
		t.Errorf("Update(TracksLoadedMsg) tracksData len = %d, want 2", len(updated.tracksData))
	}
}

func TestUpdatePlaybackStateMsg(t *testing.T) {
	m := NewModel(nil)

	state := &spotify.PlayerState{
		CurrentlyPlaying: spotify.CurrentlyPlaying{
			Progress: 5000,
			Playing:  true,
		},
	}

	msg := PlaybackStateMsg{State: state}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.playbackState != state {
		t.Error("Update(PlaybackStateMsg) should set playbackState")
	}
	if updated.localProgress != 5000 {
		t.Errorf("Update(PlaybackStateMsg) localProgress = %d, want 5000", updated.localProgress)
	}
	if !updated.isPlaying {
		t.Error("Update(PlaybackStateMsg) should set isPlaying to true")
	}
}

func TestUpdatePlaybackStateMsgNil(t *testing.T) {
	m := NewModel(nil)
	m.isPlaying = true

	msg := PlaybackStateMsg{State: nil}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.playbackState != nil {
		t.Error("Update(PlaybackStateMsg) with nil state should set playbackState to nil")
	}
}

func TestUpdateErrMsg(t *testing.T) {
	m := NewModel(nil)
	m.fetching = true

	msg := ErrMsg{Err: &testError{msg: "test error"}}

	newModel, cmd := m.Update(msg)
	updated := newModel.(Model)

	if updated.fetching {
		t.Error("Update(ErrMsg) should set fetching to false")
	}
	if !updated.showError {
		t.Error("Update(ErrMsg) should set showError to true")
	}
	if updated.errMsg != "test error" {
		t.Errorf("Update(ErrMsg) errMsg = %v, want 'test error'", updated.errMsg)
	}
	if cmd == nil {
		t.Error("Update(ErrMsg) should return cmd for error dismiss")
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func TestUpdateDismissErrorMsg(t *testing.T) {
	m := NewModel(nil)
	m.showError = true
	m.errMsg = "some error"

	msg := DismissErrorMsg{}

	newModel, cmd := m.Update(msg)
	updated := newModel.(Model)

	if updated.showError {
		t.Error("Update(DismissErrorMsg) should set showError to false")
	}
	if updated.errMsg != "" {
		t.Error("Update(DismissErrorMsg) should clear errMsg")
	}
	if cmd != nil {
		t.Error("Update(DismissErrorMsg) should return nil cmd")
	}
}

func TestUpdateVolumeChangedMsg(t *testing.T) {
	m := NewModel(nil)
	m.playbackState = &spotify.PlayerState{
		Device: spotify.PlayerDevice{Volume: 50},
	}

	msg := VolumeChangedMsg{Volume: 75}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if int(updated.playbackState.Device.Volume) != 75 {
		t.Errorf("Update(VolumeChangedMsg) volume = %d, want 75", int(updated.playbackState.Device.Volume))
	}
}

func TestUpdateShuffleToggledMsg(t *testing.T) {
	m := NewModel(nil)
	m.playbackState = &spotify.PlayerState{
		ShuffleState: false,
	}

	msg := ShuffleToggledMsg{NewState: true}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if !updated.playbackState.ShuffleState {
		t.Error("Update(ShuffleToggledMsg) should update shuffle state")
	}
}

func TestUpdateRepeatCycledMsg(t *testing.T) {
	m := NewModel(nil)
	m.playbackState = &spotify.PlayerState{
		RepeatState: "off",
	}

	msg := RepeatCycledMsg{NewState: "context"}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.playbackState.RepeatState != "context" {
		t.Errorf("Update(RepeatCycledMsg) repeat state = %v, want 'context'", updated.playbackState.RepeatState)
	}
}

func TestUpdateDevicesLoadedMsg(t *testing.T) {
	m := NewModel(nil)
	m.width = 100
	m.height = 50
	m.fetching = true

	devices := []spotify.PlayerDevice{
		{Name: "Device 1", Active: true},
		{Name: "Device 2", Active: false},
	}

	msg := DevicesLoadedMsg{Devices: devices}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.fetching {
		t.Error("Update(DevicesLoadedMsg) should set fetching to false")
	}
	if len(updated.devicesData) != 2 {
		t.Errorf("Update(DevicesLoadedMsg) devicesData len = %d, want 2", len(updated.devicesData))
	}
}

func TestUpdateProgressTickMsg(t *testing.T) {
	m := NewModel(nil)
	m.isPlaying = true
	m.view = ViewLyrics
	m.localProgress = 1000
	m.lastProgressAt = time.Now().Add(-100 * time.Millisecond)

	msg := ProgressTickMsg{}

	newModel, cmd := m.Update(msg)
	updated := newModel.(Model)

	if updated.localProgress <= 1000 {
		t.Errorf("Update(ProgressTickMsg) should increase localProgress, got %d", updated.localProgress)
	}
	if cmd == nil {
		t.Error("Update(ProgressTickMsg) should schedule next tick")
	}
}

func TestUpdateFetchingTickMsg(t *testing.T) {
	m := NewModel(nil)
	m.fetching = true
	m.fetchingDots = 1

	msg := FetchingTickMsg{}

	newModel, cmd := m.Update(msg)
	updated := newModel.(Model)

	if updated.fetchingDots != 2 {
		t.Errorf("Update(FetchingTickMsg) fetchingDots = %d, want 2", updated.fetchingDots)
	}
	if cmd == nil {
		t.Error("Update(FetchingTickMsg) should schedule next tick when fetching")
	}
}

func TestUpdateFetchingTickMsgNotFetching(t *testing.T) {
	m := NewModel(nil)
	m.fetching = false

	msg := FetchingTickMsg{}

	_, cmd := m.Update(msg)

	if cmd != nil {
		t.Error("Update(FetchingTickMsg) should return nil when not fetching")
	}
}

func TestUpdateLikeToggledMsg(t *testing.T) {
	m := NewModel(nil)

	msg := LikeToggledMsg{
		TrackID:   "track123",
		IsLiked:   true,
		TrackName: "Test Track",
	}

	newModel, cmd := m.Update(msg)
	updated := newModel.(Model)

	if !updated.showNotify {
		t.Error("Update(LikeToggledMsg) should show notification")
	}
	if updated.notifyMsg != "♥ Liked: Test Track" {
		t.Errorf("Update(LikeToggledMsg) notifyMsg = %v, want '♥ Liked: Test Track'", updated.notifyMsg)
	}
	if cmd == nil {
		t.Error("Update(LikeToggledMsg) should schedule notify dismiss")
	}
}

func TestUpdateLikeToggledMsgUnlike(t *testing.T) {
	m := NewModel(nil)

	msg := LikeToggledMsg{
		TrackID:   "track123",
		IsLiked:   false,
		TrackName: "Test Track",
	}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.notifyMsg != "♡ Unliked: Test Track" {
		t.Errorf("Update(LikeToggledMsg) notifyMsg = %v, want '♡ Unliked: Test Track'", updated.notifyMsg)
	}
}

func TestUpdateDismissNotifyMsg(t *testing.T) {
	m := NewModel(nil)
	m.showNotify = true
	m.notifyMsg = "some notification"

	msg := DismissNotifyMsg{}

	newModel, cmd := m.Update(msg)
	updated := newModel.(Model)

	if updated.showNotify {
		t.Error("Update(DismissNotifyMsg) should set showNotify to false")
	}
	if updated.notifyMsg != "" {
		t.Error("Update(DismissNotifyMsg) should clear notifyMsg")
	}
	if cmd != nil {
		t.Error("Update(DismissNotifyMsg) should return nil cmd")
	}
}

func TestUpdateTrackAddedToPlaylistMsg(t *testing.T) {
	m := NewModel(nil)

	msg := TrackAddedToPlaylistMsg{
		PlaylistName: "My Playlist",
		TrackName:    "Test Track",
	}

	newModel, cmd := m.Update(msg)
	updated := newModel.(Model)

	if !updated.showNotify {
		t.Error("Update(TrackAddedToPlaylistMsg) should show notification")
	}
	if updated.notifyMsg != "Added to My Playlist" {
		t.Errorf("Update(TrackAddedToPlaylistMsg) notifyMsg = %v, want 'Added to My Playlist'", updated.notifyMsg)
	}
	if cmd == nil {
		t.Error("Update(TrackAddedToPlaylistMsg) should schedule notify dismiss")
	}
}

func TestViewRendersTooSmall(t *testing.T) {
	m := NewModel(nil)
	m.width = 30
	m.height = 10

	output := m.View()

	if output == "" {
		t.Error("View() should render something even when too small")
	}
}

func TestViewRendersLoading(t *testing.T) {
	m := NewModel(nil)
	m.width = 100
	m.height = 50
	m.view = ViewLoading

	output := m.View()

	if output == "" {
		t.Error("View() should render loading view")
	}
}

func TestViewRendersPlaylists(t *testing.T) {
	m := NewModel(nil)
	m.width = 100
	m.height = 50
	m.view = ViewPlaylists

	output := m.View()

	if output == "" {
		t.Error("View() should render playlists view")
	}
}

func TestViewRendersTracks(t *testing.T) {
	m := NewModel(nil)
	m.width = 100
	m.height = 50
	m.view = ViewTracks

	output := m.View()

	if output == "" {
		t.Error("View() should render tracks view")
	}
}

func TestViewRendersHelp(t *testing.T) {
	m := NewModel(nil)
	m.width = 100
	m.height = 50
	m.view = ViewHelp

	output := m.View()

	if output == "" {
		t.Error("View() should render help view")
	}
}

func TestViewRendersSearch(t *testing.T) {
	m := NewModel(nil)
	m.width = 100
	m.height = 50
	m.view = ViewSearch

	output := m.View()

	if output == "" {
		t.Error("View() should render search view")
	}
}

func TestViewRendersDevices(t *testing.T) {
	m := NewModel(nil)
	m.width = 100
	m.height = 50
	m.view = ViewDevices

	output := m.View()

	if output == "" {
		t.Error("View() should render devices view")
	}
}

func TestViewRendersHistory(t *testing.T) {
	m := NewModel(nil)
	m.width = 100
	m.height = 50
	m.view = ViewHistory

	output := m.View()

	if output == "" {
		t.Error("View() should render history view")
	}
}

func TestViewRendersAlbum(t *testing.T) {
	m := NewModel(nil)
	m.width = 100
	m.height = 50
	m.view = ViewAlbum

	output := m.View()

	if output == "" {
		t.Error("View() should render album view")
	}
}

func TestViewRendersArtist(t *testing.T) {
	m := NewModel(nil)
	m.width = 100
	m.height = 50
	m.view = ViewArtist
	m.artistViewMode = "tracks"

	output := m.View()

	if output == "" {
		t.Error("View() should render artist view")
	}
}

func TestViewRendersLyrics(t *testing.T) {
	m := NewModel(nil)
	m.width = 100
	m.height = 50
	m.view = ViewLyrics
	m.lyricsData = "Test lyrics line 1\nTest lyrics line 2"

	output := m.View()

	if output == "" {
		t.Error("View() should render lyrics view")
	}
}

func TestViewRendersAddToPlaylist(t *testing.T) {
	m := NewModel(nil)
	m.width = 100
	m.height = 50
	m.view = ViewAddToPlaylist
	m.addToPlaylistList = views.CreateAddToPlaylistList(nil, m.styles, 96, 38)

	output := m.View()

	if output == "" {
		t.Error("View() should render add to playlist view")
	}
}

func TestViewRendersUnknown(t *testing.T) {
	m := NewModel(nil)
	m.width = 100
	m.height = 50
	m.view = View(99)

	output := m.View()

	if output != "Unknown view" {
		t.Errorf("View() for unknown view = %v, want 'Unknown view'", output)
	}
}

func TestUpdateSeekTickMsgPending(t *testing.T) {
	m := NewModel(nil)
	m.seekPending = true
	m.pendingSeek = 10000
	m.lastSeekRequest = time.Now().Add(-200 * time.Millisecond)

	msg := SeekTickMsg{}

	newModel, cmd := m.Update(msg)
	updated := newModel.(Model)

	if updated.seekPending {
		t.Error("Update(SeekTickMsg) should clear seekPending after delay")
	}
	if cmd == nil {
		t.Error("Update(SeekTickMsg) should return cmd to execute seek")
	}
}

func TestUpdateSeekTickMsgRecentRequest(t *testing.T) {
	m := NewModel(nil)
	m.seekPending = true
	m.pendingSeek = 10000
	m.lastSeekRequest = time.Now()

	msg := SeekTickMsg{}

	newModel, cmd := m.Update(msg)
	updated := newModel.(Model)

	if !updated.seekPending {
		t.Error("Update(SeekTickMsg) should keep seekPending for recent request")
	}
	if cmd == nil {
		t.Error("Update(SeekTickMsg) should schedule another tick")
	}
}

func TestUpdateSearchResultsMsg(t *testing.T) {
	m := NewModel(nil)
	m.width = 100
	m.height = 50
	m.fetching = true
	m.searching = true

	tracks := []spotify.FullTrack{
		{SimpleTrack: spotify.SimpleTrack{Name: "Track 1"}},
	}
	albums := []spotify.SimpleAlbum{
		{Name: "Album 1"},
	}

	msg := SearchResultsMsg{
		Tracks: tracks,
		Albums: albums,
	}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.fetching {
		t.Error("Update(SearchResultsMsg) should set fetching to false")
	}
	if updated.searching {
		t.Error("Update(SearchResultsMsg) should set searching to false")
	}
	if len(updated.searchTracksData) != 1 {
		t.Errorf("Update(SearchResultsMsg) searchTracksData len = %d, want 1", len(updated.searchTracksData))
	}
	if len(updated.searchAlbumsData) != 1 {
		t.Errorf("Update(SearchResultsMsg) searchAlbumsData len = %d, want 1", len(updated.searchAlbumsData))
	}
}

func TestUpdateHistoryLoadedMsg(t *testing.T) {
	m := NewModel(nil)
	m.width = 100
	m.height = 50
	m.fetching = true

	items := []spotify.RecentlyPlayedItem{
		{Track: spotify.SimpleTrack{Name: "Recent Track"}},
	}

	msg := HistoryLoadedMsg{Items: items}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.fetching {
		t.Error("Update(HistoryLoadedMsg) should set fetching to false")
	}
	if len(updated.historyData) != 1 {
		t.Errorf("Update(HistoryLoadedMsg) historyData len = %d, want 1", len(updated.historyData))
	}
}

func TestUpdateAlbumTracksLoadedMsg(t *testing.T) {
	m := NewModel(nil)
	m.width = 100
	m.height = 50
	m.fetching = true
	m.selectedAlbum = &spotify.SimpleAlbum{Name: "Test Album"}

	tracks := []spotify.SimpleTrack{
		{Name: "Album Track 1"},
		{Name: "Album Track 2"},
	}

	msg := AlbumTracksLoadedMsg{Tracks: tracks}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.fetching {
		t.Error("Update(AlbumTracksLoadedMsg) should set fetching to false")
	}
	if len(updated.albumTracksData) != 2 {
		t.Errorf("Update(AlbumTracksLoadedMsg) albumTracksData len = %d, want 2", len(updated.albumTracksData))
	}
}

func TestUpdateArtistLoadedMsg(t *testing.T) {
	m := NewModel(nil)
	m.width = 100
	m.height = 50
	m.fetching = true

	artist := &spotify.FullArtist{
		SimpleArtist: spotify.SimpleArtist{Name: "Test Artist"},
	}
	topTracks := []spotify.FullTrack{
		{SimpleTrack: spotify.SimpleTrack{Name: "Top Track 1"}},
	}
	albums := []spotify.SimpleAlbum{
		{Name: "Artist Album 1"},
	}

	msg := ArtistLoadedMsg{
		Artist:    artist,
		TopTracks: topTracks,
		Albums:    albums,
	}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.fetching {
		t.Error("Update(ArtistLoadedMsg) should set fetching to false")
	}
	if updated.selectedArtist != artist {
		t.Error("Update(ArtistLoadedMsg) should set selectedArtist")
	}
	if updated.artistViewMode != "tracks" {
		t.Errorf("Update(ArtistLoadedMsg) artistViewMode = %v, want 'tracks'", updated.artistViewMode)
	}
	if len(updated.artistTopTracksData) != 1 {
		t.Errorf("Update(ArtistLoadedMsg) artistTopTracksData len = %d, want 1", len(updated.artistTopTracksData))
	}
	if len(updated.artistAlbumsData) != 1 {
		t.Errorf("Update(ArtistLoadedMsg) artistAlbumsData len = %d, want 1", len(updated.artistAlbumsData))
	}
}

func TestUpdateLyricsLoadedMsg(t *testing.T) {
	m := NewModel(nil)
	m.fetching = true
	m.fetchingLyrics = true

	msg := LyricsLoadedMsg{
		Lyrics:       "Plain lyrics text",
		SyncedLyrics: "",
		TrackName:    "Test Track",
		ArtistName:   "Test Artist",
	}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.fetching {
		t.Error("Update(LyricsLoadedMsg) should set fetching to false")
	}
	if updated.fetchingLyrics {
		t.Error("Update(LyricsLoadedMsg) should set fetchingLyrics to false")
	}
	if updated.lyricsData != "Plain lyrics text" {
		t.Errorf("Update(LyricsLoadedMsg) lyricsData = %v, want 'Plain lyrics text'", updated.lyricsData)
	}
	if updated.lyricsTrackName != "Test Track" {
		t.Errorf("Update(LyricsLoadedMsg) lyricsTrackName = %v, want 'Test Track'", updated.lyricsTrackName)
	}
	if updated.view != ViewLyrics {
		t.Errorf("Update(LyricsLoadedMsg) view = %v, want ViewLyrics", updated.view)
	}
}

func TestUpdateLyricsLoadedMsgWithSynced(t *testing.T) {
	m := NewModel(nil)
	m.fetching = true

	msg := LyricsLoadedMsg{
		Lyrics:       "Plain lyrics",
		SyncedLyrics: "[00:01.00]Line one\n[00:05.00]Line two",
		TrackName:    "Test Track",
		ArtistName:   "Test Artist",
	}

	newModel, cmd := m.Update(msg)
	updated := newModel.(Model)

	if !updated.lyricsIsSynced {
		t.Error("Update(LyricsLoadedMsg) should set lyricsIsSynced to true for synced lyrics")
	}
	if len(updated.lyricsSynced) == 0 {
		t.Error("Update(LyricsLoadedMsg) should parse synced lyrics")
	}
	if cmd == nil {
		t.Error("Update(LyricsLoadedMsg) should schedule progress tick for synced lyrics")
	}
}
