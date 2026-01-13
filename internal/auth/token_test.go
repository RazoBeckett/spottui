package auth

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func TestTokenStore(t *testing.T) {
	tmpDir := t.TempDir()
	store := &TokenStore{
		path: filepath.Join(tmpDir, "token.json"),
	}

	token := &oauth2.Token{
		AccessToken:  "test-access-token",
		TokenType:    "Bearer",
		RefreshToken: "test-refresh-token",
		Expiry:       time.Now().Add(time.Hour),
	}

	if err := store.Save(token); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if _, err := os.Stat(store.path); os.IsNotExist(err) {
		t.Fatal("Save() did not create token file")
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.AccessToken != token.AccessToken {
		t.Errorf("Load() AccessToken = %v, want %v", loaded.AccessToken, token.AccessToken)
	}
	if loaded.RefreshToken != token.RefreshToken {
		t.Errorf("Load() RefreshToken = %v, want %v", loaded.RefreshToken, token.RefreshToken)
	}

	if err := store.Clear(); err != nil {
		t.Fatalf("Clear() error = %v", err)
	}

	if _, err := os.Stat(store.path); !os.IsNotExist(err) {
		t.Error("Clear() did not remove token file")
	}
}

func TestTokenStoreLoadNonExistent(t *testing.T) {
	store := &TokenStore{
		path: "/nonexistent/path/token.json",
	}

	_, err := store.Load()
	if err == nil {
		t.Error("Load() expected error for non-existent file")
	}
}
