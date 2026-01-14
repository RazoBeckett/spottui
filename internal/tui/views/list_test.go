package views

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zmb3/spotify/v2"

	"github.com/razobeckett/spottui/internal/tui/styles"
)

func TestPlaylistItem(t *testing.T) {
	item := PlaylistItem{
		Playlist: spotify.SimplePlaylist{
			Name: "Test Playlist",
			Tracks: spotify.PlaylistTracks{
				Total: 42,
			},
		},
	}

	assert.Equal(t, "Test Playlist", item.Title())
	assert.Equal(t, "42 tracks", item.Description())
	assert.Equal(t, "Test Playlist", item.FilterValue())
}

func TestLikedSongsItem(t *testing.T) {
	item := LikedSongsItem{Total: 100}

	assert.Equal(t, "Liked Songs", item.Title())
	assert.Equal(t, "100 tracks", item.Description())
	assert.Equal(t, "Liked Songs", item.FilterValue())
}

func TestTrackItem(t *testing.T) {
	item := TrackItem{
		Track: spotify.PlaylistTrack{
			Track: spotify.FullTrack{
				SimpleTrack: spotify.SimpleTrack{
					Name: "Test Song",
					Artists: []spotify.SimpleArtist{
						{Name: "Artist 1"},
						{Name: "Artist 2"},
					},
				},
			},
		},
		Index: 0,
	}

	assert.Equal(t, "Test Song", item.Title())
	assert.Equal(t, "Artist 1, Artist 2", item.Description())
	assert.Equal(t, "Test Song Artist 1, Artist 2", item.FilterValue())
}

func TestTrackItemEmptyName(t *testing.T) {
	item := TrackItem{
		Track: spotify.PlaylistTrack{
			Track: spotify.FullTrack{
				SimpleTrack: spotify.SimpleTrack{Name: ""},
			},
		},
	}

	assert.Equal(t, "Unknown Track", item.Title())
}

func TestTrackItemNoArtists(t *testing.T) {
	item := TrackItem{
		Track: spotify.PlaylistTrack{
			Track: spotify.FullTrack{
				SimpleTrack: spotify.SimpleTrack{
					Name:    "Test Song",
					Artists: []spotify.SimpleArtist{},
				},
			},
		},
	}

	assert.Equal(t, "Unknown Artist", item.Description())
}

func TestAlbumTrackItem(t *testing.T) {
	item := AlbumTrackItem{
		Track: spotify.SimpleTrack{
			Name:    "Album Track",
			Artists: []spotify.SimpleArtist{{Name: "Album Artist"}},
		},
		Index: 5,
	}

	assert.Equal(t, "Album Track", item.Title())
	assert.Equal(t, "Album Artist", item.Description())
	assert.Equal(t, "Album Track Album Artist", item.FilterValue())
}

func TestAlbumTrackItem_EmptyName(t *testing.T) {
	item := AlbumTrackItem{Track: spotify.SimpleTrack{Name: ""}}
	assert.Equal(t, "Unknown Track", item.Title())
}

func TestAlbumTrackItem_NoArtists(t *testing.T) {
	item := AlbumTrackItem{Track: spotify.SimpleTrack{Name: "Track", Artists: []spotify.SimpleArtist{}}}
	assert.Equal(t, "Unknown Artist", item.Description())
}

func TestSearchItem_Track(t *testing.T) {
	track := &spotify.FullTrack{
		SimpleTrack: spotify.SimpleTrack{
			Name:    "Search Track",
			URI:     "spotify:track:123",
			Artists: []spotify.SimpleArtist{{Name: "Artist"}},
		},
	}
	item := SearchItem{Type: SearchResultTrack, Track: track}

	assert.Equal(t, "Search Track", item.Title())
	assert.Contains(t, item.Description(), "Song")
	assert.Equal(t, "spotify:track:123", item.URI())
}

func TestSearchItem_Album(t *testing.T) {
	album := &spotify.SimpleAlbum{
		Name:    "Test Album",
		URI:     "spotify:album:456",
		Artists: []spotify.SimpleArtist{{Name: "Album Artist"}},
	}
	item := SearchItem{Type: SearchResultAlbum, Album: album}

	assert.Equal(t, "Test Album", item.Title())
	assert.Contains(t, item.Description(), "Album")
	assert.Equal(t, "spotify:album:456", item.URI())
}

