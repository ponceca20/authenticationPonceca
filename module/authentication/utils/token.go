package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateSecureToken creates a random, URL-safe string of a given length.
// It's suitable for generating tokens for password resets, email verifications, etc.
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
