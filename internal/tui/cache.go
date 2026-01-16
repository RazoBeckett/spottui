package tui

import (
	"sync"
	"time"

	"github.com/zmb3/spotify/v2"

	"github.com/razobeckett/spottui/internal/config"
)

type cacheEntry[T any] struct {
	data      T
	expiresAt time.Time
}

func (e cacheEntry[T]) isExpired() bool {
	return time.Now().After(e.expiresAt)
}

type Cache struct {
	mu sync.RWMutex

	playlistTracks map[spotify.ID]cacheEntry[[]spotify.PlaylistTrack]
	albumTracks    map[spotify.ID]cacheEntry[[]spotify.SimpleTrack]
	artistData     map[spotify.ID]cacheEntry[ArtistLoadedMsg]
	searchResults  map[string]cacheEntry[SearchResultsMsg]

	playlistTracksTTL time.Duration
	albumTracksTTL    time.Duration
	artistDataTTL     time.Duration
	searchResultsTTL  time.Duration
}

func NewCache(cfg *config.Config) *Cache {
	return &Cache{
		playlistTracks: make(map[spotify.ID]cacheEntry[[]spotify.PlaylistTrack]),
		albumTracks:    make(map[spotify.ID]cacheEntry[[]spotify.SimpleTrack]),
		artistData:     make(map[spotify.ID]cacheEntry[ArtistLoadedMsg]),
		searchResults:  make(map[string]cacheEntry[SearchResultsMsg]),

		playlistTracksTTL: time.Duration(cfg.PlaylistTracksTTL),
		albumTracksTTL:    time.Duration(cfg.AlbumTracksTTL),
		artistDataTTL:     time.Duration(cfg.ArtistDataTTL),
		searchResultsTTL:  time.Duration(cfg.SearchResultsTTL),
	}
}

func (c *Cache) GetPlaylistTracks(id spotify.ID) ([]spotify.PlaylistTrack, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.playlistTracks[id]
	if !ok || entry.isExpired() {
		return nil, false
	}
	return entry.data, true
}

func (c *Cache) SetPlaylistTracks(id spotify.ID, tracks []spotify.PlaylistTrack) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.playlistTracks[id] = cacheEntry[[]spotify.PlaylistTrack]{
		data:      tracks,
		expiresAt: time.Now().Add(c.playlistTracksTTL),
	}
}

func (c *Cache) InvalidatePlaylistTracks(id spotify.ID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.playlistTracks, id)
}

func (c *Cache) GetAlbumTracks(id spotify.ID) ([]spotify.SimpleTrack, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.albumTracks[id]
	if !ok || entry.isExpired() {
		return nil, false
	}
	return entry.data, true
}

func (c *Cache) SetAlbumTracks(id spotify.ID, tracks []spotify.SimpleTrack) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.albumTracks[id] = cacheEntry[[]spotify.SimpleTrack]{
		data:      tracks,
		expiresAt: time.Now().Add(c.albumTracksTTL),
	}
}

func (c *Cache) GetArtistData(id spotify.ID) (ArtistLoadedMsg, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.artistData[id]
	if !ok || entry.isExpired() {
		return ArtistLoadedMsg{}, false
	}
	return entry.data, true
}

func (c *Cache) SetArtistData(id spotify.ID, data ArtistLoadedMsg) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.artistData[id] = cacheEntry[ArtistLoadedMsg]{
		data:      data,
		expiresAt: time.Now().Add(c.artistDataTTL),
	}
}

func (c *Cache) GetSearchResults(query string) (SearchResultsMsg, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.searchResults[query]
	if !ok || entry.isExpired() {
		return SearchResultsMsg{}, false
	}
	return entry.data, true
}

func (c *Cache) SetSearchResults(query string, results SearchResultsMsg) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.searchResults[query] = cacheEntry[SearchResultsMsg]{
		data:      results,
		expiresAt: time.Now().Add(c.searchResultsTTL),
	}
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.playlistTracks = make(map[spotify.ID]cacheEntry[[]spotify.PlaylistTrack])
	c.albumTracks = make(map[spotify.ID]cacheEntry[[]spotify.SimpleTrack])
	c.artistData = make(map[spotify.ID]cacheEntry[ArtistLoadedMsg])
	c.searchResults = make(map[string]cacheEntry[SearchResultsMsg])
}
