package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode"
)

func GetAPIKey(headers http.Header) (string, error) {
	const authorizationHeaderPrefix = "ApiKey "
	authorization := headers.Get("Authorization")
	if authorization == "" {
		return "", errors.New("authorization header not present")
	}
	if !strings.HasPrefix(authorization, authorizationHeaderPrefix) {
		return "", fmt.Errorf("authorization header must begin with `%s `", authorizationHeaderPrefix)
	}

	key := strings.TrimPrefix(authorization, authorizationHeaderPrefix)
	if len(key) > 0 && unicode.IsSpace(rune(key[0])) {
		return "", fmt.Errorf("there should not be extra spaces between the `%s` keyword and the token", authorizationHeaderPrefix)
	}
	if len(key) == 0 {
		return "", errors.New("token must not be empty")
	}

	return key, nil
}
