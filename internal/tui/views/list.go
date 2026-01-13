package views

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zmb3/spotify/v2"

	"github.com/razobeckett/spottui/internal/tui/styles"
)

// PlaylistItem implements list.Item for playlists
type PlaylistItem struct {
	Playlist spotify.SimplePlaylist
}

func (i PlaylistItem) Title() string       { return i.Playlist.Name }
func (i PlaylistItem) Description() string { return fmt.Sprintf("%d tracks", i.Playlist.Tracks.Total) }
func (i PlaylistItem) FilterValue() string { return i.Playlist.Name }

// TrackItem implements list.Item for tracks
type TrackItem struct {
	Track spotify.PlaylistTrack
	Index int
}

func (i TrackItem) Title() string {
	if i.Track.Track.Name == "" {
		return "Unknown Track"
	}
	return i.Track.Track.Name
}

func (i TrackItem) Description() string {
	if len(i.Track.Track.Artists) == 0 {
		return "Unknown Artist"
	}
	artists := make([]string, len(i.Track.Track.Artists))
	for j, a := range i.Track.Track.Artists {
		artists[j] = a.Name
	}
	return strings.Join(artists, ", ")
}

func (i TrackItem) FilterValue() string {
	return i.Title() + " " + i.Description()
}

// PlaylistDelegate handles playlist item rendering
type PlaylistDelegate struct {
	Styles styles.Styles
}

func (d PlaylistDelegate) Height() int                             { return 2 }
func (d PlaylistDelegate) Spacing() int                            { return 0 }
func (d PlaylistDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d PlaylistDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(PlaylistItem)
	if !ok {
		return
	}

	var title, desc string
	if index == m.Index() {
		title = d.Styles.ListItemActive.Render("▶ " + i.Title())
		desc = d.Styles.Muted.Render("  " + i.Description())
	} else {
		title = d.Styles.ListItem.Render("  " + i.Title())
		desc = d.Styles.Muted.Render("  " + i.Description())
	}

	fmt.Fprint(w, title+"\n"+desc)
}

// TrackDelegate handles track item rendering
type TrackDelegate struct {
	Styles       styles.Styles
	CurrentTrack string // URI of currently playing track
}

func (d TrackDelegate) Height() int                             { return 2 }
func (d TrackDelegate) Spacing() int                            { return 0 }
func (d TrackDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d TrackDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(TrackItem)
	if !ok {
		return
	}

	isPlaying := string(i.Track.Track.URI) == d.CurrentTrack
	isSelected := index == m.Index()

	var titleStyle, descStyle lipgloss.Style

	if isPlaying || isSelected {
		titleStyle = d.Styles.ListItemActive
	} else {
		titleStyle = d.Styles.ListItem
	}

	descStyle = d.Styles.Muted

	prefix := "  "
	if isPlaying {
		prefix = "♫ "
	} else if isSelected {
		prefix = "▶ "
	}

	title := titleStyle.Render(prefix + i.Title())
	desc := descStyle.Render("  󰳩 " + i.Description())

	fmt.Fprint(w, title+"\n"+desc)
}

// CreatePlaylistList creates a new list model for playlists
func CreatePlaylistList(playlists []spotify.SimplePlaylist, s styles.Styles, width, height int) list.Model {
	items := make([]list.Item, len(playlists))
	for i, p := range playlists {
		items[i] = PlaylistItem{Playlist: p}
	}

	delegate := PlaylistDelegate{Styles: s}

	l := list.New(items, delegate, width, height)
	l.Title = "Your Playlists"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = s.ListTitle
	l.SetShowHelp(false)

	return l
}

// CreateTrackList creates a new list model for tracks
func CreateTrackList(tracks []spotify.PlaylistTrack, s styles.Styles, currentTrackURI string, width, height int) list.Model {
	items := make([]list.Item, len(tracks))
	for i, t := range tracks {
		items[i] = TrackItem{Track: t, Index: i}
	}

	delegate := TrackDelegate{Styles: s, CurrentTrack: currentTrackURI}

	l := list.New(items, delegate, width, height)
	l.Title = "Tracks"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = s.ListTitle
	l.SetShowHelp(false)

	return l
}

type AlbumTrackItem struct {
	Track spotify.SimpleTrack
	Index int
}

func (i AlbumTrackItem) Title() string {
	if i.Track.Name == "" {
		return "Unknown Track"
	}
	return i.Track.Name
}

