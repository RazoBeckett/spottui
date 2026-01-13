package tui

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zmb3/spotify/v2"
)

func TestCache_PlaylistTracks(t *testing.T) {
	cache := NewCache()
	playlistID := spotify.ID("test-playlist")

	_, ok := cache.GetPlaylistTracks(playlistID)
	assert.False(t, ok, "should not find uncached playlist")

	tracks := []spotify.PlaylistTrack{
		{Track: spotify.FullTrack{SimpleTrack: spotify.SimpleTrack{Name: "Track 1"}}},
		{Track: spotify.FullTrack{SimpleTrack: spotify.SimpleTrack{Name: "Track 2"}}},
	}
	cache.SetPlaylistTracks(playlistID, tracks)

	cached, ok := cache.GetPlaylistTracks(playlistID)
	assert.True(t, ok, "should find cached playlist")
	assert.Len(t, cached, 2)
	assert.Equal(t, "Track 1", cached[0].Track.Name)

	cache.InvalidatePlaylistTracks(playlistID)
	_, ok = cache.GetPlaylistTracks(playlistID)
	assert.False(t, ok, "should not find invalidated playlist")
}

func TestCache_AlbumTracks(t *testing.T) {
	cache := NewCache()
	albumID := spotify.ID("test-album")

	_, ok := cache.GetAlbumTracks(albumID)
	assert.False(t, ok, "should not find uncached album")

	tracks := []spotify.SimpleTrack{
		{Name: "Track 1"},
		{Name: "Track 2"},
	}
	cache.SetAlbumTracks(albumID, tracks)

	cached, ok := cache.GetAlbumTracks(albumID)
	assert.True(t, ok, "should find cached album")
	assert.Len(t, cached, 2)
}

func TestCache_ArtistData(t *testing.T) {
	cache := NewCache()
	artistID := spotify.ID("test-artist")

	_, ok := cache.GetArtistData(artistID)
	assert.False(t, ok, "should not find uncached artist")

	data := ArtistLoadedMsg{
		Artist:    &spotify.FullArtist{SimpleArtist: spotify.SimpleArtist{Name: "Test Artist"}},
		TopTracks: []spotify.FullTrack{{SimpleTrack: spotify.SimpleTrack{Name: "Hit Song"}}},
		Albums:    []spotify.SimpleAlbum{{Name: "Album 1"}},
	}
	cache.SetArtistData(artistID, data)

	cached, ok := cache.GetArtistData(artistID)
	assert.True(t, ok, "should find cached artist")
	assert.Equal(t, "Test Artist", cached.Artist.Name)
	assert.Len(t, cached.TopTracks, 1)
}

func TestCache_SearchResults(t *testing.T) {
	cache := NewCache()
	query := "test query"

	_, ok := cache.GetSearchResults(query)
	assert.False(t, ok, "should not find uncached search")

	results := SearchResultsMsg{
		Tracks:  []spotify.FullTrack{{SimpleTrack: spotify.SimpleTrack{Name: "Found Track"}}},
		Albums:  []spotify.SimpleAlbum{{Name: "Found Album"}},
		Artists: []spotify.FullArtist{{SimpleArtist: spotify.SimpleArtist{Name: "Found Artist"}}},
	}
	cache.SetSearchResults(query, results)

	cached, ok := cache.GetSearchResults(query)
	assert.True(t, ok, "should find cached search")
	assert.Len(t, cached.Tracks, 1)
	assert.Equal(t, "Found Track", cached.Tracks[0].Name)
}

func TestCache_Clear(t *testing.T) {
	cache := NewCache()

	cache.SetPlaylistTracks(spotify.ID("p1"), []spotify.PlaylistTrack{})
	cache.SetAlbumTracks(spotify.ID("a1"), []spotify.SimpleTrack{})
	cache.SetArtistData(spotify.ID("ar1"), ArtistLoadedMsg{})
	cache.SetSearchResults("query", SearchResultsMsg{})

	cache.Clear()

	_, ok := cache.GetPlaylistTracks(spotify.ID("p1"))
	assert.False(t, ok)
	_, ok = cache.GetAlbumTracks(spotify.ID("a1"))
	assert.False(t, ok)
	_, ok = cache.GetArtistData(spotify.ID("ar1"))
	assert.False(t, ok)
	_, ok = cache.GetSearchResults("query")
	assert.False(t, ok)
}

func TestCacheEntry_Expiration(t *testing.T) {
	entry := cacheEntry[string]{
		data:      "test",
		expiresAt: time.Now().Add(-1 * time.Second),
	}
	assert.True(t, entry.isExpired(), "past entry should be expired")

	entry = cacheEntry[string]{
		data:      "test",
		expiresAt: time.Now().Add(1 * time.Hour),
	}
	assert.False(t, entry.isExpired(), "future entry should not be expired")
}
