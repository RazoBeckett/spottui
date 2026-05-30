package tui

import (
	"testing"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"github.com/zmb3/spotify/v2"

	"github.com/razobeckett/spottui/internal/tui/views"
)

func TestNewModel(t *testing.T) {
	m := NewModel(nil, getTestConfig())

	if m.Nav.Current != ViewLoading {
		t.Errorf("NewModel() view = %v, want %v", m.Nav.Current, ViewLoading)
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
	m := NewModel(nil, getTestConfig())
	m.UI.Width = 100
	m.UI.Height = 50

	if m.listWidth() != 96 {
		t.Errorf("listWidth() = %d, want 96", m.listWidth())
	}

	expectedHeight := 50 - ListHeightSub
	if m.listHeight() != expectedHeight {
		t.Errorf("listHeight() = %d, want %d", m.listHeight(), expectedHeight)
	}
}

func TestModelListHeightMinimum(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.UI.Height = 10

	if m.listHeight() < MinListHeight {
		t.Errorf("listHeight() = %d, should be at least MinListHeight(%d)", m.listHeight(), MinListHeight)
	}
}

func TestUpdateWindowSizeMsg(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	msg := tea.WindowSizeMsg{Width: 120, Height: 40}

	newModel, cmd := m.Update(msg)
	updated := newModel.(Model)

	if updated.UI.Width != 120 {
		t.Errorf("Update(WindowSizeMsg) width = %d, want 120", updated.UI.Width)
	}
	if updated.UI.Height != 40 {
		t.Errorf("Update(WindowSizeMsg) height = %d, want 40", updated.UI.Height)
	}
	if cmd != nil {
		t.Error("Update(WindowSizeMsg) should return nil cmd")
	}
}

func TestUpdateSpinnerTickMsg(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	msg := spinner.TickMsg{}

	newModel, cmd := m.Update(msg)
	_ = newModel.(Model)

	if cmd == nil {
		t.Error("Update(spinner.TickMsg) should return a cmd for next tick")
	}
}

func TestUpdateUserDataMsg(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.UI.Width = 100
	m.UI.Height = 50

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
	if updated.Nav.Current != ViewPlaylists {
		t.Errorf("Update(UserDataMsg) view = %v, want ViewPlaylists", updated.Nav.Current)
	}
	if cmd == nil {
		t.Error("Update(UserDataMsg) should return cmd for playback polling")
	}
}

func TestUpdateTracksLoadedMsg(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.UI.Width = 100
	m.UI.Height = 50
	m.UI.Fetching = true
	m.selectedPlaylist = &spotify.SimplePlaylist{Name: "Test Playlist"}

	tracks := []spotify.PlaylistTrack{
		{Track: spotify.FullTrack{SimpleTrack: spotify.SimpleTrack{Name: "Track 1"}}},
		{Track: spotify.FullTrack{SimpleTrack: spotify.SimpleTrack{Name: "Track 2"}}},
	}

	msg := TracksLoadedMsg{Tracks: tracks}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.UI.Fetching {
		t.Error("Update(TracksLoadedMsg) should set fetching to false")
	}
	if len(updated.tracksData) != 2 {
		t.Errorf("Update(TracksLoadedMsg) tracksData len = %d, want 2", len(updated.tracksData))
	}
}

func TestUpdatePlaybackStateMsg(t *testing.T) {
	m := NewModel(nil, getTestConfig())

	state := &spotify.PlayerState{
		CurrentlyPlaying: spotify.CurrentlyPlaying{
			Progress: 5000,
			Playing:  true,
		},
	}

	msg := PlaybackStateMsg{State: state}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.Playback.State != state {
		t.Error("Update(PlaybackStateMsg) should set playbackState")
	}
	if updated.Playback.LocalProgress != 5000 {
		t.Errorf("Update(PlaybackStateMsg) localProgress = %d, want 5000", updated.Playback.LocalProgress)
	}
	if !updated.Playback.IsPlaying {
		t.Error("Update(PlaybackStateMsg) should set isPlaying to true")
	}
}

func TestUpdatePlaybackStateMsgNil(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.Playback.IsPlaying = true

	msg := PlaybackStateMsg{State: nil}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.Playback.State != nil {
		t.Error("Update(PlaybackStateMsg) with nil state should set playbackState to nil")
	}
}

func TestUpdateErrMsg(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.UI.Fetching = true

	msg := ErrMsg{Err: &testError{msg: "test error"}}

	newModel, cmd := m.Update(msg)
	updated := newModel.(Model)

	if updated.UI.Fetching {
		t.Error("Update(ErrMsg) should set fetching to false")
	}
	if !updated.UI.ShowError {
		t.Error("Update(ErrMsg) should set showError to true")
	}
	if updated.UI.ErrMsg != "test error" {
		t.Errorf("Update(ErrMsg) errMsg = %v, want 'test error'", updated.UI.ErrMsg)
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
	m := NewModel(nil, getTestConfig())
	m.UI.ShowError = true
	m.UI.ErrMsg = "some error"

	msg := DismissErrorMsg{}

	newModel, cmd := m.Update(msg)
	updated := newModel.(Model)

	if updated.UI.ShowError {
		t.Error("Update(DismissErrorMsg) should set showError to false")
	}
	if updated.UI.ErrMsg != "" {
		t.Error("Update(DismissErrorMsg) should clear errMsg")
	}
	if cmd != nil {
		t.Error("Update(DismissErrorMsg) should return nil cmd")
	}
}

func TestUpdateVolumeChangedMsg(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.Playback.State = &spotify.PlayerState{
		Device: spotify.PlayerDevice{ID: "test-device", Volume: 50},
	}

	msg := VolumeChangedMsg{Volume: 75}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if int(updated.Playback.State.Device.Volume) != 75 {
		t.Errorf("Update(VolumeChangedMsg) volume = %d, want 75", int(updated.Playback.State.Device.Volume))
	}
}

func TestUpdateShuffleToggledMsg(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.Playback.State = &spotify.PlayerState{
		ShuffleState: false,
	}

	msg := ShuffleToggledMsg{NewState: true}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if !updated.Playback.State.ShuffleState {
		t.Error("Update(ShuffleToggledMsg) should update shuffle state")
	}
}

func TestUpdateRepeatCycledMsg(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.Playback.State = &spotify.PlayerState{
		RepeatState: "off",
	}

	msg := RepeatCycledMsg{NewState: "context"}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.Playback.State.RepeatState != "context" {
		t.Errorf("Update(RepeatCycledMsg) repeat state = %v, want 'context'", updated.Playback.State.RepeatState)
	}
}

func TestUpdateDevicesLoadedMsg(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.UI.Width = 100
	m.UI.Height = 50
	m.UI.Fetching = true

	devices := []spotify.PlayerDevice{
		{Name: "Device 1", Active: true},
		{Name: "Device 2", Active: false},
	}

	msg := DevicesLoadedMsg{Devices: devices}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.UI.Fetching {
		t.Error("Update(DevicesLoadedMsg) should set fetching to false")
	}
	if len(updated.devicesData) != 2 {
		t.Errorf("Update(DevicesLoadedMsg) devicesData len = %d, want 2", len(updated.devicesData))
	}
}

func TestUpdateProgressTickMsg(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.Playback.IsPlaying = true
	m.Nav.Current = ViewLyrics
	m.Playback.LocalProgress = 1000
	m.Playback.LastProgressAt = time.Now().Add(-100 * time.Millisecond)

	msg := ProgressTickMsg{}

	newModel, cmd := m.Update(msg)
	updated := newModel.(Model)

	if updated.Playback.LocalProgress <= 1000 {
		t.Errorf("Update(ProgressTickMsg) should increase localProgress, got %d", updated.Playback.LocalProgress)
	}
	if cmd == nil {
		t.Error("Update(ProgressTickMsg) should schedule next tick")
	}
}

func TestUpdateFetchingTickMsg(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.UI.Fetching = true
	m.UI.FetchingDots = 1

	msg := FetchingTickMsg{}

	newModel, cmd := m.Update(msg)
	updated := newModel.(Model)

	if updated.UI.FetchingDots != 2 {
		t.Errorf("Update(FetchingTickMsg) fetchingDots = %d, want 2", updated.UI.FetchingDots)
	}
	if cmd == nil {
		t.Error("Update(FetchingTickMsg) should schedule next tick when fetching")
	}
}

func TestUpdateFetchingTickMsgNotFetching(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.UI.Fetching = false

	msg := FetchingTickMsg{}

	_, cmd := m.Update(msg)

	if cmd != nil {
		t.Error("Update(FetchingTickMsg) should return nil when not fetching")
	}
}

func TestUpdateLikeToggledMsg(t *testing.T) {
	m := NewModel(nil, getTestConfig())

	msg := LikeToggledMsg{
		TrackID:   "track123",
		IsLiked:   true,
		TrackName: "Test Track",
	}

	newModel, cmd := m.Update(msg)
	updated := newModel.(Model)

	if !updated.UI.ShowNotify {
		t.Error("Update(LikeToggledMsg) should show notification")
	}
	if updated.UI.NotifyMsg != "♥ Liked: Test Track" {
		t.Errorf("Update(LikeToggledMsg) notifyMsg = %v, want '♥ Liked: Test Track'", updated.UI.NotifyMsg)
	}
	if cmd == nil {
		t.Error("Update(LikeToggledMsg) should schedule notify dismiss")
	}
}

func TestUpdateLikeToggledMsgUnlike(t *testing.T) {
	m := NewModel(nil, getTestConfig())

	msg := LikeToggledMsg{
		TrackID:   "track123",
		IsLiked:   false,
		TrackName: "Test Track",
	}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.UI.NotifyMsg != "♡ Unliked: Test Track" {
		t.Errorf("Update(LikeToggledMsg) notifyMsg = %v, want '♡ Unliked: Test Track'", updated.UI.NotifyMsg)
	}
}

func TestUpdateDismissNotifyMsg(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.UI.ShowNotify = true
	m.UI.NotifyMsg = "some notification"

	msg := DismissNotifyMsg{}

	newModel, cmd := m.Update(msg)
	updated := newModel.(Model)

	if updated.UI.ShowNotify {
		t.Error("Update(DismissNotifyMsg) should set showNotify to false")
	}
	if updated.UI.NotifyMsg != "" {
		t.Error("Update(DismissNotifyMsg) should clear notifyMsg")
	}
	if cmd != nil {
		t.Error("Update(DismissNotifyMsg) should return nil cmd")
	}
}

func TestUpdateTrackAddedToPlaylistMsg(t *testing.T) {
	m := NewModel(nil, getTestConfig())

	msg := TrackAddedToPlaylistMsg{
		PlaylistName: "My Playlist",
		TrackName:    "Test Track",
	}

	newModel, cmd := m.Update(msg)
	updated := newModel.(Model)

	if !updated.UI.ShowNotify {
		t.Error("Update(TrackAddedToPlaylistMsg) should show notification")
	}
	if updated.UI.NotifyMsg != "Added to My Playlist" {
		t.Errorf("Update(TrackAddedToPlaylistMsg) notifyMsg = %v, want 'Added to My Playlist'", updated.UI.NotifyMsg)
	}
	if cmd == nil {
		t.Error("Update(TrackAddedToPlaylistMsg) should schedule notify dismiss")
	}
}

func TestViewRendersTooSmall(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.UI.Width = 30
	m.UI.Height = 10

	output := m.renderContent()

	if output == "" {
		t.Error("View() should render something even when too small")
	}
}

func TestViewRendersLoading(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.UI.Width = 100
	m.UI.Height = 50
	m.Nav.Current = ViewLoading

	output := m.renderContent()

	if output == "" {
		t.Error("View() should render loading view")
	}
}

func TestViewRendersPlaylists(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.UI.Width = 100
	m.UI.Height = 50
	m.Nav.Current = ViewPlaylists

	output := m.renderContent()

	if output == "" {
		t.Error("View() should render playlists view")
	}
}

func TestViewRendersTracks(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.UI.Width = 100
	m.UI.Height = 50
	m.Nav.Current = ViewTracks

	output := m.renderContent()

	if output == "" {
		t.Error("View() should render tracks view")
	}
}

func TestViewRendersHelp(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.UI.Width = 100
	m.UI.Height = 50
	m.Nav.Current = ViewHelp

	output := m.renderContent()

	if output == "" {
		t.Error("View() should render help view")
	}
}

func TestViewRendersSearch(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.UI.Width = 100
	m.UI.Height = 50
	m.Nav.Current = ViewSearch

	output := m.renderContent()

	if output == "" {
		t.Error("View() should render search view")
	}
}

func TestViewRendersDevices(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.UI.Width = 100
	m.UI.Height = 50
	m.Nav.Current = ViewDevices

	output := m.renderContent()

	if output == "" {
		t.Error("View() should render devices view")
	}
}

func TestViewRendersHistory(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.UI.Width = 100
	m.UI.Height = 50
	m.Nav.Current = ViewHistory

	output := m.renderContent()

	if output == "" {
		t.Error("View() should render history view")
	}
}

func TestViewRendersAlbum(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.UI.Width = 100
	m.UI.Height = 50
	m.Nav.Current = ViewAlbum

	output := m.renderContent()

	if output == "" {
		t.Error("View() should render album view")
	}
}

func TestViewRendersArtist(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.UI.Width = 100
	m.UI.Height = 50
	m.Nav.Current = ViewArtist
	m.artistViewMode = "tracks"

	output := m.renderContent()

	if output == "" {
		t.Error("View() should render artist view")
	}
}

func TestViewRendersLyrics(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.UI.Width = 100
	m.UI.Height = 50
	m.Nav.Current = ViewLyrics
	m.lyricsData = "Test lyrics line 1\nTest lyrics line 2"

	output := m.renderContent()

	if output == "" {
		t.Error("View() should render lyrics view")
	}
}

func TestViewRendersAddToPlaylist(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.UI.Width = 100
	m.UI.Height = 50
	m.Nav.Current = ViewAddToPlaylist
	m.addToPlaylistList = views.CreateAddToPlaylistList(nil, m.styles, 96, 38)

	output := m.renderContent()

	if output == "" {
		t.Error("View() should render add to playlist view")
	}
}

func TestViewRendersUnknown(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.UI.Width = 100
	m.UI.Height = 50
	m.Nav.Current = View(99)

	output := m.renderContent()

	if output != "Unknown view" {
		t.Errorf("View() for unknown view = %v, want 'Unknown view'", output)
	}
}

func TestUpdateSeekTickMsgPending(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.Playback.SeekPending = true
	m.Playback.PendingSeek = 10000
	m.Playback.LastSeekRequest = time.Now().Add(-200 * time.Millisecond)

	msg := SeekTickMsg{}

	newModel, cmd := m.Update(msg)
	updated := newModel.(Model)

	if updated.Playback.SeekPending {
		t.Error("Update(SeekTickMsg) should clear seekPending after delay")
	}
	if cmd == nil {
		t.Error("Update(SeekTickMsg) should return cmd to execute seek")
	}
}

func TestUpdateSeekTickMsgRecentRequest(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.Playback.SeekPending = true
	m.Playback.PendingSeek = 10000
	m.Playback.LastSeekRequest = time.Now()

	msg := SeekTickMsg{}

	newModel, cmd := m.Update(msg)
	updated := newModel.(Model)

	if !updated.Playback.SeekPending {
		t.Error("Update(SeekTickMsg) should keep seekPending for recent request")
	}
	if cmd == nil {
		t.Error("Update(SeekTickMsg) should schedule another tick")
	}
}

func TestUpdateSearchResultsMsg(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.UI.Width = 100
	m.UI.Height = 50
	m.UI.Fetching = true
	m.UI.Searching = true

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

	if updated.UI.Fetching {
		t.Error("Update(SearchResultsMsg) should set fetching to false")
	}
	if updated.UI.Searching {
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
	m := NewModel(nil, getTestConfig())
	m.UI.Width = 100
	m.UI.Height = 50
	m.UI.Fetching = true

	items := []spotify.RecentlyPlayedItem{
		{Track: spotify.SimpleTrack{Name: "Recent Track"}},
	}

	msg := HistoryLoadedMsg{Items: items}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.UI.Fetching {
		t.Error("Update(HistoryLoadedMsg) should set fetching to false")
	}
	if len(updated.historyData) != 1 {
		t.Errorf("Update(HistoryLoadedMsg) historyData len = %d, want 1", len(updated.historyData))
	}
}

func TestUpdateAlbumTracksLoadedMsg(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.UI.Width = 100
	m.UI.Height = 50
	m.UI.Fetching = true
	m.selectedAlbum = &spotify.SimpleAlbum{Name: "Test Album"}

	tracks := []spotify.SimpleTrack{
		{Name: "Album Track 1"},
		{Name: "Album Track 2"},
	}

	msg := AlbumTracksLoadedMsg{Tracks: tracks}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.UI.Fetching {
		t.Error("Update(AlbumTracksLoadedMsg) should set fetching to false")
	}
	if len(updated.albumTracksData) != 2 {
		t.Errorf("Update(AlbumTracksLoadedMsg) albumTracksData len = %d, want 2", len(updated.albumTracksData))
	}
}

func TestUpdateArtistLoadedMsg(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.UI.Width = 100
	m.UI.Height = 50
	m.UI.Fetching = true

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

	if updated.UI.Fetching {
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
	m := NewModel(nil, getTestConfig())
	m.UI.Fetching = true
	m.fetchingLyrics = true

	msg := LyricsLoadedMsg{
		Lyrics:       "Plain lyrics text",
		SyncedLyrics: "",
		TrackName:    "Test Track",
		ArtistName:   "Test Artist",
	}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.UI.Fetching {
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
	if updated.Nav.Current != ViewLyrics {
		t.Errorf("Update(LyricsLoadedMsg) view = %v, want ViewLyrics", updated.Nav.Current)
	}
}

func TestUpdateLyricsLoadedMsgWithSynced(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.UI.Fetching = true

	msg := LyricsLoadedMsg{
		Lyrics:       "Plain lyrics",
		SyncedLyrics: "[00:01.00]Line one\n[00:05.00]Line two",
		TrackName:    "Test Track",
		ArtistName:   "Test Artist",
	}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if !updated.lyricsIsSynced {
		t.Error("Update(LyricsLoadedMsg) should set lyricsIsSynced to true for synced lyrics")
	}
	if len(updated.lyricsSynced) == 0 {
		t.Error("Update(LyricsLoadedMsg) should parse synced lyrics")
	}
}

func TestProgressTickAdvancesOutsideLyrics(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.Playback.IsPlaying = true
	m.Nav.Current = ViewPlaylists
	m.Playback.LocalProgress = 1000
	m.Playback.LastProgressAt = time.Now().Add(-100 * time.Millisecond)

	newModel, _ := m.Update(ProgressTickMsg{})
	updated := newModel.(Model)

	if updated.Playback.LocalProgress <= 1000 {
		t.Errorf("ProgressTickMsg should advance progress in any view, got %d", updated.Playback.LocalProgress)
	}
}

func TestPlaybackStateMsgKeepsOptimisticDuringSeek(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.Playback.SeekPending = true
	m.Playback.LocalProgress = 42000

	msg := PlaybackStateMsg{State: &spotify.PlayerState{
		CurrentlyPlaying: spotify.CurrentlyPlaying{Progress: 1000, Playing: true},
	}}
	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.Playback.LocalProgress != 42000 {
		t.Errorf("PlaybackStateMsg should not clobber optimistic progress during seek, got %d", updated.Playback.LocalProgress)
	}
}

func TestAdjustVolumeOptimistic(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.Playback.State = &spotify.PlayerState{
		Device: spotify.PlayerDevice{ID: "d", Volume: 95},
	}

	m.adjustVolumeOptimistic(10)
	if int(m.Playback.State.Device.Volume) != 100 {
		t.Errorf("adjustVolumeOptimistic should clamp to 100, got %d", int(m.Playback.State.Device.Volume))
	}

	m.adjustVolumeOptimistic(-200)
	if int(m.Playback.State.Device.Volume) != 0 {
		t.Errorf("adjustVolumeOptimistic should clamp to 0, got %d", int(m.Playback.State.Device.Volume))
	}
}

func TestNextRepeatState(t *testing.T) {
	cases := map[string]string{"off": "context", "context": "track", "track": "off", "": "off"}
	for in, want := range cases {
		if got := nextRepeatState(in); got != want {
			t.Errorf("nextRepeatState(%q) = %q, want %q", in, got, want)
		}
	}
}
