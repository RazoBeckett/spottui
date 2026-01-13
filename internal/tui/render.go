package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/razobeckett/spottui/internal/tui/views"
)

func (m Model) renderLoading() string {
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		m.styles.Spinner.Render(m.spinner.View()+" Loading your Spotify data..."),
	)
}

func (m Model) renderPlaylists() string {
	header := m.renderHeader()
	notification := m.renderNotification()
	player := views.RenderNowPlaying(m.playbackState, m.styles, m.width)
	help := m.renderHelpBar()

	headerHeight := lipgloss.Height(header)
	notificationHeight := lipgloss.Height(notification)
	playerHeight := lipgloss.Height(player)
	helpHeight := lipgloss.Height(help)
	contentHeight := m.height - headerHeight - notificationHeight - playerHeight - helpHeight

	if contentHeight < 1 {
		contentHeight = 1
	}

	content := lipgloss.NewStyle().Height(contentHeight).Render(m.playlists.View())

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		content,
		notification,
		player,
		help,
	)
}

func (m Model) renderTracks() string {
	header := m.renderHeader()
	notification := m.renderNotification()
	player := views.RenderNowPlaying(m.playbackState, m.styles, m.width)
	help := m.renderHelpBar()

	headerHeight := lipgloss.Height(header)
	notificationHeight := lipgloss.Height(notification)
	playerHeight := lipgloss.Height(player)
	helpHeight := lipgloss.Height(help)
	contentHeight := m.height - headerHeight - notificationHeight - playerHeight - helpHeight

	if contentHeight < 1 {
		contentHeight = 1
	}

	content := lipgloss.NewStyle().Height(contentHeight).Render(m.tracks.View())

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		content,
		notification,
		player,
		help,
	)
}

func (m Model) renderAlbum() string {
	header := m.renderHeader()
	notification := m.renderNotification()
	player := views.RenderNowPlaying(m.playbackState, m.styles, m.width)
	help := m.renderHelpBar()

	headerHeight := lipgloss.Height(header)
	notificationHeight := lipgloss.Height(notification)
	playerHeight := lipgloss.Height(player)
	helpHeight := lipgloss.Height(help)
	contentHeight := m.height - headerHeight - notificationHeight - playerHeight - helpHeight

	if contentHeight < 1 {
		contentHeight = 1
	}

	content := lipgloss.NewStyle().Height(contentHeight).Render(m.albumTracks.View())

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		content,
		notification,
		player,
		help,
	)
}

func (m Model) renderDevices() string {
	header := m.renderHeader()
	help := m.styles.HelpBar.Render("↑/↓ navigate • enter select • esc back")

	headerHeight := lipgloss.Height(header)
	helpHeight := lipgloss.Height(help)
	contentHeight := m.height - headerHeight - helpHeight

	if contentHeight < 1 {
		contentHeight = 1
	}

	content := lipgloss.NewStyle().Height(contentHeight).Render(m.devices.View())

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		content,
		help,
	)
}

func (m Model) renderSearch() string {
	header := m.renderHeader()

	inputStyle := m.styles.Header.Copy().Padding(0, 1)
	searchBox := inputStyle.Render("🔍 " + m.searchInput.View())

	notification := m.renderNotification()
	player := views.RenderNowPlaying(m.playbackState, m.styles, m.width)
	help := m.styles.HelpBar.Render("enter search • ↑/↓ navigate results • esc back")

	headerHeight := lipgloss.Height(header)
	searchBoxHeight := lipgloss.Height(searchBox)
	notificationHeight := lipgloss.Height(notification)
	playerHeight := lipgloss.Height(player)
	helpHeight := lipgloss.Height(help)
	contentHeight := m.height - headerHeight - searchBoxHeight - notificationHeight - playerHeight - helpHeight

	if contentHeight < 1 {
		contentHeight = 1
	}

	var contentView string
	if m.hasSearchResults() {
		contentView = m.searchResults.View()
	} else if m.searching {
		contentView = m.styles.Muted.Render("\n  Searching...")
	} else if m.searchInput.Value() != "" && !m.searchInput.Focused() {
		contentView = m.styles.Muted.Render("\n  No results found")
	} else {
		contentView = m.styles.Muted.Render("\n  Type your query and press Enter to search...")
	}

	content := lipgloss.NewStyle().Height(contentHeight).Render(contentView)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		searchBox,
		content,
		notification,
		player,
		help,
	)
}

func (m Model) renderHistory() string {
	header := m.renderHeader()
	notification := m.renderNotification()
	player := views.RenderNowPlaying(m.playbackState, m.styles, m.width)
	help := m.renderHelpBar()

	headerHeight := lipgloss.Height(header)
	notificationHeight := lipgloss.Height(notification)
	playerHeight := lipgloss.Height(player)
	helpHeight := lipgloss.Height(help)
	contentHeight := m.height - headerHeight - notificationHeight - playerHeight - helpHeight

	if contentHeight < 1 {
		contentHeight = 1
	}

	content := lipgloss.NewStyle().Height(contentHeight).Render(m.historyTracks.View())

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		content,
		notification,
		player,
		help,
	)
}

