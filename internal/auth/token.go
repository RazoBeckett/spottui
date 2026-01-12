package auth

import (
	"encoding/json"
	"os"
	"path/filepath"

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

// Save persists the OAuth token to disk
func (s *TokenStore) Save(token *oauth2.Token) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
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
