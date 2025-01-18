package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func MakeRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("unable to create random data: %w", err)
	}

	encodedString := hex.EncodeToString(b)
	return encodedString, nil
}