func (m Model) renderArtist() string {
	header := m.renderHeader()

	artistName := "Unknown Artist"
	artistInfo := ""
	if m.selectedArtist != nil {
		artistName = m.selectedArtist.Name
		if len(m.selectedArtist.Genres) > 0 {
			genres := m.selectedArtist.Genres
			if len(genres) > 3 {
				genres = genres[:3]
			}
			artistInfo = " • " + strings.Join(genres, ", ")
		}
	}

	artistHeader := m.styles.ListTitle.Render("󰠃 " + artistName + artistInfo)

	tabStyle := m.styles.Muted
	activeTabStyle := m.styles.ListItemActive
	tracksTab := "Top Tracks"
	albumsTab := "Albums"
	if m.artistViewMode == "tracks" {
		tracksTab = activeTabStyle.Render("[" + tracksTab + "]")
		albumsTab = tabStyle.Render(" " + albumsTab + " ")
	} else {
		tracksTab = tabStyle.Render(" " + tracksTab + " ")
		albumsTab = activeTabStyle.Render("[" + albumsTab + "]")
	}
	tabs := "  " + tracksTab + "  " + albumsTab

	notification := m.renderNotification()
	player := views.RenderNowPlaying(m.playbackState, m.styles, m.width)
	help := m.styles.HelpBar.Render("tab switch view • enter play/select • a view artist • esc back")

	headerHeight := lipgloss.Height(header)
	artistHeaderHeight := lipgloss.Height(artistHeader)
	tabsHeight := 1
	notificationHeight := lipgloss.Height(notification)
	playerHeight := lipgloss.Height(player)
	helpHeight := lipgloss.Height(help)
	contentHeight := m.height - headerHeight - artistHeaderHeight - tabsHeight - notificationHeight - playerHeight - helpHeight

	if contentHeight < 1 {
		contentHeight = 1
	}

	var listView string
	if m.artistViewMode == "tracks" {
		listView = m.artistTopTracks.View()
	} else {
		listView = m.artistAlbums.View()
	}

	content := lipgloss.NewStyle().Height(contentHeight).Render(listView)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		artistHeader,
		tabs,
		content,
		notification,
		player,
		help,
	)
}

func (m Model) renderHeader() string {
	userName := "Spotify User"
	if m.currentUser != nil && m.currentUser.DisplayName != "" {
		userName = m.currentUser.DisplayName
	}

	title := m.styles.Header.Render("♫ SpotTUI")
	user := m.styles.Muted.Render(" - " + userName)

	return title + user
}

func (m Model) renderHelpBar() string {
	return m.styles.HelpBar.Render(
		"↑/↓ navigate • enter select • esc back • space play/pause • ? help",
	)
}

func (m Model) renderNotification() string {
	if m.showError {
		return m.styles.Error.Render("⚠ " + m.errMsg)
	}
	if m.showNotify {
		return m.styles.Success.Render(m.notifyMsg)
	}
	return ""
}

func (m Model) renderHelp() string {
	title := m.styles.DialogTitle.Render("SpotTUI Help")

	navSection := m.styles.ListItemActive.Render("Navigation") + "\n" +
		"  ↑/k       Move up\n" +
		"  ↓/j       Move down\n" +
		"  enter     Select item\n" +
		"  esc       Go back\n" +
		"  /         Filter list\n" +
		"  S         Global search\n" +
		"  A         View artist"

	playbackSection := m.styles.ListItemActive.Render("Playback") + "\n" +
		"  space     Play/Pause\n" +
		"  n/>       Next track\n" +
		"  p/<       Previous track\n" +
		"  [/]       Seek -/+5s\n" +
		"  +/=       Volume up\n" +
		"  -         Volume down\n" +
		"  s         Toggle shuffle\n" +
		"  r         Cycle repeat mode\n" +
		"  l         Like/unlike track\n" +
		"  L         Show lyrics\n" +
		"  a         Add to playlist"

	generalSection := m.styles.ListItemActive.Render("General") + "\n" +
		"  H         Recently played\n" +
		"  d         Device selector\n" +
		"  ?         Toggle help\n" +
		"  ctrl+r    Refresh\n" +
		"  q         Quit"

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		navSection,
		"",
		playbackSection,
		"",
		generalSection,
	)

	footer := m.styles.Muted.Render("Press ? or esc to close")

	dialog := m.styles.Dialog.Render(content)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, dialog, "", footer),
	)
}

