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

func TestTokenStoreSaveCreatesDir(t *testing.T) {
	tmpDir := t.TempDir()
	store := &TokenStore{
		path: filepath.Join(tmpDir, "nested", "dir", "token.json"),
	}

	token := &oauth2.Token{
		AccessToken: "test-token",
	}

	if err := store.Save(token); err != nil {
		t.Fatalf("Save() should create nested directories, got error = %v", err)
	}

	if _, err := os.Stat(store.path); os.IsNotExist(err) {
		t.Error("Save() did not create token file in nested directory")
	}
}

func TestTokenStoreLoadInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	tokenPath := filepath.Join(tmpDir, "token.json")

	if err := os.WriteFile(tokenPath, []byte("not valid json"), 0600); err != nil {
		t.Fatalf("failed to write invalid json: %v", err)
	}

	store := &TokenStore{path: tokenPath}

	_, err := store.Load()
	if err == nil {
		t.Error("Load() should return error for invalid JSON")
	}
}

func TestTokenStoreClearNonExistent(t *testing.T) {
	store := &TokenStore{
		path: "/nonexistent/path/token.json",
	}

	err := store.Clear()
	if err == nil {
		t.Error("Clear() should return error for non-existent file")
	}
}

func TestNewTokenStore(t *testing.T) {
	store := NewTokenStore()

	if store == nil {
		t.Fatal("NewTokenStore() returned nil")
	}

	if store.path == "" {
		t.Error("NewTokenStore() path should not be empty")
	}

	expectedSuffix := filepath.Join(".config", "spottui", "token.json")
	if !containsSuffix(store.path, expectedSuffix) {
		t.Errorf("NewTokenStore() path = %v, should end with %v", store.path, expectedSuffix)
	}
}

func containsSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}
