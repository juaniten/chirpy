package auth

import (
	"errors"
	"fmt"
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
		return "", fmt.Errorf("authorization header must begin with `%s `", bearerPrefix)
	}

	token := strings.TrimPrefix(authorization, bearerPrefix)
	if len(token) > 0 && unicode.IsSpace(rune(token[0])) {
		return "", fmt.Errorf("there should not be extra spaces between the `%s` keyword and the token", bearerPrefix)
	}
	if len(token) == 0 {
		return "", errors.New("token must not be empty")
	}

	return token, nil
}
