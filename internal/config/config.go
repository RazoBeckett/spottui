package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Duration is a custom type that unmarshals duration strings like "5m" or "10s" from JSON
type Duration time.Duration

// UnmarshalJSON implements json.Unmarshaler for Duration
func (d *Duration) UnmarshalJSON(data []byte) error {
	// Remove quotes
	s := string(data)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	// Parse duration
	duration, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", s, err)
	}

	*d = Duration(duration)
	return nil
}

// MarshalJSON implements json.Marshaler for Duration
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

// Config represents all configurable settings for SpotTUI
type Config struct {
	// Spotify Settings
	ClientID    string `json:"client_id"`    // Spotify client ID (env: SPOTIFY_CLIENT)
	CallbackURL string `json:"callback_url"` // OAuth callback URL (env: SPOTIFY_CALLBACK, default: http://localhost:8080/callback)
	APITimeout  int    `json:"api_timeout"`  // API request timeout in seconds (default: 15)

	// Cache Settings
	CacheEnabled      bool     `json:"cache_enabled"`       // Enable/disable caching (default: true)
	PlaylistTracksTTL Duration `json:"playlist_tracks_ttl"` // Playlist tracks cache TTL (default: 5m)
	AlbumTracksTTL    Duration `json:"album_tracks_ttl"`    // Album tracks cache TTL (default: 10m)
	ArtistDataTTL     Duration `json:"artist_data_ttl"`     // Artist data cache TTL (default: 10m)
	SearchResultsTTL  Duration `json:"search_results_ttl"`  // Search results cache TTL (default: 2m)

	// Retry Settings
	RetryEnabled  bool     `json:"retry_enabled"`  // Enable/disable retry logic (default: true)
	MaxRetries    int      `json:"max_retries"`    // Maximum retry attempts (default: 3)
	BaseDelay     Duration `json:"base_delay"`     // Initial delay between retries (default: 500ms)
	MaxDelay      Duration `json:"max_delay"`      // Maximum delay between retries (default: 5s)
	BackoffFactor float64  `json:"backoff_factor"` // Exponential backoff factor (default: 2.0)

	// UI Settings
	Theme             string `json:"theme"`              // Color theme: "spotify", "minimal", "high-contrast" (default: "spotify")
	DefaultView       string `json:"default_view"`       // Default view on startup: "playlists", "search", "history" (default: "playlists")
	ShowNotifications bool   `json:"show_notifications"` // Show toast notifications (default: true)

	// Layout Settings
	MinWidth  int `json:"min_width"`  // Minimum terminal width (default: 60)
	MinHeight int `json:"min_height"` // Minimum terminal height (default: 15)

	// Playback Settings
	SeekIncrement int `json:"seek_increment"` // Seek increment in milliseconds (default: 5000)
	VolumeStep    int `json:"volume_step"`    // Volume step percentage (default: 10)
}

// DefaultConfig returns the default configuration values
func DefaultConfig() Config {
	return Config{
		// Spotify Settings
		ClientID:    "",
		CallbackURL: "http://localhost:8080/callback",
		APITimeout:  15,

		// Cache Settings
		CacheEnabled:      true,
		PlaylistTracksTTL: Duration(5 * time.Minute),
		AlbumTracksTTL:    Duration(10 * time.Minute),
		ArtistDataTTL:     Duration(10 * time.Minute),
		SearchResultsTTL:  Duration(2 * time.Minute),

		// Retry Settings
		RetryEnabled:  true,
		MaxRetries:    3,
		BaseDelay:     Duration(500 * time.Millisecond),
		MaxDelay:      Duration(5 * time.Second),
		BackoffFactor: 2.0,

		// UI Settings
		Theme:             "spotify",
		DefaultView:       "playlists",
		ShowNotifications: true,

		// Layout Settings
		MinWidth:  60,
		MinHeight: 15,

		// Playback Settings
		SeekIncrement: 5000,
		VolumeStep:    10,
	}
}

// configPath returns the path to the config file
func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".config", "spottui", "config.json"), nil
}

