## Project Snapshot

SpotTUI is a terminal-based Spotify client built with Go 1.25+ using the Charmbracelet TUI ecosystem. It provides playlist browsing, playback controls, search, and device management through a Spotify-themed terminal interface.

This repository is a VERY EARLY WIP. Proposing sweeping changes that improve long-term maintainability is encouraged.

## Core Priorities

1. Performance first.
2. Reliability first.
3. Keep behavior predictable under load and during failures (session restarts, reconnects, partial streams).

If a tradeoff is required, choose correctness and robustness over short-term convenience.

## Maintainability

Long term maintainability is a core priority. If you add new functionality, first check if there is shared logic that can be extracted to a separate module. Duplicate logic across multiple files is a code smell and should be avoided. Don't be afraid to change existing code. Don't take shortcuts by just adding local logic to solve a problem.

## Tech Stack

- **TUI Framework**: Bubble Tea v2
- **Styling**: Lip Gloss
- **Components**: Bubbles (lists, spinners, text inputs)
- **Spotify API**: zmb3/spotify/v2
- **Auth**: OAuth 2.0 PKCE flow (no client secret required)

## Reference Repos

- spotify_player: https://github.com/aome510/spotify-player (prefered resource)
- spotify-tui: https://github.com/Rigellute/spotify-tui

Use these as implementation references when designing protocol handling, UX flows, and operational safeguards.

## Key Patterns

### Bubble Tea Model Structure
- Single root `Model` in `internal/tui/model.go`
- Views are enum states (`ViewLoading`, `ViewPlaylists`, `ViewTracks`, etc.)
- Sub-models embedded for Bubbles components (list.Model, spinner.Model, textinput.Model)

### Message/Command Pattern
- All Spotify API calls wrapped in `tea.Cmd` functions in `commands.go`
- Results returned as typed messages (`TracksLoadedMsg`, `PlaybackStateMsg`, etc.)
- Error handling via `ErrMsg` with auto-dismiss after 3s

### Delegate Pattern for Lists
- Each list type has its own Item + Delegate implementation
- Delegates handle rendering with current track highlighting
- Examples: `PlaylistDelegate`, `TrackDelegate`, `SearchItemDelegate`

### Style Organization
- All styles centralized in `styles.Styles` struct
- Single theme (`SpotifyTheme`) with adaptive colors
- Styles passed to views, not imported directly

## Conventions

### Naming
- Message types: `*Msg` suffix (e.g., `TracksLoadedMsg`)
- Command functions: verb-based (e.g., `fetchTracks`, `pollPlaybackState`)
- View render functions: `render*` prefix (e.g., `renderPlaylists`)

### Error Handling
- API errors wrapped in `ErrMsg{Err: err}` and returned
- Non-critical errors (e.g., no active player) return nil or empty state
- User-facing errors displayed in styled error box

### Context Usage
- Background contexts for user-initiated actions
- Timeout contexts (5-15s) for data fetching operations
- Context stored on Model for reuse

## Build & Run

```bash
# Build
go build -o spottui .

# Run
./spottui
# or
go run main.go
```

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `SPOTIFY_CLIENT` | Yes | Spotify app client ID |
| `SPOTIFY_CALLBACK` | No | OAuth redirect URI (default: `http://localhost:8080/callback`) |

## Token Storage

Tokens persisted to `~/.config/spottui/token.json` with 0600 permissions. Auto-refresh handled by `autoSaveTokenSource`.

## No Tests

This project currently has no test files. When adding tests:
- Use standard Go testing (`*_test.go` files)
- testify is available as a dependency
- Run with `go test ./...`

## Common Tasks

### Adding a New View
1. Add view constant to `View` enum in `model.go`
2. Add case to `View()` switch statement
3. Create `render*` function
4. Add navigation logic in `handleKeyPress`

### Adding a New Spotify API Call
1. Create message type in `commands.go`
2. Create `tea.Cmd` function that calls API and returns message
3. Handle message in `Update()` switch

### Adding a New List Type
1. Create Item struct implementing `list.Item` interface in `views/list.go`
2. Create Delegate struct with `Render` method
3. Create `Create*List` factory function
