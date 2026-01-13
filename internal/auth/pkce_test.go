package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"testing"
)

func TestGenerateCodeVerifier(t *testing.T) {
	verifier, err := generateCodeVerifier()
	if err != nil {
		t.Fatalf("generateCodeVerifier() error = %v", err)
	}

	if len(verifier) < 43 || len(verifier) > 128 {
		t.Errorf("verifier length = %d, want between 43-128", len(verifier))
	}

	verifier2, _ := generateCodeVerifier()
	if verifier == verifier2 {
		t.Error("generateCodeVerifier() should produce unique values")
	}
}

func TestGenerateCodeChallenge(t *testing.T) {
	verifier := "test-verifier-string"
	challenge := generateCodeChallenge(verifier)

	expected := sha256.Sum256([]byte(verifier))
	expectedChallenge := base64.RawURLEncoding.EncodeToString(expected[:])

	if challenge != expectedChallenge {
		t.Errorf("generateCodeChallenge() = %v, want %v", challenge, expectedChallenge)
	}
}

func TestRandomString(t *testing.T) {
	length := 16
	s := randomString(length)

	if len(s) != length {
		t.Errorf("randomString(%d) length = %d, want %d", length, len(s), length)
	}

	s2 := randomString(length)
	if s == s2 {
		t.Error("randomString() should produce unique values")
	}

	for _, c := range s {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
			t.Errorf("randomString() contains invalid character: %c", c)
		}
	}
}
