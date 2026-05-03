package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func randomBytes(length uint32) ([]byte, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}

	return bytes, nil
}

func GenerateRawToken() (string, error) {
	bytes, err := randomBytes(32)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func HashToken(rawToken string) string {
	hashedToken := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(hashedToken[:])
}
