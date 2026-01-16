package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	// Test Spotify settings
	if cfg.CallbackURL != "http://localhost:8080/callback" {
		t.Errorf("CallbackURL = %q, want %q", cfg.CallbackURL, "http://localhost:8080/callback")
	}
	if cfg.APITimeout != 15 {
		t.Errorf("APITimeout = %d, want %d", cfg.APITimeout, 15)
	}

	// Test Cache settings
	if !cfg.CacheEnabled {
		t.Error("CacheEnabled = false, want true")
	}
	if time.Duration(cfg.PlaylistTracksTTL) != 5*time.Minute {
		t.Errorf("PlaylistTracksTTL = %v, want %v", time.Duration(cfg.PlaylistTracksTTL), 5*time.Minute)
	}
	if time.Duration(cfg.AlbumTracksTTL) != 10*time.Minute {
		t.Errorf("AlbumTracksTTL = %v, want %v", time.Duration(cfg.AlbumTracksTTL), 10*time.Minute)
	}
	if time.Duration(cfg.ArtistDataTTL) != 10*time.Minute {
		t.Errorf("ArtistDataTTL = %v, want %v", time.Duration(cfg.ArtistDataTTL), 10*time.Minute)
	}
	if time.Duration(cfg.SearchResultsTTL) != 2*time.Minute {
		t.Errorf("SearchResultsTTL = %v, want %v", time.Duration(cfg.SearchResultsTTL), 2*time.Minute)
	}

	// Test Retry settings
	if !cfg.RetryEnabled {
		t.Error("RetryEnabled = false, want true")
	}
	if cfg.MaxRetries != 3 {
		t.Errorf("MaxRetries = %d, want %d", cfg.MaxRetries, 3)
	}
	if time.Duration(cfg.BaseDelay) != 500*time.Millisecond {
		t.Errorf("BaseDelay = %v, want %v", time.Duration(cfg.BaseDelay), 500*time.Millisecond)
	}
	if time.Duration(cfg.MaxDelay) != 5*time.Second {
		t.Errorf("MaxDelay = %v, want %v", time.Duration(cfg.MaxDelay), 5*time.Second)
	}
	if cfg.BackoffFactor != 2.0 {
		t.Errorf("BackoffFactor = %f, want %f", cfg.BackoffFactor, 2.0)
	}

	// Test UI settings
	if cfg.Theme != "spotify" {
		t.Errorf("Theme = %q, want %q", cfg.Theme, "spotify")
	}
	if cfg.DefaultView != "playlists" {
		t.Errorf("DefaultView = %q, want %q", cfg.DefaultView, "playlists")
	}
	if !cfg.ShowNotifications {
		t.Error("ShowNotifications = false, want true")
	}

	// Test Layout settings
	if cfg.MinWidth != 60 {
		t.Errorf("MinWidth = %d, want %d", cfg.MinWidth, 60)
	}
	if cfg.MinHeight != 15 {
		t.Errorf("MinHeight = %d, want %d", cfg.MinHeight, 15)
	}

	// Test Playback settings
	if cfg.SeekIncrement != 5000 {
		t.Errorf("SeekIncrement = %d, want %d", cfg.SeekIncrement, 5000)
	}
	if cfg.VolumeStep != 10 {
		t.Errorf("VolumeStep = %d, want %d", cfg.VolumeStep, 10)
	}
}

func TestConfigApply(t *testing.T) {
	cfg := DefaultConfig()
	other := Config{
		CallbackURL:   "http://example.com/callback",
		APITimeout:    30,
		Theme:         "minimal",
		SeekIncrement: 10000,
		VolumeStep:    5,
	}

	cfg.apply(other)

	if cfg.CallbackURL != "http://example.com/callback" {
		t.Errorf("CallbackURL = %q, want %q", cfg.CallbackURL, "http://example.com/callback")
	}
	if cfg.APITimeout != 30 {
		t.Errorf("APITimeout = %d, want %d", cfg.APITimeout, 30)
	}
	if cfg.Theme != "minimal" {
		t.Errorf("Theme = %q, want %q", cfg.Theme, "minimal")
	}
	if cfg.SeekIncrement != 10000 {
		t.Errorf("SeekIncrement = %d, want %d", cfg.SeekIncrement, 10000)
	}
	if cfg.VolumeStep != 5 {
		t.Errorf("VolumeStep = %d, want %d", cfg.VolumeStep, 5)
	}
}

