# SpotTUI - Spotify Terminal UI

A beautiful terminal-based Spotify client built with Go, using the Charmbracelet ecosystem (Bubble Tea, Lip Gloss, Bubbles).

> [!NOTE]
> This project is under active development. Features may change and bugs may exist. Contributions and feedback are welcome!

![SpotTUI Demo](demo.gif)

## Features

- Browse your playlists and tracks
- Global search for tracks, albums, and playlists
- Recently played history
- Playback controls (play/pause, next, previous, volume, shuffle, repeat)
- Live now-playing display with progress bar
- Device selector for switching playback devices
- Fuzzy filtering for playlists and tracks
- Vim-style keybindings
- Spotify-inspired color theme

## Prerequisites

- Go 1.21+
- Spotify Premium account (required for playback control)
- Spotify Developer App credentials

## Setup

### 1. Create Spotify App

1. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Create a new app
3. Add `http://localhost:8080/callback` as a Redirect URI in your app settings
4. Copy your **Client ID** (you don't need the client secret for PKCE flow)

### 2. Configure Environment

Create a `.env` file in the project root:

```bash
cp .env.example .env
```

Edit `.env` and add your Client ID:

```
SPOTIFY_CLIENT=your_client_id_here
SPOTIFY_CALLBACK=http://localhost:8080/callback
```

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

## Keybindings

### Navigation
| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Enter` | Select item |
| `Esc` / `Backspace` | Go back |
| `/` | Filter list |
| `S` | Global search |
| `H` | Recently played |
| `d` | Device selector |

### Playback
| Key | Action |
|-----|--------|
| `Space` | Play/Pause |
| `n` / `>` | Next track |
| `p` / `<` | Previous track |
| `+` / `=` | Volume up |
| `-` | Volume down |
| `s` | Toggle shuffle |
| `r` | Cycle repeat mode (off → playlist → track) |

### General
| Key | Action |
|-----|--------|
| `?` | Toggle help |
| `Ctrl+R` | Refresh |
| `q` | Quit |

## Project Structure

```
spottui/
├── main.go                     # Entry point
├── internal/
│   ├── auth/
│   │   ├── oauth.go           # PKCE OAuth flow
│   │   ├── pkce.go            # PKCE helpers
│   │   └── token.go           # Token storage
│   └── tui/
│       ├── model.go           # Main Bubble Tea model
│       ├── commands.go        # Async commands
│       ├── keys.go            # Keybindings
│       ├── styles/
│       │   └── theme.go       # Lipgloss styles
│       └── views/
│           ├── list.go        # List components
│           └── player.go      # Now playing component
├── demo.tape                   # VHS demo script
├── .env.example
├── go.mod
└── go.sum
```

## How It Works

- **OAuth PKCE Flow**: Secure authentication without client secret storage
- **Bubble Tea**: Elm-architecture TUI framework for state management
- **Bubbles**: Pre-built components (lists, spinners)
- **Lip Gloss**: Terminal styling
- **zmb3/spotify**: Spotify Web API client library

## Generate Demo GIF

To regenerate the demo GIF, install [VHS](https://github.com/charmbracelet/vhs) and run:

```bash
vhs demo.tape
```

## License

MIT
