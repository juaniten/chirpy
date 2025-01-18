package auth

import (
	"errors"
	"net/http"
	"strings"
	"unicode"
)

func GetBearerToken(headers http.Header) (string, error) {
	const bearerPrefix = "Bearer "
	authorization := headers.Get("Authorization")
	if authorization == "" {
		return "", errors.New("authorization header not present")
	}
	if !strings.HasPrefix(authorization, bearerPrefix) {
		return "", errors.New("authorization header must begin with `Bearer `")
	}

	token := strings.TrimPrefix(authorization, bearerPrefix)
	if len(token) > 0 && unicode.IsSpace(rune(token[0])) {
		return "", errors.New("there should not be extra spaces between the `Bearer` keyword and the token")
	}
	if len(token) == 0 {
		return "", errors.New("token must not be empty")
	}

	return token, nil
}