func TestSearchItem_Playlist(t *testing.T) {
	playlist := &spotify.SimplePlaylist{
		Name:  "Test Playlist",
		URI:   "spotify:playlist:789",
		Owner: spotify.User{DisplayName: "Owner Name"},
	}
	item := SearchItem{Type: SearchResultPlaylist, Playlist: playlist}

	assert.Equal(t, "Test Playlist", item.Title())
	assert.Contains(t, item.Description(), "Playlist")
	assert.Contains(t, item.Description(), "Owner Name")
	assert.Equal(t, "spotify:playlist:789", item.URI())
}

func TestSearchItem_Artist(t *testing.T) {
	artist := &spotify.FullArtist{
		SimpleArtist: spotify.SimpleArtist{
			Name: "Test Artist",
			URI:  "spotify:artist:abc",
		},
		Followers: spotify.Followers{Count: 1500000},
	}
	item := SearchItem{Type: SearchResultArtist, Artist: artist}

	assert.Equal(t, "Test Artist", item.Title())
	assert.Contains(t, item.Description(), "Artist")
	assert.Contains(t, item.Description(), "1.5M")
	assert.Equal(t, "spotify:artist:abc", item.URI())
}

func TestSearchItem_EmptyNames(t *testing.T) {
	trackItem := SearchItem{Type: SearchResultTrack, Track: &spotify.FullTrack{}}
	assert.Equal(t, "Unknown Track", trackItem.Title())

	albumItem := SearchItem{Type: SearchResultAlbum, Album: &spotify.SimpleAlbum{}}
	assert.Equal(t, "Unknown Album", albumItem.Title())

	playlistItem := SearchItem{Type: SearchResultPlaylist, Playlist: &spotify.SimplePlaylist{}}
	assert.Equal(t, "Unknown Playlist", playlistItem.Title())

	artistItem := SearchItem{Type: SearchResultArtist, Artist: &spotify.FullArtist{}}
	assert.Equal(t, "Unknown Artist", artistItem.Title())
}

func TestDeviceItem(t *testing.T) {
	item := DeviceItem{
		Device: spotify.PlayerDevice{
			Name:   "My Speaker",
			Type:   "Speaker",
			Active: true,
		},
	}

	assert.Equal(t, "My Speaker (active)", item.Title())
	assert.Equal(t, "Speaker", item.Description())
	assert.Equal(t, "My Speaker", item.FilterValue())
}

func TestDeviceItemInactive(t *testing.T) {
	item := DeviceItem{
		Device: spotify.PlayerDevice{
			Name:   "My Speaker",
			Type:   "Speaker",
			Active: false,
		},
	}

	assert.Equal(t, "My Speaker", item.Title())
}

func TestArtistTopTrackItem(t *testing.T) {
	item := ArtistTopTrackItem{
		Track: spotify.FullTrack{
			SimpleTrack: spotify.SimpleTrack{Name: "Hit Song"},
			Album:       spotify.SimpleAlbum{Name: "Greatest Hits"},
		},
		Index: 0,
	}

	assert.Equal(t, "Hit Song", item.Title())
	assert.Equal(t, "Greatest Hits", item.Description())
	assert.Equal(t, "Hit Song Greatest Hits", item.FilterValue())
}

func TestArtistTopTrackItem_Empty(t *testing.T) {
	item := ArtistTopTrackItem{Track: spotify.FullTrack{}}
	assert.Equal(t, "Unknown Track", item.Title())
	assert.Equal(t, "Unknown Album", item.Description())
}

func TestArtistAlbumItem(t *testing.T) {
	item := ArtistAlbumItem{
		Album: spotify.SimpleAlbum{
			Name:        "Album Name",
			AlbumType:   "album",
			ReleaseDate: "2023-05-15",
		},
	}

	assert.Equal(t, "Album Name", item.Title())
	assert.Contains(t, item.Description(), "album")
	assert.Contains(t, item.Description(), "2023")
	assert.Equal(t, "Album Name", item.FilterValue())
}

func TestArtistAlbumItem_Empty(t *testing.T) {
	item := ArtistAlbumItem{Album: spotify.SimpleAlbum{}}
	assert.Equal(t, "Unknown Album", item.Title())
}

func TestArtistAlbumItem_NoYear(t *testing.T) {
	item := ArtistAlbumItem{
		Album: spotify.SimpleAlbum{
			Name:      "Album",
			AlbumType: "single",
		},
	}
	assert.Equal(t, "single", item.Description())
}

func TestFormatFollowers(t *testing.T) {
	tests := []struct {
		count    spotify.Numeric
		expected string
	}{
		{500, "500"},
		{1500, "1.5K"},
		{10000, "10.0K"},
		{1500000, "1.5M"},
		{10000000, "10.0M"},
	}

	for _, tt := range tests {
		result := formatFollowers(tt.count)
		assert.Equal(t, tt.expected, result)
	}
}

