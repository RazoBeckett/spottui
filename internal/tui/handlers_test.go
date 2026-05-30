package tui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/zmb3/spotify/v2"

	"github.com/razobeckett/spottui/internal/tui/views"
)

func TestHandleKeyPress_Quit(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.width = 100
	m.height = 50
	m.view = ViewPlaylists

	msg := tea.KeyPressMsg{Code: 'q', Text: "q"}
	_, cmd := m.handleKeyPress(msg)

	if cmd == nil {
		t.Error("handleKeyPress('q') should return quit cmd")
	}
}

func TestHandleKeyPress_Help(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewPlaylists

	msg := tea.KeyPressMsg{Code: '?', Text: "?"}
	newModel, _ := m.handleKeyPress(msg)
	updated := newModel.(Model)

	if updated.view != ViewHelp {
		t.Errorf("handleKeyPress('?') view = %v, want ViewHelp", updated.view)
	}
	if updated.prevView != ViewPlaylists {
		t.Errorf("handleKeyPress('?') prevView = %v, want ViewPlaylists", updated.prevView)
	}
}

func TestHandleKeyPress_HelpToggle(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewHelp
	m.prevView = ViewTracks

	msg := tea.KeyPressMsg{Code: '?', Text: "?"}
	newModel, _ := m.handleKeyPress(msg)
	updated := newModel.(Model)

	if updated.view != ViewTracks {
		t.Errorf("handleKeyPress('?') from help should return to prevView, got %v", updated.view)
	}
}

func TestHandleKeyPress_BackFromTracks(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewTracks

	msg := tea.KeyPressMsg{Code: tea.KeyEsc}
	newModel, _ := m.handleKeyPress(msg)
	updated := newModel.(Model)

	if updated.view != ViewPlaylists {
		t.Errorf("handleKeyPress(Esc) from tracks should go to playlists, got %v", updated.view)
	}
}

func TestHandleKeyPress_BackFromAlbum(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewAlbum
	m.prevView = ViewSearch

	msg := tea.KeyPressMsg{Code: tea.KeyEsc}
	newModel, _ := m.handleKeyPress(msg)
	updated := newModel.(Model)

	if updated.view != ViewSearch {
		t.Errorf("handleKeyPress(Esc) from album should return to prevView, got %v", updated.view)
	}
}

func TestHandleKeyPress_BackFromArtist(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewArtist
	m.prevView = ViewTracks

	msg := tea.KeyPressMsg{Code: tea.KeyEsc}
	newModel, _ := m.handleKeyPress(msg)
	updated := newModel.(Model)

	if updated.view != ViewTracks {
		t.Errorf("handleKeyPress(Esc) from artist should return to prevView, got %v", updated.view)
	}
}

func TestHandleKeyPress_BackFromLyrics(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewLyrics
	m.prevView = ViewPlaylists

	msg := tea.KeyPressMsg{Code: tea.KeyEsc}
	newModel, _ := m.handleKeyPress(msg)
	updated := newModel.(Model)

	if updated.view != ViewPlaylists {
		t.Errorf("handleKeyPress(Esc) from lyrics should return to prevView, got %v", updated.view)
	}
}

func TestHandleKeyPress_BackFromDevices(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewDevices
	m.prevView = ViewPlaylists

	msg := tea.KeyPressMsg{Code: tea.KeyEsc}
	newModel, _ := m.handleKeyPress(msg)
	updated := newModel.(Model)

	if updated.view != ViewPlaylists {
		t.Errorf("handleKeyPress(Esc) from devices should return to prevView, got %v", updated.view)
	}
}

func TestHandleKeyPress_BackFromHistory(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewHistory
	m.prevView = ViewPlaylists

	msg := tea.KeyPressMsg{Code: tea.KeyEsc}
	newModel, _ := m.handleKeyPress(msg)
	updated := newModel.(Model)

	if updated.view != ViewPlaylists {
		t.Errorf("handleKeyPress(Esc) from history should return to prevView, got %v", updated.view)
	}
}

