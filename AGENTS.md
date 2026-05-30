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

- **TUI Framework**: Bubble Tea v2 (`charm.land/bubbletea/v2`)
- **Styling**: Lip Gloss v2 (`charm.land/lipgloss/v2`)
- **Components**: Bubbles v2 (lists, spinners, text inputs)
- **ANSI-safe text**: `github.com/charmbracelet/x/ansi` (truncation, width)
- **Spotify API**: `github.com/zmb3/spotify/v2`
- **Auth**: OAuth 2.0 PKCE flow (no client secret required)

## Reference Repos

- crush: https://github.com/charmbracelet/crush (reference for Charm UI conventions, gradient rendering, ANSI-safe text handling)
- spotify_player: https://github.com/aome510/spotify-player (preferred resource for UX flows & operational safeguards)
- spotify-tui: https://github.com/Rigellute/spotify-tui

Use these as implementation references when designing protocol handling, UX flows, and operational safeguards.

## Key Patterns

### Bubble Tea Model Structure
- Single root `Model` in `internal/tui/model.go` is the **sole** Bubble Tea model.
- State is grouped into sub-state structs in `internal/tui/state.go`, not a flat field list:
  - `Nav NavState` — active view + back-navigation stack.
  - `UI UIState` — transient presentation state (width/height, error/notify, fetching/searching).
  - `Playback PlaybackState` — now-playing snapshot plus locally tracked progress/seek.
- Domain data (list models, data slices, selected items, lyrics) stays flat on `Model`. Group a new concern into a struct only when it is shared across views or has its own behavior.
- Views are enum states (`ViewLoading`, `ViewPlaylists`, `ViewTracks`, etc.).
- Sub-models are embedded Bubbles components (`list.Model`, `spinner.Model`, `textinput.Model`).

### Navigation (NavState)
- `Nav.Push(v)` enters a view and remembers the current one; `Nav.Pop()` returns to it; `Nav.Previous()` peeks.
- Use `Push`/`Pop` for overlay-style views (Devices, Search, History, Artist, Album, AddToPlaylist, Lyrics, Help). Do **not** track a separate "previous view" field.
- Async view switches (e.g. lyrics opening in `LyricsLoadedMsg`) should `Push` only when not already in that view, to avoid stack growth on refetch.

### Message/Command Pattern
- All Spotify API calls and network IO are wrapped in `tea.Cmd` functions in `commands.go`.
- Results return as typed messages (`TracksLoadedMsg`, `PlaybackStateMsg`, etc.).
- Error handling via `ErrMsg` with auto-dismiss after 3s.
- **Never do IO or expensive work in `Update`** — always defer to a `tea.Cmd`.
- **Never mutate `Model` state inside a command closure.** Commands return messages; `Update` mutates state. (A command closure captures a copy of `Model`; writes are discarded.)
- Return multiple commands with `tea.Batch(...)`.
- Prefer methods on `Model` over inline closures when writing a `tea.Cmd`.

### Optimistic Updates
- User actions that change playback (volume, shuffle, repeat, seek) mutate the local `Playback` snapshot **immediately** in the key handler, then fire the confirming `tea.Cmd`. This keeps the UI responsive and flicker-free.
- The progress bar/time render from `Playback.LocalProgress` (advanced by a single `ProgressTickMsg` loop started at launch), not from the raw 5s-polled `state.Progress`.
- While a seek is pending (`Playback.SeekPending`), `PlaybackStateMsg` must **not** clobber `LocalProgress` with the stale polled value.
- Keep optimistic values consistent with the command (e.g. volume step must match `cfg.VolumeStep`, repeat cycling shares `nextRepeatState`).

### Delegate Pattern for Lists
- Each list type has its own Item + Delegate implementation in `views/list.go`.
- All delegates render via the shared `renderRow` helper — do not duplicate the highlight/prefix/title/desc block per delegate.
- `renderRow` truncates both lines to the list width with `ansi.Truncate`. Long names must never overflow the layout.

### Rendering & ANSI Safety
- Use `github.com/charmbracelet/x/ansi` for any width-bounded or ANSI-aware string work: `ansi.Truncate`, `ansi.StringWidth`, `ansi.Strip`. **Never manipulate styled strings at the byte level.**
- Always account for padding/borders in width calculations.
- Gradient helpers live in `styles/grad.go` (`ApplyForegroundGrad`, `GradientLogo`); use the theme's `GradStart`/`GradEnd` endpoints rather than hardcoded colors.

### Style Organization
- All styles centralized in the `styles.Styles` struct.
- Single theme (`SpotifyTheme`); gradient endpoints exposed as `GradStart`/`GradEnd`.
- Styles are passed to views, not imported ad hoc.

## Conventions

### Naming
- Message types: `*Msg` suffix (e.g., `TracksLoadedMsg`).
- Command functions: verb-based (e.g., `fetchTracks`, `pollPlaybackState`).
- View render functions: `render*` prefix (e.g., `renderPlaylists`).
- State sub-structs and their fields are exported (`Nav`, `UI`, `Playback`) so tests and views can read them.

### Error Handling
- API errors wrapped in `ErrMsg{Err: err}` and returned from commands.
- Non-critical errors (e.g., no active player) return nil or empty state.
- User-facing errors displayed in the styled error box.

### Context Usage
- All network commands use a timeout context (5–15s) derived from `m.ctx`; do not pass bare `context.Background()` into API calls.
- Context stored on `Model` for reuse.

## Build & Run

```bash
# Build
go build -o spottui .

# Run
./spottui
# or
go run main.go
```

## Verification

Before considering a change complete:

```bash
go build ./...
go vet ./...
go test ./...
gofmt -l internal/   # must print nothing
```

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `SPOTIFY_CLIENT` | Yes | Spotify app client ID |
| `SPOTIFY_CALLBACK` | No | OAuth redirect URI (default: `http://localhost:8080/callback`) |

## Token Storage

Tokens persisted to `~/.config/spottui/token.json` with 0600 permissions. Auto-refresh handled by `autoSaveTokenSource`.

## Testing

The project has a substantial test suite (`*_test.go` across `internal/tui`, `views`, `styles`, `auth`, `config`).
- Use standard Go testing; `testify` is available.
- Tests assert directly on grouped state (`m.UI.*`, `m.Playback.*`, `m.Nav.*`) — keep these fields exported and stable.
- Add tests for new features and bug fixes (e.g. width truncation, optimistic-update edge cases).
- Run with `go test ./...`.

## Common Tasks

### Adding a New View
1. Add a view constant to the `View` enum in `model.go`.
2. Add a case to the `renderContent()` switch.
3. Create a `render*` function.
4. Wire navigation via `Nav.Push`/`Nav.Pop` in `handleKeyPress`.

### Adding a New Spotify API Call
1. Create a message type in `commands.go`.
2. Create a `tea.Cmd` method on `Model` that calls the API (with a timeout context) and returns the message.
3. Handle the message in the `Update()` switch (this is where state is mutated).

### Adding a New List Type
1. Create an Item struct implementing `list.Item` in `views/list.go`.
2. Create a Delegate whose `Render` delegates to the shared `renderRow` helper.
3. Create a `Create*List` factory function.
