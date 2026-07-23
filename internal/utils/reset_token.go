package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// GenerateResetToken returns a cryptographically random, URL-safe hex
// string (32 bytes -> 64 hex chars) to hand back to the caller as the
// raw, one-time password reset token.
func GenerateResetToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// HashToken returns the SHA-256 hex digest of a raw token, for storage/lookup.
// Reset tokens are single-use, short-lived, high-entropy random values (not
// passwords), so a fast hash is appropriate — deliberately not bcrypt.
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