func TestHandleKeyPress_BackFromAddToPlaylist(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewAddToPlaylist
	m.prevView = ViewTracks

	msg := tea.KeyPressMsg{Code: tea.KeyEsc}
	newModel, _ := m.handleKeyPress(msg)
	updated := newModel.(Model)

	if updated.view != ViewTracks {
		t.Errorf("handleKeyPress(Esc) from add to playlist should return to prevView, got %v", updated.view)
	}
}

func TestHandleKeyPress_PlayPause(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewPlaylists

	msg := tea.KeyPressMsg{Code: tea.KeySpace}
	_, cmd := m.handleKeyPress(msg)

	if cmd == nil {
		t.Error("handleKeyPress(Space) should return playback toggle cmd")
	}
}

func TestHandleKeyPress_Next(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewPlaylists

	msg := tea.KeyPressMsg{Code: 'n', Text: "n"}
	_, cmd := m.handleKeyPress(msg)

	if cmd == nil {
		t.Error("handleKeyPress('n') should return next track cmd")
	}
}

func TestHandleKeyPress_Prev(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewPlaylists

	msg := tea.KeyPressMsg{Code: 'p', Text: "p"}
	_, cmd := m.handleKeyPress(msg)

	if cmd == nil {
		t.Error("handleKeyPress('p') should return prev track cmd")
	}
}

func TestHandleKeyPress_VolumeUp(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewPlaylists

	msg := tea.KeyPressMsg{Code: '+', Text: "+"}
	_, cmd := m.handleKeyPress(msg)

	if cmd == nil {
		t.Error("handleKeyPress('+') should return volume up cmd")
	}
}

func TestHandleKeyPress_VolumeDown(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewPlaylists

	msg := tea.KeyPressMsg{Code: '-', Text: "-"}
	_, cmd := m.handleKeyPress(msg)

	if cmd == nil {
		t.Error("handleKeyPress('-') should return volume down cmd")
	}
}

func TestHandleKeyPress_Shuffle(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewPlaylists

	msg := tea.KeyPressMsg{Code: 's', Text: "s"}
	_, cmd := m.handleKeyPress(msg)

	if cmd == nil {
		t.Error("handleKeyPress('s') should return shuffle toggle cmd")
	}
}

func TestHandleKeyPress_Repeat(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewPlaylists

	msg := tea.KeyPressMsg{Code: 'r', Text: "r"}
	_, cmd := m.handleKeyPress(msg)

	if cmd == nil {
		t.Error("handleKeyPress('r') should return repeat cycle cmd")
	}
}

func TestHandleKeyPress_SeekBackward(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewHelp
	m.localProgress = 10000

	msg := tea.KeyPressMsg{Code: '[', Text: "["}
	newModel, cmd := m.handleKeyPress(msg)
	updated := newModel.(Model)

	if !updated.seekPending {
		t.Error("handleKeyPress('[') should set seekPending")
	}
	if updated.pendingSeek != 5000 {
		t.Errorf("handleKeyPress('[') pendingSeek = %d, want 5000", updated.pendingSeek)
	}
	if cmd == nil {
		t.Error("handleKeyPress('[') should return seek tick cmd")
	}
}

func TestHandleKeyPress_SeekForward(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewHelp
	m.localProgress = 10000

	msg := tea.KeyPressMsg{Code: ']', Text: "]"}
	newModel, cmd := m.handleKeyPress(msg)
	updated := newModel.(Model)

	if !updated.seekPending {
		t.Error("handleKeyPress(']') should set seekPending")
	}
	if updated.pendingSeek != 15000 {
		t.Errorf("handleKeyPress(']') pendingSeek = %d, want 15000", updated.pendingSeek)
	}
	if cmd == nil {
		t.Error("handleKeyPress(']') should return seek tick cmd")
	}
}

