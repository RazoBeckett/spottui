# SpotTUI Progress Tracker

## Recently Completed (Feb 03, 2026)

- [x] **README keybindings documentation** - @razobeckett
  - Documented Tab key in Navigation keybindings table
  - Added Tab description: "Toggle tabs (in artist view)"
  - Matches implementation in `handlers.go` and help text in `render.go`

## Recently Completed (Jan 16, 2026)

- [x] **Configuration support** - @razobeckett
  - Created `internal/config/config.go` with comprehensive Config struct
  - Config file location: `~/.config/spottui/config.json` (XDG standard)
  - Supports all configurable settings: Spotify credentials, cache TTLs, retry params, UI theme, layout, playback
  - Environment variable override: `SPOTIFY_CLIENT`, `SPOTIFY_CALLBACK`
  - JSON file format with sensible defaults for all settings
  - Custom Duration type for human-readable duration strings ("5m", "10s")
  - Created `internal/config/config_test.go` with 8 test cases
  - Updated `main.go` to use new config system
  - Updated `model.go` to embed config and pass to cache
  - Updated `cache.go` to use config TTL values
  - Updated all test files to use new config parameter
  - Created `config.example.json` for user reference
  - Updated `.env.example` to document new config approach
  - All tests pass: `go test ./...` (6 packages, all passing)
  - Commit: `3dbeec4`

- [x] **README documentation update** - @razobeckett
  - Updated setup section with config file instructions
  - Added comprehensive Configuration section with all options documented
  - Added Configuration to "How It Works" section
  - Updated Project Structure to include config package
  - Renamed Token Storage to "Token & Configuration Storage"
  - Documented environment variable overrides

## Recently Completed (Jan 15, 2026)

- [x] **README documentation update** - @razobeckett
  - Added new features to README: lyrics display, seek forward/backward, artist view, like tracks, add to playlist
  - Updated keybindings table with all new keys ([, ], L, A, l, a)
  - Updated project structure (model.go refactored into model.go, handlers.go, render.go, retry.go, cache.go, lrc.go)
  - Updated Go version requirement to 1.25+
  - Added testing section
  - Added contributing section
  - Added smart caching and error handling to "How It Works" section
  - Added lyrics integration details

- [x] **Auto-fetch lyrics on track change** - @razobeckett
  - Lyrics view now detects when track changes (via next/prev or external control)
  - Compares current playback track with displayed lyrics track/artist
  - Automatically fetches new lyrics when mismatch detected
  - Shows "fetching..." indicator centered above player while loading
  - Files: `internal/tui/model.go`, `internal/tui/render.go`
  - Commit: `7889f5e`

## Recently Completed (Jan 14, 2026)

- [x] **Test coverage improvements** - @razobeckett
  - Added comprehensive tests for model.go, commands.go, handlers.go
  - Added tests for auth/token.go edge cases
  - Coverage improvements:
    - `tui`: 7.5% → 43.9%
    - `views`: 9.1% → 56.7%
    - `auth`: 20.8% → 23.6%
    - `styles`: 100% (maintained)
  - Test files added/modified:
    - `internal/tui/model_test.go` (NEW)
    - `internal/tui/commands_test.go` (NEW)
    - `internal/tui/handlers_test.go` (NEW)
    - `internal/auth/token_test.go` (expanded)

- [x] **Error retry logic with exponential backoff** - @razobeckett
  - Created `retry.go` with generic `withRetry[T]` function
  - Exponential backoff: 500ms base, 2x factor, 5s max, 3 retries
  - Retries on: rate limits, 5xx errors, timeouts, connection issues
  - Added `friendlyError()` for user-friendly error messages
  - Wrapped all data fetching functions: initial data, playlists, tracks, albums, search, history, artist
  - Files: `internal/tui/retry.go`, `internal/tui/commands.go`

- [x] **Optimistic UI updates for shuffle/repeat** - @razobeckett
  - Added `ShuffleToggledMsg` and `RepeatCycledMsg` message types
  - UI updates instantly on toggle, falls back to polling on API failure
  - Eliminates delay when toggling shuffle/repeat modes
  - Files: `internal/tui/commands.go`, `internal/tui/model.go`
  - Commit: `48e9eef`

- [x] **Move shuffle/repeat icons to right side** - @razobeckett
  - Icons now display before volume percentage, always aligned
  - Added `ActiveIcon` style (green when active, muted when inactive)
  - Changed to Unicode symbols (⤮ for shuffle, ⟳ for repeat) for compatibility
  - Files: `internal/tui/views/player.go`, `internal/tui/styles/theme.go`
  - Commit: `7cd9850`