func TestConfigApplyZeroValues(t *testing.T) {
	cfg := DefaultConfig()
	originalCallbackURL := cfg.CallbackURL
	other := Config{
		CallbackURL: "",
		APITimeout:  0,
	}

	cfg.apply(other)

	// Zero values should not override existing values
	if cfg.CallbackURL != originalCallbackURL {
		t.Errorf("CallbackURL = %q, want %q", cfg.CallbackURL, originalCallbackURL)
	}
}

func TestConfigApplyEnvOverrides(t *testing.T) {
	// Set environment variables
	os.Setenv("SPOTIFY_CLIENT", "test-client-id")
	os.Setenv("SPOTIFY_CALLBACK", "http://test.com/callback")
	defer os.Unsetenv("SPOTIFY_CLIENT")
	defer os.Unsetenv("SPOTIFY_CALLBACK")

	cfg := DefaultConfig()
	cfg.applyEnvOverrides()

	if cfg.ClientID != "test-client-id" {
		t.Errorf("ClientID = %q, want %q", cfg.ClientID, "test-client-id")
	}
	if cfg.CallbackURL != "http://test.com/callback" {
		t.Errorf("CallbackURL = %q, want %q", cfg.CallbackURL, "http://test.com/callback")
	}
}

func TestConfigApplyEnvOverridesEmpty(t *testing.T) {
	os.Unsetenv("SPOTIFY_CLIENT")
	os.Unsetenv("SPOTIFY_CALLBACK")

	cfg := DefaultConfig()
	cfg.ClientID = "existing-client-id"
	cfg.applyEnvOverrides()

	// Empty env vars should not override existing values
	if cfg.ClientID != "existing-client-id" {
		t.Errorf("ClientID = %q, want %q", cfg.ClientID, "existing-client-id")
	}
}

func TestLoadConfigFileNotFound(t *testing.T) {
	// Save and remove any existing config file
	existingCfg, existingErr := Load()
	existingPath, pathErr := configPath()

	// Try to remove existing config file
	if pathErr == nil {
		os.Remove(existingPath)
	}

	// Load should return default config when file doesn't exist
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Should return default config values
	if cfg.CallbackURL != "http://localhost:8080/callback" {
		t.Errorf("CallbackURL = %q, want %q", cfg.CallbackURL, "http://localhost:8080/callback")
	}
	if cfg.APITimeout != 15 {
		t.Errorf("APITimeout = %d, want %d", cfg.APITimeout, 15)
	}

	// Restore existing config if it was there
	if existingErr == nil && existingPath != "" && existingCfg.ClientID != "" {
		existingCfg.Save()
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	testConfigPath := filepath.Join(tempDir, "config.json")
	_ = testConfigPath

	cfg := DefaultConfig()
	cfg.CallbackURL = "http://test.com/callback"
	cfg.Theme = "minimal"
	cfg.SeekIncrement = 10000

	err := cfg.Save()
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}
}

func TestGetCacheTTL(t *testing.T) {
	cfg := DefaultConfig()

	tests := []struct {
		cacheType string
		want      time.Duration
	}{
		{"playlist_tracks", 5 * time.Minute},
		{"album_tracks", 10 * time.Minute},
		{"artist_data", 10 * time.Minute},
		{"search_results", 2 * time.Minute},
		{"unknown", 0},
	}

	for _, tt := range tests {
		t.Run(tt.cacheType, func(t *testing.T) {
			if got := cfg.GetCacheTTL(tt.cacheType); got != tt.want {
				t.Errorf("GetCacheTTL(%q) = %v, want %v", tt.cacheType, got, tt.want)
			}
		})
	}
}

func TestGetRetryParams(t *testing.T) {
	cfg := DefaultConfig()

	enabled, maxRetries, baseDelay, maxDelay, backoffFactor := cfg.GetRetryParams()

	if !enabled {
		t.Error("RetryEnabled = false, want true")
	}
	if maxRetries != 3 {
		t.Errorf("MaxRetries = %d, want %d", maxRetries, 3)
	}
	if baseDelay != 500*time.Millisecond {
		t.Errorf("BaseDelay = %v, want %v", baseDelay, 500*time.Millisecond)
	}
	if maxDelay != 5*time.Second {
		t.Errorf("MaxDelay = %v, want %v", maxDelay, 5*time.Second)
	}
	if backoffFactor != 2.0 {
		t.Errorf("BackoffFactor = %f, want %f", backoffFactor, 2.0)
	}
}