func TestHandleKeyPress_SeekBackwardClampToZero(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewHelp
	m.localProgress = 2000

	msg := tea.KeyPressMsg{Code: '[', Text: "["}
	newModel, _ := m.handleKeyPress(msg)
	updated := newModel.(Model)

	if updated.pendingSeek != 0 {
		t.Errorf("handleKeyPress('[') should clamp pendingSeek to 0, got %d", updated.pendingSeek)
	}
}

func TestHandleKeyPress_Refresh(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewPlaylists

	msg := tea.KeyPressMsg{Code: 'r', Mod: tea.ModCtrl}
	_, cmd := m.handleKeyPress(msg)

	if cmd == nil {
		t.Error("handleKeyPress(Ctrl+R) should return poll playback cmd")
	}
}

func TestHandleKeyPress_GlobalSearch(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewPlaylists

	msg := tea.KeyPressMsg{Code: 'S', Text: "S"}
	newModel, cmd := m.handleKeyPress(msg)
	updated := newModel.(Model)

	if updated.view != ViewSearch {
		t.Errorf("handleKeyPress('S') view = %v, want ViewSearch", updated.view)
	}
	if updated.prevView != ViewPlaylists {
		t.Errorf("handleKeyPress('S') prevView = %v, want ViewPlaylists", updated.prevView)
	}
	if !updated.searchInput.Focused() {
		t.Error("handleKeyPress('S') should focus search input")
	}
	if cmd == nil {
		t.Error("handleKeyPress('S') should return blink cmd")
	}
}

func TestHandleKeyPress_History(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewPlaylists

	msg := tea.KeyPressMsg{Code: 'H', Text: "H"}
	newModel, cmd := m.handleKeyPress(msg)
	updated := newModel.(Model)

	if updated.view != ViewHistory {
		t.Errorf("handleKeyPress('H') view = %v, want ViewHistory", updated.view)
	}
	if !updated.fetching {
		t.Error("handleKeyPress('H') should set fetching")
	}
	if cmd == nil {
		t.Error("handleKeyPress('H') should return fetch cmd")
	}
}

func TestHandleKeyPress_Devices(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewPlaylists

	msg := tea.KeyPressMsg{Code: 'd', Text: "d"}
	newModel, cmd := m.handleKeyPress(msg)
	updated := newModel.(Model)

	if updated.view != ViewDevices {
		t.Errorf("handleKeyPress('d') view = %v, want ViewDevices", updated.view)
	}
	if !updated.fetching {
		t.Error("handleKeyPress('d') should set fetching")
	}
	if cmd == nil {
		t.Error("handleKeyPress('d') should return fetch cmd")
	}
}

func TestHandleKeyPress_LyricsWithPlayback(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewPlaylists
	m.playbackState = &spotify.PlayerState{
		CurrentlyPlaying: spotify.CurrentlyPlaying{
			Item: &spotify.FullTrack{
				SimpleTrack: spotify.SimpleTrack{
					Name:    "Test Track",
					Artists: []spotify.SimpleArtist{{Name: "Test Artist"}},
				},
			},
		},
	}

	msg := tea.KeyPressMsg{Code: 'L', Text: "L"}
	newModel, cmd := m.handleKeyPress(msg)
	updated := newModel.(Model)

	if !updated.fetchingLyrics {
		t.Error("handleKeyPress('L') should set fetchingLyrics")
	}
	if cmd == nil {
		t.Error("handleKeyPress('L') should return fetch lyrics cmd")
	}
}

func TestHandleKeyPress_LyricsWithoutPlayback(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewPlaylists

	msg := tea.KeyPressMsg{Code: 'L', Text: "L"}
	_, cmd := m.handleKeyPress(msg)

	if cmd != nil {
		t.Error("handleKeyPress('L') without playback should return nil cmd")
	}
}

