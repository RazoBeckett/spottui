package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/oauth2"
)

// TokenStore handles persistent token storage
type TokenStore struct {
	path string
}

// NewTokenStore creates a new token store at ~/.config/spottui/token.json
func NewTokenStore() *TokenStore {
	home, _ := os.UserHomeDir()
	return &TokenStore{
		path: filepath.Join(home, ".config", "spottui", "token.json"),
	}
}

// Save persists the OAuth token to disk using atomic file operations
func (s *TokenStore) Save(token *oauth2.Token) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	// Write to temp file first with restricted permissions
	tempPath := s.path + ".tmp"
	if err := os.WriteFile(tempPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write temp token file: %w", err)
	}

	// Ensure temp file is cleaned up if something goes wrong
	cleanup := func() {
		_ = os.Remove(tempPath)
	}

	// Windows doesn't support atomic rename over existing files
	// Need to remove target first, then rename
	if runtime.GOOS == "windows" {
		// Check if target exists and remove it
		if _, err := os.Stat(s.path); err == nil {
			if err := os.Remove(s.path); err != nil {
				cleanup()
				return fmt.Errorf("failed to remove existing token file: %w", err)
			}
		}
	}

	// Atomically rename temp to target
	if err := os.Rename(tempPath, s.path); err != nil {
		cleanup()
		return fmt.Errorf("failed to rename token file: %w", err)
	}

	return nil
}

// Load retrieves the OAuth token from disk
func (s *TokenStore) Load() (*oauth2.Token, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return nil, err
	}
	var token oauth2.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, err
	}
	return &token, nil
}

// Clear removes the stored token
func (s *TokenStore) Clear() error {
	return os.Remove(s.path)
}