func TestArtistNames(t *testing.T) {
	artists := []spotify.SimpleArtist{
		{Name: "Artist A"},
		{Name: "Artist B"},
	}
	result := artistNames(artists)
	assert.Equal(t, "Artist A, Artist B", result)

	result = artistNames([]spotify.SimpleArtist{})
	assert.Equal(t, "Unknown Artist", result)
}

func TestCreatePlaylistList(t *testing.T) {
	s := styles.DefaultStyles()
	playlists := []spotify.SimplePlaylist{
		{Name: "Playlist 1", Tracks: spotify.PlaylistTracks{Total: 10}},
		{Name: "Playlist 2", Tracks: spotify.PlaylistTracks{Total: 20}},
	}

	l := CreatePlaylistList(playlists, 50, s, 80, 20)

	assert.Equal(t, "Your Playlists", l.Title)
	assert.Equal(t, 3, len(l.Items()))
}

func TestCreateTrackList(t *testing.T) {
	s := styles.DefaultStyles()
	tracks := []spotify.PlaylistTrack{
		{Track: spotify.FullTrack{SimpleTrack: spotify.SimpleTrack{Name: "Track 1"}}},
		{Track: spotify.FullTrack{SimpleTrack: spotify.SimpleTrack{Name: "Track 2"}}},
	}

	l := CreateTrackList(tracks, s, "", 80, 20)

	assert.Equal(t, "Tracks", l.Title)
	assert.Equal(t, 2, len(l.Items()))
}

func TestCreateDeviceList(t *testing.T) {
	s := styles.DefaultStyles()
	devices := []spotify.PlayerDevice{
		{Name: "Device 1", Type: "Computer"},
		{Name: "Device 2", Type: "Smartphone"},
	}

	l := CreateDeviceList(devices, s, 80, 20)

	assert.Equal(t, "Select Device", l.Title)
	assert.Equal(t, 2, len(l.Items()))
}

func TestCreateSearchResultsList(t *testing.T) {
	s := styles.DefaultStyles()
	tracks := []spotify.FullTrack{{SimpleTrack: spotify.SimpleTrack{Name: "Track"}}}
	albums := []spotify.SimpleAlbum{{Name: "Album"}}
	playlists := []spotify.SimplePlaylist{{Name: "Playlist"}}
	artists := []spotify.FullArtist{{SimpleArtist: spotify.SimpleArtist{Name: "Artist"}}}

	l := CreateSearchResultsList(tracks, albums, playlists, artists, s, "", 80, 20)

	assert.Equal(t, "Search Results", l.Title)
	assert.Equal(t, 4, len(l.Items()))
}

func TestCreateHistoryList(t *testing.T) {
	s := styles.DefaultStyles()
	items := []spotify.RecentlyPlayedItem{
		{Track: spotify.SimpleTrack{Name: "Recent 1"}},
		{Track: spotify.SimpleTrack{Name: "Recent 2"}},
	}

	l := CreateHistoryList(items, s, "", 80, 20)

	assert.Equal(t, "Recently Played", l.Title)
	assert.Equal(t, 2, len(l.Items()))
}

func TestDelegateHeightAndSpacing(t *testing.T) {
	s := styles.DefaultStyles()

	pd := PlaylistDelegate{Styles: s}
	assert.Equal(t, 2, pd.Height())
	assert.Equal(t, 0, pd.Spacing())

	td := TrackDelegate{Styles: s}
	assert.Equal(t, 2, td.Height())
	assert.Equal(t, 0, td.Spacing())

	dd := DeviceDelegate{Styles: s}
	assert.Equal(t, 2, dd.Height())
	assert.Equal(t, 0, dd.Spacing())

	atd := AlbumTrackDelegate{Styles: s}
	assert.Equal(t, 2, atd.Height())
	assert.Equal(t, 0, atd.Spacing())

	sid := SearchItemDelegate{Styles: s}
	assert.Equal(t, 2, sid.Height())
	assert.Equal(t, 0, sid.Spacing())

	attd := ArtistTopTrackDelegate{Styles: s}
	assert.Equal(t, 2, attd.Height())
	assert.Equal(t, 0, attd.Spacing())

	aad := ArtistAlbumDelegate{Styles: s}
	assert.Equal(t, 2, aad.Height())
	assert.Equal(t, 0, aad.Spacing())
}