- [x] **Loading skeletons - layout fix** - @razobeckett
  - Fixed nil pointer panic: clear data slices instead of uninitialized list items
  - Added titles to skeleton views for consistent layout (playlist name, album name, etc.)
  - Skeleton views now match loaded state structure
  - Files: `internal/tui/handlers.go`, `internal/tui/render.go`
  - Commit: `6868b93`

- [x] **Loading skeletons fix** - @razobeckett
  - Fixed skeleton not displaying: view now switches BEFORE fetch commands
  - Updated all handlers: playlist, album, artist, search, history, devices
  - Removed redundant view assignments from message handlers
  - Files: `internal/tui/handlers.go`, `internal/tui/model.go`
  - Commit: `7480fbc`

- [x] **Loading skeletons** - @razobeckett
  - Added `renderSkeleton()` function using ░ block characters
  - Skeleton displays placeholder bars mimicking list item structure
  - Integrated into 6 views: tracks, album, devices, search, history, artist
  - Shows skeleton when `fetching=true` AND data is empty
  - Skeleton count adapts to available content height
  - Files: `internal/tui/render.go`

- [x] **ASCII art logo in header** - @razobeckett
  - Replaced "♫ SpotTUI" text with compact Unicode box-drawing ASCII logo
  - 3-line tall, 21-char wide logo using ┏━┓ style characters
  - Added top padding via leading newline in logo constant
  - Files: `internal/tui/render.go`
  - Commit: `95b5cd9`

- [x] **Progress tracking guidelines in AGENTS.md** - @razobeckett
  - Added rule to always update `progress/` directory
  - Documented all 5 progress files and their purposes
  - Included format guidelines for entries
  - Commit: `13e8cd1`

## Recently Completed (Jan 13, 2026)

- [x] **Responsive layout improvements** - @razobeckett
  - Added layout constants (MinWidth=60, MinHeight=15, HorizontalPad, ListHeightSub, MinListHeight)
  - Shows "Terminal too small" warning when below minimum size
  - Replaced hardcoded magic numbers with constants and helper methods
  - All lists now resize properly on window resize
  - Safe width calculations in player component using max()
  - Commit: `ddb907d`

- [x] **Animated loading indicator** - @razobeckett
  - Added `fetching` and `fetchingDots` fields to Model
  - `FetchingTickMsg` cycles dots every 300ms ("fetching.", "fetching..", "fetching...")
  - All fetch calls in handlers.go trigger loading state
  - All result messages clear loading state
  - Center-aligned indicator in notification area
  - Commits: `6a58350`, `4988da1`

- [x] **TTL-based API response caching** - @razobeckett
  - Added `cache.go` with generic `CacheEntry[T]` and `Cache` struct
  - Cached: playlist tracks (5min), album tracks (10min), artist data (10min), search results (2min)
  - Cache invalidation on `addTrackToPlaylist`
  - Comprehensive tests in `cache_test.go`
  - Commit: `2301739`

- [x] **model.go refactor** - @razobeckett
  - Split 1455-line file into 4 focused files:
    - `model.go` (425 lines) - Core Model, Init, Update, View
    - `handlers.go` (527 lines) - All key handlers
    - `render.go` (521 lines) - All render functions
    - `lrc.go` (35 lines) - LRC parsing
  - Commit: `fbaf8cd`

- [x] **Unit tests** - @razobeckett
  - Added tests for auth, tui, styles, and views packages
  - Commit: `bc56899`

- [x] **Seek forward/backward** - @razobeckett
  - `[` and `]` keys for 5s seek
  - Debounced to prevent API rate limiting (100-150ms)
  - Commits: `4bc9b6b`, `9e871b2`, `0ee7e4b`

- [x] **Help bar improvements** - @razobeckett
  - Simplified to essential keybinds
  - Removed double box styling
  - Commits: `077d56e`, `5e72739`

---

## In Progress

_None currently_

---

## Backlog

### Features
- [ ] Queue view - View and manage playback queue
- [ ] Playlist creation/editing - Create new playlists, add/remove tracks
- [ ] Album art display - Sixel/Kitty graphics protocol support
- [ ] Keyboard shortcuts customization - User-configurable keybindings
- [ ] Offline mode handling - Graceful degradation when no connection

### Polish
- [x] ~~Responsive layout - Better handling of small terminal sizes~~ (completed Jan 13, 2026)
- [x] ~~Error retry logic - Automatic retry with exponential backoff~~ (completed Jan 14, 2026)
- [x] ~~Loading skeletons - Placeholder UI while content loads~~ (completed Jan 14, 2026)

### Tech Debt
- [ ] Increase test coverage - Target 80%+ coverage
- [ ] CI/CD setup - GitHub Actions for build/test/release
- [ ] Documentation - API docs, architecture overview
