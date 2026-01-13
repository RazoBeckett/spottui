# SpotTUI Progress Tracker

## Recently Completed (Jan 14, 2026)

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
  - Shows skeleton when `fetching=true` AND list is empty
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
- [ ] Error retry logic - Automatic retry with exponential backoff
- [x] ~~Loading skeletons - Placeholder UI while content loads~~ (completed Jan 14, 2026)

### Tech Debt
- [ ] Increase test coverage - Target 80%+ coverage
- [ ] CI/CD setup - GitHub Actions for build/test/release
- [ ] Documentation - API docs, architecture overview
