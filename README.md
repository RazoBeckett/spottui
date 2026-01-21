# SpotTUI - Spotify Terminal UI

A beautiful terminal-based Spotify client built with Go, using the Charmbracelet ecosystem (Bubble Tea, Lip Gloss, Bubbles).

> [!NOTE]
> This project is under active development. Features may change and bugs may exist. Contributions and feedback are welcome!

![SpotTUI Demo](demo.gif)

## Features

- **Browse**: Playlists, tracks, albums, and artists
- **Search**: Global search for tracks, albums, playlists, and artists
- **History**: Recently played tracks
- **Lyrics**: View synced and plain lyrics for current track
- **Playback Controls**: Play/pause, next, previous, volume, shuffle, repeat
- **Seeking**: Jump forward/backward 5 seconds
- **Now Playing**: Live display with progress bar and synced lyrics highlighting
- **Device Management**: Switch playback devices on the fly
- **Fuzzy Filtering**: Filter playlists and tracks quickly
- **Vim-style Keybindings**: Efficient keyboard-first navigation
- **Responsive Layout**: Adapts to terminal size with minimum size warnings
- **Smart Caching**: TTL-based API caching for faster navigation
- **Error Handling**: Automatic retry with exponential backoff
- **Spotify-inspired Theme**: Dark green color scheme

## Prerequisites

- Go 1.25+
- Spotify Premium account (required for playback control)
- Spotify Developer App credentials

## Setup

### 1. Create Spotify App

1. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Create a new app
3. Add `http://localhost:8080/callback` as a Redirect URI in your app settings
4. Copy your **Client ID** (you don't need the client secret for PKCE flow)

### 2. Configure SpotTUI

SpotTUI uses a configuration file at `~/.config/spottui/config.json` (following the XDG Base Directory Specification).

**Option A: Create config file**

```bash
mkdir -p ~/.config/spottui
cp config.example.json ~/.config/spottui/config.json
```

Edit the config file and add your Client ID:

```json
{
  "client_id": "your_client_id_here"
}
```

**Option B: Use environment variables** (legacy method)

```bash
cp .env.example .env
# Edit .env and add your SPOTIFY_CLIENT
```

Environment variables (`SPOTIFY_CLIENT`, `SPOTIFY_CALLBACK`) still work and take precedence over the config file.

### 3. Build and Run

```bash
go build -o spottui .
./spottui
```

Or run directly:

```bash
go run main.go
```

On first run, your browser will open for Spotify authentication. After authorizing, the token is cached for future sessions.

## Configuration

SpotTUI uses a JSON configuration file at `~/.config/spottui/config.json`. See `config.example.json` for all available options.

### Configuration Options

| Setting | Default | Description |
|---------|---------|-------------|
| `client_id` | - | Spotify client ID |
| `callback_url` | `http://localhost:8080/callback` | OAuth callback URL |
| `api_timeout` | 15 | API request timeout in seconds |
| `cache_enabled` | true | Enable TTL-based caching |
| `playlist_tracks_ttl` | `5m` | Playlist tracks cache duration |
| `album_tracks_ttl` | `10m` | Album tracks cache duration |
| `artist_data_ttl` | `10m` | Artist data cache duration |
| `search_results_ttl` | `2m` | Search results cache duration |
| `retry_enabled` | true | Enable retry with exponential backoff |
| `max_retries` | 3 | Maximum retry attempts |
| `base_delay` | `500ms` | Initial delay between retries |
| `max_delay` | `5s` | Maximum delay between retries |
| `backoff_factor` | 2.0 | Exponential backoff multiplier |
| `theme` | `spotify` | Color theme |
| `default_view` | `playlists` | Default view on startup |
| `show_notifications` | true | Show toast notifications |
| `min_width` | 60 | Minimum terminal width |
| `min_height` | 15 | Minimum terminal height |
| `seek_increment` | 5000 | Seek increment in milliseconds |
| `volume_step` | 10 | Volume step percentage |

### Environment Variables

The following environment variables override config file values:
- `SPOTIFY_CLIENT` - Spotify client ID
- `SPOTIFY_CALLBACK` - OAuth callback URL

## Keybindings

### Navigation
| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `←` / `h` | Scroll lyrics up (in lyrics view) |
| `→` / `l` | Scroll lyrics down (in lyrics view) |
| `Enter` | Select item |
| `Esc` / `Backspace` | Go back |
| `/` | Filter list |
| `S` | Global search |
| `H` | Recently played |
| `d` | Device selector |
| `A` | View artist |
| `L` | Show lyrics |

### Playback
| Key | Action |
|-----|--------|
| `Space` | Play/Pause |
| `n` / `>` | Next track |
| `p` / `<` | Previous track |
| `+` / `=` | Volume up |
| `-` / `_` | Volume down |
| `s` | Toggle shuffle |
| `r` | Cycle repeat mode (off → playlist → track) |
| `[` | Seek backward 5 seconds |
| `]` | Seek forward 5 seconds |

### Library Management
| Key | Action |
|-----|--------|
| `l` | Like/Unlike current track |
| `a` | Add track to playlist |

### General
| Key | Action |
|-----|--------|
| `?` | Toggle help |
| `Ctrl+R` | Refresh |
| `q` | Quit |

## Project Structure

```
spottui/
├── main.go                     # Entry point, env loading, auth init
├── config.example.json         # Configuration file template
├── .env.example                # Environment variable template
├── internal/
│   ├── auth/
│   │   ├── oauth.go           # PKCE OAuth flow, browser auth
│   │   ├── pkce.go            # Code verifier/challenge generation
│   │   └── token.go           # Token persistence (~/.config/spottui/)
│   ├── config/
│   │   ├── config.go          # Configuration struct and load/save
│   │   └── config_test.go     # Configuration tests
│   └── tui/
│       ├── model.go           # Root Bubble Tea model, Update/View
│       ├── commands.go        # Async tea.Cmd functions (API calls)
│       ├── keys.go            # KeyMap definitions
│       ├── handlers.go        # Key press handlers for all views
│       ├── render.go          # Render functions for all views
│       ├── retry.go           # Retry logic with exponential backoff
│       ├── cache.go           # TTL-based API response caching
│       ├── lrc.go             # LRC lyrics parsing
│       ├── styles/
│       │   └── theme.go       # Spotify-themed Lip Gloss styles
│       └── views/
│           ├── list.go        # List delegates (playlist, track, device, search)
│           └── player.go      # Now playing component with progress bar
├── demo.tape                   # VHS demo script
├── go.mod
└── go.sum
```

## How It Works

- **OAuth PKCE Flow**: Secure authentication without client secret storage
- **Configuration**: JSON config file at `~/.config/spottui/config.json` with XDG standard path
- **Bubble Tea**: Elm-architecture TUI framework for state management
- **Bubbles**: Pre-built components (lists, spinners)
- **Lip Gloss**: Terminal styling
- **zmb3/spotify**: Spotify Web API client library
- **Lyrics Integration**: Fetches synced (LRC) and plain lyrics from external API
- **Smart Caching**: Reduces API calls with TTL-based caching for playlists, albums, artists, and search
- **Error Resilience**: Automatic retry with exponential backoff on API failures

## Generate Demo GIF

To regenerate the demo GIF, install [VHS](https://github.com/charmbracelet/vhs) and run:

```bash
vhs demo.tape
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run specific package tests
go test ./internal/tui/...
go test ./internal/auth/...
```

## Token & Configuration Storage

- **Tokens**: Persisted to `~/.config/spottui/token.json` with 0600 permissions. Automatically refreshed using `autoSaveTokenSource`.
- **Configuration**: Stored at `~/.config/spottui/config.json`. See `config.example.json` for all available options.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run `go test ./...` to ensure all tests pass
6. Submit a pull request

## Legal

- **[Privacy Policy](PRIVACY.md)** - Data collection, storage, and user rights
- **[User Agreement](USER_AGREEMENT.md)** - Terms of use for using SpotTUI
- **[License](LICENSE)** - MIT License for this project

### Disclaimer

[Spotify](https://www.spotify.com/) is a trademark of [Spotify AB](https://www.spotify.com/us/about-us/contact/). This project is not affiliated with or endorsed by Spotify.
