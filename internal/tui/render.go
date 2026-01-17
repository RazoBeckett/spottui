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

func (m Model) renderTooSmall() string {
	msg := m.styles.Error.Render("Terminal too small") + "\n" +
		m.styles.Muted.Render("Minimum: 60x15")
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		msg,
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

	if contentHeight < MinListHeight {
		contentHeight = MinListHeight
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

	if contentHeight < MinListHeight {
		contentHeight = MinListHeight
	}

	var contentView string
	if m.fetching && len(m.tracksData) == 0 {
		title := "Liked Songs"
		if m.selectedPlaylist != nil {
			title = m.selectedPlaylist.Name
		}
		titleLine := m.styles.ListItemActive.Render("  " + title)
		skeletonCount := (contentHeight - 2) / 2
		skeleton := m.renderSkeleton(skeletonCount, m.listWidth())
		contentView = titleLine + "\n\n" + skeleton
	} else {
		contentView = m.tracks.View()
	}
	content := lipgloss.NewStyle().Height(contentHeight).Render(contentView)

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

	if contentHeight < MinListHeight {
		contentHeight = MinListHeight
	}

	var contentView string
	if m.fetching && len(m.albumTracksData) == 0 {
		title := "Album"
		if m.selectedAlbum != nil {
			title = m.selectedAlbum.Name
		}
		titleLine := m.styles.ListItemActive.Render("  " + title)
		skeletonCount := (contentHeight - 2) / 2
		skeleton := m.renderSkeleton(skeletonCount, m.listWidth())
		contentView = titleLine + "\n\n" + skeleton
	} else {
		contentView = m.albumTracks.View()
	}
	content := lipgloss.NewStyle().Height(contentHeight).Render(contentView)

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

	if contentHeight < MinListHeight {
		contentHeight = MinListHeight
	}

	var contentView string
	if m.fetching && len(m.devicesData) == 0 {
		titleLine := m.styles.ListItemActive.Render("  Devices")
		skeletonCount := (contentHeight - 2) / 2
		skeleton := m.renderSkeleton(skeletonCount, m.listWidth())
		contentView = titleLine + "\n\n" + skeleton
	} else {
		contentView = m.devices.View()
	}
	content := lipgloss.NewStyle().Height(contentHeight).Render(contentView)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		content,
		help,
	)
}

