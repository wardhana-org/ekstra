package auth

import (
	"crypto/rand"
)

func randomBytes(length uint32) ([]byte, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}

	return bytes, nil
}
