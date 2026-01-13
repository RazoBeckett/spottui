package views

import (
	"testing"

	"github.com/zmb3/spotify/v2"
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

	if item.Title() != "Test Playlist" {
		t.Errorf("Title() = %q, want %q", item.Title(), "Test Playlist")
	}

	if item.Description() != "42 tracks" {
		t.Errorf("Description() = %q, want %q", item.Description(), "42 tracks")
	}

	if item.FilterValue() != "Test Playlist" {
		t.Errorf("FilterValue() = %q, want %q", item.FilterValue(), "Test Playlist")
	}
}

func TestLikedSongsItem(t *testing.T) {
	item := LikedSongsItem{Total: 100}

	if item.Title() != "Liked Songs" {
		t.Errorf("Title() = %q, want %q", item.Title(), "Liked Songs")
	}

	if item.Description() != "100 tracks" {
		t.Errorf("Description() = %q, want %q", item.Description(), "100 tracks")
	}

	if item.FilterValue() != "Liked Songs" {
		t.Errorf("FilterValue() = %q, want %q", item.FilterValue(), "Liked Songs")
	}
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

	if item.Title() != "Test Song" {
		t.Errorf("Title() = %q, want %q", item.Title(), "Test Song")
	}

	if item.Description() != "Artist 1, Artist 2" {
		t.Errorf("Description() = %q, want %q", item.Description(), "Artist 1, Artist 2")
	}

	expectedFilter := "Test Song Artist 1, Artist 2"
	if item.FilterValue() != expectedFilter {
		t.Errorf("FilterValue() = %q, want %q", item.FilterValue(), expectedFilter)
	}
}

func TestTrackItemEmptyName(t *testing.T) {
	item := TrackItem{
		Track: spotify.PlaylistTrack{
			Track: spotify.FullTrack{
				SimpleTrack: spotify.SimpleTrack{
					Name: "",
				},
			},
		},
	}

	if item.Title() != "Unknown Track" {
		t.Errorf("Title() = %q, want %q", item.Title(), "Unknown Track")
	}
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

	if item.Description() != "Unknown Artist" {
		t.Errorf("Description() = %q, want %q", item.Description(), "Unknown Artist")
	}
}

func TestAlbumTrackItem(t *testing.T) {
	item := AlbumTrackItem{
		Track: spotify.SimpleTrack{
			Name: "Album Track",
			Artists: []spotify.SimpleArtist{
				{Name: "Album Artist"},
			},
		},
		Index: 5,
	}

	if item.Title() != "Album Track" {
		t.Errorf("Title() = %q, want %q", item.Title(), "Album Track")
	}

	if item.Description() != "Album Artist" {
		t.Errorf("Description() = %q, want %q", item.Description(), "Album Artist")
	}
}

func TestSearchItem(t *testing.T) {
	track := &spotify.FullTrack{
		SimpleTrack: spotify.SimpleTrack{
			Name: "Search Track",
		},
	}
	item := SearchItem{
		Type:  SearchResultTrack,
		Track: track,
	}

	if item.Title() != "Search Track" {
		t.Errorf("Title() = %q, want %q", item.Title(), "Search Track")
	}

	expectedFilter := "Search Track Unknown Artist • Song"
	if item.FilterValue() != expectedFilter {
		t.Errorf("FilterValue() = %q, want %q", item.FilterValue(), expectedFilter)
	}
}

func TestDeviceItem(t *testing.T) {
	item := DeviceItem{
		Device: spotify.PlayerDevice{
			Name:   "My Speaker",
			Type:   "Speaker",
			Active: true,
		},
	}

	if item.Title() != "My Speaker (active)" {
		t.Errorf("Title() = %q, want %q", item.Title(), "My Speaker (active)")
	}

	if item.Description() != "Speaker" {
		t.Errorf("Description() = %q, want %q", item.Description(), "Speaker")
	}
}

func TestDeviceItemInactive(t *testing.T) {
	item := DeviceItem{
		Device: spotify.PlayerDevice{
			Name:   "My Speaker",
			Type:   "Speaker",
			Active: false,
		},
	}

	if item.Title() != "My Speaker" {
		t.Errorf("Title() = %q, want %q", item.Title(), "My Speaker")
	}
}