func (i AlbumTrackItem) Description() string {
	if len(i.Track.Artists) == 0 {
		return "Unknown Artist"
	}
	artists := make([]string, len(i.Track.Artists))
	for j, a := range i.Track.Artists {
		artists[j] = a.Name
	}
	return strings.Join(artists, ", ")
}

func (i AlbumTrackItem) FilterValue() string {
	return i.Title() + " " + i.Description()
}

type AlbumTrackDelegate struct {
	Styles       styles.Styles
	CurrentTrack string
}

func (d AlbumTrackDelegate) Height() int                             { return 2 }
func (d AlbumTrackDelegate) Spacing() int                            { return 0 }
func (d AlbumTrackDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d AlbumTrackDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(AlbumTrackItem)
	if !ok {
		return
	}

	isPlaying := string(i.Track.URI) == d.CurrentTrack
	isSelected := index == m.Index()

	var titleStyle, descStyle lipgloss.Style

	if isPlaying || isSelected {
		titleStyle = d.Styles.ListItemActive
	} else {
		titleStyle = d.Styles.ListItem
	}

	descStyle = d.Styles.Muted

	prefix := "  "
	if isPlaying {
		prefix = "♫ "
	} else if isSelected {
		prefix = "▶ "
	}

	title := titleStyle.Render(prefix + i.Title())
	desc := descStyle.Render("  󰳩 " + i.Description())

	fmt.Fprint(w, title+"\n"+desc)
}

func CreateAlbumTrackList(tracks []spotify.SimpleTrack, s styles.Styles, currentTrackURI string, width, height int) list.Model {
	items := make([]list.Item, len(tracks))
	for i, t := range tracks {
		items[i] = AlbumTrackItem{Track: t, Index: i}
	}

	delegate := AlbumTrackDelegate{Styles: s, CurrentTrack: currentTrackURI}

	l := list.New(items, delegate, width, height)
	l.Title = "Album"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = s.ListTitle
	l.SetShowHelp(false)

	return l
}

type DeviceItem struct {
	Device spotify.PlayerDevice
}

func (i DeviceItem) Title() string {
	name := i.Device.Name
	if i.Device.Active {
		name += " (active)"
	}
	return name
}

func (i DeviceItem) Description() string {
	return string(i.Device.Type)
}

func (i DeviceItem) FilterValue() string {
	return i.Device.Name
}

type DeviceDelegate struct {
	Styles styles.Styles
}

func (d DeviceDelegate) Height() int                             { return 2 }
func (d DeviceDelegate) Spacing() int                            { return 0 }
func (d DeviceDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d DeviceDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(DeviceItem)
	if !ok {
		return
	}

	isActive := i.Device.Active
	isSelected := index == m.Index()

	var titleStyle, descStyle lipgloss.Style

	if isActive || isSelected {
		titleStyle = d.Styles.ListItemActive
	} else {
		titleStyle = d.Styles.ListItem
	}

	descStyle = d.Styles.Muted

	prefix := "  "
	if isActive {
		prefix = "● "
	} else if isSelected {
		prefix = "▶ "
	}

	title := titleStyle.Render(prefix + i.Title())
	desc := descStyle.Render("  " + i.Description())

	fmt.Fprint(w, title+"\n"+desc)
}

func CreateDeviceList(devices []spotify.PlayerDevice, s styles.Styles, width, height int) list.Model {
	items := make([]list.Item, len(devices))
	for i, d := range devices {
		items[i] = DeviceItem{Device: d}
	}

	delegate := DeviceDelegate{Styles: s}

	l := list.New(items, delegate, width, height)
	l.Title = "Select Device"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false)
	l.Styles.Title = s.ListTitle
	l.SetShowHelp(false)

	return l
}

type SearchResultType int

const (
	SearchResultTrack SearchResultType = iota
	SearchResultAlbum
	SearchResultPlaylist
	SearchResultArtist
)

type SearchItem struct {
	Type     SearchResultType
	Track    *spotify.FullTrack
	Album    *spotify.SimpleAlbum
	Playlist *spotify.SimplePlaylist
	Artist   *spotify.FullArtist
}

func (i SearchItem) Title() string {
	switch i.Type {
	case SearchResultTrack:
		if i.Track.Name == "" {
			return "Unknown Track"
		}
		return i.Track.Name
	case SearchResultAlbum:
		if i.Album.Name == "" {
			return "Unknown Album"
		}
		return i.Album.Name
	case SearchResultPlaylist:
		if i.Playlist.Name == "" {
			return "Unknown Playlist"
		}
		return i.Playlist.Name
	case SearchResultArtist:
		if i.Artist.Name == "" {
			return "Unknown Artist"
		}
		return i.Artist.Name
	}
	return "Unknown"
}