func TestHandleLyricsKeys_ScrollUp(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewLyrics
	m.lyricsData = "Line 1\nLine 2\nLine 3\nLine 4\nLine 5"
	m.lyricsScrollOffset = 2
	m.height = 20

	msg := tea.KeyPressMsg{Code: tea.KeyUp}
	newModel, _ := m.handleLyricsKeys(msg)
	updated := newModel.(Model)

	if updated.lyricsScrollOffset != 1 {
		t.Errorf("handleLyricsKeys(Up) lyricsScrollOffset = %d, want 1", updated.lyricsScrollOffset)
	}
}

func TestHandleLyricsKeys_ScrollDown(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewLyrics
	m.lyricsData = "Line 1\nLine 2\nLine 3\nLine 4\nLine 5\nLine 6\nLine 7\nLine 8\nLine 9\nLine 10\nLine 11\nLine 12\nLine 13\nLine 14\nLine 15\nLine 16\nLine 17\nLine 18\nLine 19\nLine 20"
	m.lyricsScrollOffset = 0
	m.height = 10

	msg := tea.KeyPressMsg{Code: tea.KeyDown}
	newModel, _ := m.handleLyricsKeys(msg)
	updated := newModel.(Model)

	if updated.lyricsScrollOffset != 1 {
		t.Errorf("handleLyricsKeys(Down) lyricsScrollOffset = %d, want 1", updated.lyricsScrollOffset)
	}
}

func TestHandleLyricsKeys_ScrollUpAtTop(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewLyrics
	m.lyricsData = "Line 1\nLine 2"
	m.lyricsScrollOffset = 0
	m.height = 20

	msg := tea.KeyPressMsg{Code: tea.KeyUp}
	newModel, _ := m.handleLyricsKeys(msg)
	updated := newModel.(Model)

	if updated.lyricsScrollOffset != 0 {
		t.Errorf("handleLyricsKeys(Up) at top should stay at 0, got %d", updated.lyricsScrollOffset)
	}
}

func TestHandleArtistKeys_TabToggle(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewArtist
	m.artistViewMode = "tracks"
	m.width = 100
	m.height = 50
	m.artistTopTracks = views.CreateArtistTopTracksList(nil, m.styles, "", 96, 38)
	m.artistAlbums = views.CreateArtistAlbumsList(nil, m.styles, 96, 38)

	msg := tea.KeyPressMsg{Code: tea.KeyTab}
	newModel, _ := m.handleArtistKeys(msg)
	updated := newModel.(Model)

	if updated.artistViewMode != "albums" {
		t.Errorf("handleArtistKeys(Tab) artistViewMode = %s, want 'albums'", updated.artistViewMode)
	}
}

func TestHandleArtistKeys_TabToggleBack(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewArtist
	m.artistViewMode = "albums"
	m.width = 100
	m.height = 50
	m.artistTopTracks = views.CreateArtistTopTracksList(nil, m.styles, "", 96, 38)
	m.artistAlbums = views.CreateArtistAlbumsList(nil, m.styles, 96, 38)

	msg := tea.KeyPressMsg{Code: tea.KeyTab}
	newModel, _ := m.handleArtistKeys(msg)
	updated := newModel.(Model)

	if updated.artistViewMode != "tracks" {
		t.Errorf("handleArtistKeys(Tab) artistViewMode = %s, want 'tracks'", updated.artistViewMode)
	}
}

func TestHandleArtistKeys_Back(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewArtist
	m.prevView = ViewSearch
	m.artistViewMode = "tracks"
	m.width = 100
	m.height = 50
	m.artistTopTracks = views.CreateArtistTopTracksList(nil, m.styles, "", 96, 38)
	m.artistAlbums = views.CreateArtistAlbumsList(nil, m.styles, 96, 38)

	msg := tea.KeyPressMsg{Code: tea.KeyEsc}
	newModel, _ := m.handleArtistKeys(msg)
	updated := newModel.(Model)

	if updated.view != ViewSearch {
		t.Errorf("handleArtistKeys(Esc) view = %v, want ViewSearch", updated.view)
	}
}

