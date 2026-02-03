package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"math/big"
)

// generateCodeVerifier creates a cryptographically random code verifier
// RFC 7636 requires 43-128 characters
func generateCodeVerifier() (string, error) {
	b := make([]byte, 96)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// generateCodeChallenge creates SHA256 hash of verifier for PKCE
func generateCodeChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

// randomString generates a random alphanumeric string of given length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err) // Crypto rand failure is unrecoverable
		}
		result[i] = charset[n.Int64()]
	}
	return string(result)
}