func (i SearchItem) Description() string {
	switch i.Type {
	case SearchResultTrack:
		artists := artistNames(i.Track.Artists)
		return artists + " • Song"
	case SearchResultAlbum:
		artists := simpleArtistNames(i.Album.Artists)
		return artists + " • Album"
	case SearchResultPlaylist:
		owner := "Unknown"
		if i.Playlist.Owner.DisplayName != "" {
			owner = i.Playlist.Owner.DisplayName
		}
		return owner + " • Playlist"
	case SearchResultArtist:
		followers := ""
		if i.Artist.Followers.Count > 0 {
			followers = formatFollowers(i.Artist.Followers.Count) + " followers • "
		}
		return followers + "Artist"
	}
	return ""
}

func formatFollowers(count spotify.Numeric) string {
	if count >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(count)/1000000)
	}
	if count >= 1000 {
		return fmt.Sprintf("%.1fK", float64(count)/1000)
	}
	return fmt.Sprintf("%d", count)
}

func (i SearchItem) FilterValue() string {
	return i.Title() + " " + i.Description()
}

func (i SearchItem) URI() string {
	switch i.Type {
	case SearchResultTrack:
		return string(i.Track.URI)
	case SearchResultAlbum:
		return string(i.Album.URI)
	case SearchResultPlaylist:
		return string(i.Playlist.URI)
	case SearchResultArtist:
		return string(i.Artist.URI)
	}
	return ""
}

func artistNames(artists []spotify.SimpleArtist) string {
	if len(artists) == 0 {
		return "Unknown Artist"
	}
	names := make([]string, len(artists))
	for j, a := range artists {
		names[j] = a.Name
	}
	return strings.Join(names, ", ")
}

func simpleArtistNames(artists []spotify.SimpleArtist) string {
	return artistNames(artists)
}

type SearchItemDelegate struct {
	Styles       styles.Styles
	CurrentTrack string
}

func (d SearchItemDelegate) Height() int                             { return 2 }
func (d SearchItemDelegate) Spacing() int                            { return 0 }
func (d SearchItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d SearchItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(SearchItem)
	if !ok {
		return
	}

	isPlaying := i.URI() == d.CurrentTrack
	isSelected := index == m.Index()

	var titleStyle, descStyle lipgloss.Style

	if isPlaying || isSelected {
		titleStyle = d.Styles.ListItemActive
	} else {
		titleStyle = d.Styles.ListItem
	}

	descStyle = d.Styles.Muted

	prefix := "  "
	if isPlaying {
		prefix = "♫ "
	} else if isSelected {
		prefix = "▶ "
	}

	title := titleStyle.Render(prefix + i.Title())
	desc := descStyle.Render("  " + i.Description())

	fmt.Fprint(w, title+"\n"+desc)
}

func CreateSearchResultsList(tracks []spotify.FullTrack, albums []spotify.SimpleAlbum, playlists []spotify.SimplePlaylist, artists []spotify.FullArtist, s styles.Styles, currentTrackURI string, width, height int) list.Model {
	var items []list.Item

	for _, t := range tracks {
		track := t
		items = append(items, SearchItem{Type: SearchResultTrack, Track: &track})
	}
	for _, a := range albums {
		album := a
		items = append(items, SearchItem{Type: SearchResultAlbum, Album: &album})
	}
	for _, p := range playlists {
		playlist := p
		items = append(items, SearchItem{Type: SearchResultPlaylist, Playlist: &playlist})
	}
	for _, ar := range artists {
		artist := ar
		items = append(items, SearchItem{Type: SearchResultArtist, Artist: &artist})
	}

	delegate := SearchItemDelegate{Styles: s, CurrentTrack: currentTrackURI}

	l := list.New(items, delegate, width, height)
	l.Title = "Search Results"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false)
	l.Styles.Title = s.ListTitle
	l.SetShowHelp(false)

	return l
}

func CreateHistoryList(items []spotify.RecentlyPlayedItem, s styles.Styles, currentTrackURI string, width, height int) list.Model {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = AlbumTrackItem{Track: item.Track, Index: i}
	}

	delegate := AlbumTrackDelegate{Styles: s, CurrentTrack: currentTrackURI}

	l := list.New(listItems, delegate, width, height)
	l.Title = "Recently Played"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = s.ListTitle
	l.SetShowHelp(false)

	return l
}