func TestHasSearchResults_Empty(t *testing.T) {
	m := NewModel(nil, getTestConfig())

	if m.hasSearchResults() {
		t.Error("hasSearchResults() should return false when no results")
	}
}

func TestHasSearchResults_WithTracks(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.searchTracksData = []spotify.FullTrack{{}}

	if !m.hasSearchResults() {
		t.Error("hasSearchResults() should return true when tracks exist")
	}
}

func TestHasSearchResults_WithAlbums(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.searchAlbumsData = []spotify.SimpleAlbum{{}}

	if !m.hasSearchResults() {
		t.Error("hasSearchResults() should return true when albums exist")
	}
}

func TestHasSearchResults_WithPlaylists(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.searchPlaylistsData = []spotify.SimplePlaylist{{}}

	if !m.hasSearchResults() {
		t.Error("hasSearchResults() should return true when playlists exist")
	}
}

func TestHasSearchResults_WithArtists(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.searchArtistsData = []spotify.FullArtist{{}}

	if !m.hasSearchResults() {
		t.Error("hasSearchResults() should return true when artists exist")
	}
}

func TestHandleSearchKeys_FocusOnSlash(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewSearch
	m.searchInput.Blur()

	msg := tea.KeyPressMsg{Code: '/', Text: "/"}
	newModel, cmd := m.handleSearchKeys(msg)
	updated := newModel.(Model)

	if !updated.searchInput.Focused() {
		t.Error("handleSearchKeys('/') should focus search input")
	}
	if cmd == nil {
		t.Error("handleSearchKeys('/') should return blink cmd")
	}
}

func TestHandleSearchKeys_FocusOnI(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewSearch
	m.searchInput.Blur()

	msg := tea.KeyPressMsg{Code: 'i', Text: "i"}
	newModel, cmd := m.handleSearchKeys(msg)
	updated := newModel.(Model)

	if !updated.searchInput.Focused() {
		t.Error("handleSearchKeys('i') should focus search input")
	}
	if cmd == nil {
		t.Error("handleSearchKeys('i') should return blink cmd")
	}
}

func TestHandleSearchKeys_EscBlursInput(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewSearch
	m.searchInput.Focus()

	msg := tea.KeyPressMsg{Code: 'q', Text: "q"}
	newModel, _ := m.handleSearchKeys(msg)
	updated := newModel.(Model)

	if updated.searchInput.Focused() {
		t.Error("handleSearchKeys('q') should blur search input")
	}
}

func TestHandleSearchKeys_EnterWithQuery(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewSearch
	m.searchInput.Focus()
	m.searchInput.SetValue("test query")

	msg := tea.KeyPressMsg{Code: tea.KeyEnter}
	newModel, cmd := m.handleSearchKeys(msg)
	updated := newModel.(Model)

	if updated.searchInput.Focused() {
		t.Error("handleSearchKeys(Enter) should blur search input")
	}
	if !updated.searching {
		t.Error("handleSearchKeys(Enter) should set searching to true")
	}
	if cmd == nil {
		t.Error("handleSearchKeys(Enter) should return search cmd")
	}
}

func TestHandleSearchKeys_EnterWithEmptyQuery(t *testing.T) {
	m := NewModel(nil, getTestConfig())
	m.view = ViewSearch
	m.searchInput.Focus()
	m.searchInput.SetValue("")

	msg := tea.KeyPressMsg{Code: tea.KeyEnter}
	newModel, cmd := m.handleSearchKeys(msg)
	updated := newModel.(Model)

	if !updated.searching {
		if cmd != nil {
			t.Error("handleSearchKeys(Enter) with empty query should return nil cmd")
		}
	}
}