func (m Model) renderLyrics() string {
	player := views.RenderNowPlaying(m.playbackState, m.styles, m.width)
	help := m.styles.HelpBar.Render("↑/↓ scroll • esc back")

	playerHeight := lipgloss.Height(player)
	helpHeight := lipgloss.Height(help)

	noLyricsAvailable := m.lyricsData == "" && len(m.lyricsSynced) == 0
	if noLyricsAvailable {
		noLyrics := m.styles.Muted.Render("No lyrics available for this track")
		trackInfo := ""
		if m.lyricsTrackName != "" {
			trackInfo = m.styles.ListTitle.Render(m.lyricsTrackName+" - "+m.lyricsArtistName) + "\n\n"
		}
		content := trackInfo + noLyrics
		contentHeight := m.height - playerHeight - helpHeight
		centeredContent := lipgloss.Place(
			m.width, contentHeight,
			lipgloss.Center, lipgloss.Center,
			content,
		)
		return lipgloss.JoinVertical(lipgloss.Left,
			centeredContent,
			player,
			help,
		)
	}

	header := m.styles.ListTitle.Render("♫ " + m.lyricsTrackName + " - " + m.lyricsArtistName)
	header = lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(header)

	headerHeight := lipgloss.Height(header)
	lyricsAreaHeight := m.height - headerHeight - playerHeight - helpHeight - 2

	if lyricsAreaHeight < 1 {
		lyricsAreaHeight = 1
	}

	if m.lyricsIsSynced && len(m.lyricsSynced) > 0 {
		return m.renderSyncedLyrics(header, player, help, lyricsAreaHeight)
	}

	return m.renderPlainLyrics(header, player, help, lyricsAreaHeight)
}

func (m Model) renderPlainLyrics(header, player, help string, lyricsAreaHeight int) string {
	lines := strings.Split(m.lyricsData, "\n")

	startLine := m.lyricsScrollOffset
	if startLine > len(lines)-lyricsAreaHeight {
		startLine = len(lines) - lyricsAreaHeight
	}
	if startLine < 0 {
		startLine = 0
	}

	endLine := startLine + lyricsAreaHeight
	if endLine > len(lines) {
		endLine = len(lines)
	}

	visibleLines := lines[startLine:endLine]

	var styledLines []string
	for _, line := range visibleLines {
		centered := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(line)
		styledLines = append(styledLines, centered)
	}

	lyricsContent := strings.Join(styledLines, "\n")
	lyricsBox := lipgloss.NewStyle().Height(lyricsAreaHeight).Render(lyricsContent)

	return lipgloss.JoinVertical(lipgloss.Left,
		"",
		header,
		"",
		lyricsBox,
		player,
		help,
	)
}

func (m Model) renderSyncedLyrics(header, player, help string, lyricsAreaHeight int) string {
	currentTimeMs := m.localProgress

	currentLineIdx := 0
	for i, line := range m.lyricsSynced {
		if line.TimeMs <= currentTimeMs {
			currentLineIdx = i
		} else {
			break
		}
	}

	centerOffset := lyricsAreaHeight / 2
	startLine := currentLineIdx - centerOffset
	if startLine < 0 {
		startLine = 0
	}

	endLine := startLine + lyricsAreaHeight
	if endLine > len(m.lyricsSynced) {
		endLine = len(m.lyricsSynced)
		startLine = endLine - lyricsAreaHeight
		if startLine < 0 {
			startLine = 0
		}
	}

	activeStyle := lipgloss.NewStyle().
		Foreground(m.styles.ListItemActive.GetForeground()).
		Bold(true)
	mutedStyle := m.styles.Muted

	var styledLines []string
	for i := startLine; i < endLine; i++ {
		line := m.lyricsSynced[i]
		var styled string
		if i == currentLineIdx {
			styled = activeStyle.Render(line.Text)
		} else {
			styled = mutedStyle.Render(line.Text)
		}
		centered := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(styled)
		styledLines = append(styledLines, centered)
	}

	for len(styledLines) < lyricsAreaHeight {
		styledLines = append(styledLines, "")
	}

	lyricsContent := strings.Join(styledLines, "\n")
	lyricsBox := lipgloss.NewStyle().Height(lyricsAreaHeight).Render(lyricsContent)

	return lipgloss.JoinVertical(lipgloss.Left,
		"",
		header,
		"",
		lyricsBox,
		player,
		help,
	)
}

func (m Model) renderAddToPlaylist() string {
	header := m.styles.Header.Render("Add to Playlist")
	notification := m.renderNotification()
	player := views.RenderNowPlaying(m.playbackState, m.styles, m.width)
	help := m.renderHelpBar()

	headerHeight := lipgloss.Height(header)
	notificationHeight := lipgloss.Height(notification)
	playerHeight := lipgloss.Height(player)
	helpHeight := lipgloss.Height(help)
	availableHeight := m.height - headerHeight - notificationHeight - playerHeight - helpHeight - 2

	if availableHeight < 5 {
		availableHeight = 5
	}

	m.addToPlaylistList.SetSize(m.width-4, availableHeight)

	var content strings.Builder
	content.WriteString(header)
	content.WriteString("\n")
	content.WriteString(m.addToPlaylistList.View())

	if notification != "" {
		return lipgloss.JoinVertical(lipgloss.Left,
			content.String(),
			notification,
			player,
			help,
		)
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		content.String(),
		player,
		help,
	)
}