func (m Model) renderSearch() string {
	header := m.renderHeader()

	inputStyle := m.styles.Header.Copy().Padding(0, 1)
	searchBox := inputStyle.Render("  " + m.searchInput.View())

	notification := m.renderNotification()
	player := views.RenderNowPlaying(m.playbackState, m.styles, m.width)
	help := m.styles.HelpBar.Render("enter search • ↑/↓ navigate results • esc back")

	headerHeight := lipgloss.Height(header)
	searchBoxHeight := lipgloss.Height(searchBox)
	notificationHeight := lipgloss.Height(notification)
	playerHeight := lipgloss.Height(player)
	helpHeight := lipgloss.Height(help)
	contentHeight := m.height - headerHeight - searchBoxHeight - notificationHeight - playerHeight - helpHeight

	if contentHeight < MinListHeight {
		contentHeight = MinListHeight
	}

	var contentView string
	if m.hasSearchResults() {
		contentView = m.searchResults.View()
	} else if m.searching {
		titleLine := m.styles.ListItemActive.Render("  Search Results")
		skeletonCount := (contentHeight - 2) / 2
		skeleton := m.renderSkeleton(skeletonCount, m.listWidth())
		contentView = titleLine + "\n\n" + skeleton
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

	if contentHeight < MinListHeight {
		contentHeight = MinListHeight
	}

	var contentView string
	if m.fetching && len(m.historyData) == 0 {
		titleLine := m.styles.ListItemActive.Render("  Recently Played")
		skeletonCount := (contentHeight - 2) / 2
		skeleton := m.renderSkeleton(skeletonCount, m.listWidth())
		contentView = titleLine + "\n\n" + skeleton
	} else {
		contentView = m.historyTracks.View()
	}
	content := lipgloss.NewStyle().Height(contentHeight).Render(contentView)

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

	if contentHeight < MinListHeight {
		contentHeight = MinListHeight
	}

	var listView string
	if m.fetching && len(m.artistTopTracksData) == 0 && len(m.artistAlbumsData) == 0 {
		skeletonCount := (contentHeight - 2) / 2
		listView = m.renderSkeleton(skeletonCount, m.listWidth())
	} else if m.artistViewMode == "tracks" {
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

const asciiLogo = `
┏━┓┏━┓┏━┓╺┳╸╺┳╸╻ ╻╻
┗━┓┣━┛┃ ┃ ┃  ┃ ┃ ┃┃
┗━┛╹  ┗━┛ ╹  ╹ ┗━┛╹`

func (m Model) renderHeader() string {
	userName := "Spotify User"
	if m.currentUser != nil && m.currentUser.DisplayName != "" {
		userName = m.currentUser.DisplayName
	}

	logo := m.styles.Header.Render(asciiLogo)
	user := m.styles.Muted.Render("Welcome, " + userName)

	return lipgloss.JoinVertical(lipgloss.Left, logo, user)
}

func (m Model) renderHelpBar() string {
	return m.styles.HelpBar.Render(
		"↑/↓ navigate • enter select • esc back • space play/pause • ? help",
	)
}

func (m Model) renderSkeleton(count, width int) string {
	var lines []string
	for i := 0; i < count; i++ {
		titleWidth := width / 3
		descWidth := width / 4
		if titleWidth > 30 {
			titleWidth = 30
		}
		if descWidth > 20 {
			descWidth = 20
		}
		title := m.styles.Muted.Render("  " + strings.Repeat("░", titleWidth))
		desc := m.styles.Muted.Render("  " + strings.Repeat("░", descWidth))
		lines = append(lines, title, desc)
	}
	return strings.Join(lines, "\n")
}

func (m Model) renderNotification() string {
	if m.showError {
		return m.styles.Error.Render("⚠ " + m.errMsg)
	}
	if m.showNotify {
		return m.styles.Success.Render(m.notifyMsg)
	}
	if m.fetching {
		dots := strings.Repeat(".", m.fetchingDots+1)
		msg := m.styles.Muted.Render("fetching" + dots)
		return lipgloss.PlaceHorizontal(m.width, lipgloss.Center, msg)
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

	if m.fetching {
		dots := strings.Repeat(".", m.fetchingDots+1)
		msg := m.styles.Muted.Render("fetching" + dots)
		contentHeight := m.height - playerHeight - helpHeight
		centeredContent := lipgloss.Place(
			m.width, contentHeight,
			lipgloss.Center, lipgloss.Center,
			msg,
		)
		return lipgloss.JoinVertical(lipgloss.Left,
			centeredContent,
			player,
			help,
		)
	}

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

	if lyricsAreaHeight < MinListHeight {
		lyricsAreaHeight = MinListHeight
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

	if availableHeight < MinListHeight {
		availableHeight = MinListHeight
	}

	m.addToPlaylistList.SetSize(m.width-HorizontalPad, availableHeight)

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

func (m Model) renderCreatePlaylist() string {
	header := m.styles.Header.Render("Create New Playlist")
	notification := m.renderNotification()
	player := views.RenderNowPlaying(m.playbackState, m.styles, m.width)
	help := m.renderHelpBar()

	publicText := "Private"
	if m.createPlaylistIsPublic {
		publicText = "Public"
	}

	labelStyle := m.styles.Header.Copy().Padding(0, 1)

	var content strings.Builder
	content.WriteString(header)
	content.WriteString("\n\n")
	content.WriteString(labelStyle.Render("Name:"))
	content.WriteString("\n")
	content.WriteString(m.createPlaylistName.View())
	content.WriteString("\n\n")
	content.WriteString(labelStyle.Render("Description (optional):"))
	content.WriteString("\n")
	content.WriteString(m.createPlaylistDesc.View())
	content.WriteString("\n\n")
	content.WriteString(labelStyle.Render("Visibility: " + publicText))
	content.WriteString("\n")
	content.WriteString(m.styles.Muted.Render("Press 'p' to toggle public/private, Enter to save, Esc to cancel"))

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

func (m Model) renderEditPlaylist() string {
	header := m.styles.Header.Render("Edit Playlist")
	notification := m.renderNotification()
	player := views.RenderNowPlaying(m.playbackState, m.styles, m.width)
	help := m.renderHelpBar()

	publicText := "Private"
	if m.editPlaylistIsPublic {
		publicText = "Public"
	}

	labelStyle := m.styles.Header.Copy().Padding(0, 1)

	var content strings.Builder
	content.WriteString(header)
	content.WriteString("\n\n")
	content.WriteString(labelStyle.Render("Name:"))
	content.WriteString("\n")
	content.WriteString(m.editPlaylistName.View())
	content.WriteString("\n\n")
	content.WriteString(labelStyle.Render("Description (optional):"))
	content.WriteString("\n")
	content.WriteString(m.editPlaylistDesc.View())
	content.WriteString("\n\n")
	content.WriteString(labelStyle.Render("Visibility: " + publicText))
	content.WriteString("\n")
	content.WriteString(m.styles.Muted.Render("Press 'p' to toggle public/private, Enter to save, Esc to cancel"))

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