type ArtistTopTrackItem struct {
	Track spotify.FullTrack
	Index int
}

func (i ArtistTopTrackItem) Title() string {
	if i.Track.Name == "" {
		return "Unknown Track"
	}
	return i.Track.Name
}

func (i ArtistTopTrackItem) Description() string {
	if i.Track.Album.Name == "" {
		return "Unknown Album"
	}
	return i.Track.Album.Name
}

func (i ArtistTopTrackItem) FilterValue() string {
	return i.Title() + " " + i.Description()
}

type ArtistTopTrackDelegate struct {
	Styles       styles.Styles
	CurrentTrack string
}

func (d ArtistTopTrackDelegate) Height() int                             { return 2 }
func (d ArtistTopTrackDelegate) Spacing() int                            { return 0 }
func (d ArtistTopTrackDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d ArtistTopTrackDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(ArtistTopTrackItem)
	if !ok {
		return
	}

	isPlaying := string(i.Track.URI) == d.CurrentTrack
	isSelected := index == m.Index()

	var titleStyle, descStyle lipgloss.Style

	if isPlaying || isSelected {
		titleStyle = d.Styles.ListItemActive
	} else {
		titleStyle = d.Styles.ListItem
	}

	descStyle = d.Styles.Muted

	prefix := "  "
	if isPlaying {
		prefix = "♫ "
	} else if isSelected {
		prefix = "▶ "
	}

	title := titleStyle.Render(prefix + i.Title())
	desc := descStyle.Render("  󰀥 " + i.Description())

	fmt.Fprint(w, title+"\n"+desc)
}

func CreateArtistTopTracksList(tracks []spotify.FullTrack, s styles.Styles, currentTrackURI string, width, height int) list.Model {
	items := make([]list.Item, len(tracks))
	for i, t := range tracks {
		items[i] = ArtistTopTrackItem{Track: t, Index: i}
	}

	delegate := ArtistTopTrackDelegate{Styles: s, CurrentTrack: currentTrackURI}

	l := list.New(items, delegate, width, height)
	l.Title = "Top Tracks"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false)
	l.Styles.Title = s.ListTitle
	l.SetShowHelp(false)

	return l
}

type ArtistAlbumItem struct {
	Album spotify.SimpleAlbum
}

func (i ArtistAlbumItem) Title() string {
	if i.Album.Name == "" {
		return "Unknown Album"
	}
	return i.Album.Name
}

func (i ArtistAlbumItem) Description() string {
	albumType := string(i.Album.AlbumType)
	if albumType == "" {
		albumType = "Album"
	}
	year := ""
	if i.Album.ReleaseDate != "" && len(i.Album.ReleaseDate) >= 4 {
		year = i.Album.ReleaseDate[:4]
	}
	if year != "" {
		return albumType + " • " + year
	}
	return albumType
}

func (i ArtistAlbumItem) FilterValue() string {
	return i.Album.Name
}

type ArtistAlbumDelegate struct {
	Styles styles.Styles
}

func (d ArtistAlbumDelegate) Height() int                             { return 2 }
func (d ArtistAlbumDelegate) Spacing() int                            { return 0 }
func (d ArtistAlbumDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d ArtistAlbumDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(ArtistAlbumItem)
	if !ok {
		return
	}

	isSelected := index == m.Index()

	var titleStyle, descStyle lipgloss.Style

	if isSelected {
		titleStyle = d.Styles.ListItemActive
	} else {
		titleStyle = d.Styles.ListItem
	}

	descStyle = d.Styles.Muted

	prefix := "  "
	if isSelected {
		prefix = "▶ "
	}

	title := titleStyle.Render(prefix + i.Title())
	desc := descStyle.Render("  " + i.Description())

	fmt.Fprint(w, title+"\n"+desc)
}

func CreateArtistAlbumsList(albums []spotify.SimpleAlbum, s styles.Styles, width, height int) list.Model {
	items := make([]list.Item, len(albums))
	for i, a := range albums {
		items[i] = ArtistAlbumItem{Album: a}
	}

	delegate := ArtistAlbumDelegate{Styles: s}

	l := list.New(items, delegate, width, height)
	l.Title = "Albums"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = s.ListTitle
	l.SetShowHelp(false)

	return l
}
