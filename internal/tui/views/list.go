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
		title = d.Styles.ListItem.Render(i.Title())
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

	prefix := ""
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