// Load loads configuration from the config file
// Environment variables take precedence over config file values for certain settings
func Load() (Config, error) {
	cfg := DefaultConfig()

	path, err := configPath()
	if err != nil {
		return cfg, err
	}

	// Try to read config file
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Config file doesn't exist, return defaults
			return cfg, nil
		}
		return cfg, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse config file
	var fileCfg Config
	if err := json.Unmarshal(data, &fileCfg); err != nil {
		return cfg, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply file config
	cfg.apply(fileCfg)

	// Environment variables override config file values
	cfg.applyEnvOverrides()

	return cfg, nil
}

// isZeroDuration checks if a Duration is zero
func (d Duration) isZero() bool {
	return time.Duration(d) == 0
}

// apply applies configuration from another Config struct (non-zero values only)
func (c *Config) apply(other Config) {
	if other.ClientID != "" {
		c.ClientID = other.ClientID
	}
	if other.CallbackURL != "" {
		c.CallbackURL = other.CallbackURL
	}
	if other.APITimeout != 0 {
		c.APITimeout = other.APITimeout
	}
	c.CacheEnabled = other.CacheEnabled
	if !other.PlaylistTracksTTL.isZero() {
		c.PlaylistTracksTTL = other.PlaylistTracksTTL
	}
	if !other.AlbumTracksTTL.isZero() {
		c.AlbumTracksTTL = other.AlbumTracksTTL
	}
	if !other.ArtistDataTTL.isZero() {
		c.ArtistDataTTL = other.ArtistDataTTL
	}
	if !other.SearchResultsTTL.isZero() {
		c.SearchResultsTTL = other.SearchResultsTTL
	}
	c.RetryEnabled = other.RetryEnabled
	if other.MaxRetries > 0 {
		c.MaxRetries = other.MaxRetries
	}
	if !other.BaseDelay.isZero() {
		c.BaseDelay = other.BaseDelay
	}
	if !other.MaxDelay.isZero() {
		c.MaxDelay = other.MaxDelay
	}
	if other.BackoffFactor > 0 {
		c.BackoffFactor = other.BackoffFactor
	}
	if other.Theme != "" {
		c.Theme = other.Theme
	}
	if other.DefaultView != "" {
		c.DefaultView = other.DefaultView
	}
	c.ShowNotifications = other.ShowNotifications
	if other.MinWidth > 0 {
		c.MinWidth = other.MinWidth
	}
	if other.MinHeight > 0 {
		c.MinHeight = other.MinHeight
	}
	if other.SeekIncrement > 0 {
		c.SeekIncrement = other.SeekIncrement
	}
	if other.VolumeStep > 0 {
		c.VolumeStep = other.VolumeStep
	}
}

// applyEnvOverrides applies environment variable overrides to config
func (c *Config) applyEnvOverrides() {
	// SPOTIFY_CLIENT overrides client_id
	if val := os.Getenv("SPOTIFY_CLIENT"); val != "" {
		c.ClientID = val
	}

	// SPOTIFY_CALLBACK overrides callback_url
	if val := os.Getenv("SPOTIFY_CALLBACK"); val != "" {
		c.CallbackURL = val
	}
}

// Save saves the configuration to the config file
func (c *Config) Save() error {
	path, err := configPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write config file
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetCacheTTL returns the cache TTL for a given cache type
func (c *Config) GetCacheTTL(cacheType string) time.Duration {
	switch cacheType {
	case "playlist_tracks":
		return time.Duration(c.PlaylistTracksTTL)
	case "album_tracks":
		return time.Duration(c.AlbumTracksTTL)
	case "artist_data":
		return time.Duration(c.ArtistDataTTL)
	case "search_results":
		return time.Duration(c.SearchResultsTTL)
	default:
		return 0
	}
}

// GetRetryParams returns the retry parameters
func (c *Config) GetRetryParams() (enabled bool, maxRetries int, baseDelay, maxDelay time.Duration, backoffFactor float64) {
	return c.RetryEnabled, c.MaxRetries, time.Duration(c.BaseDelay), time.Duration(c.MaxDelay), c.BackoffFactor
}
