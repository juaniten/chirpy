package auth

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	id := uuid.New()
	tokenSecret := "a token secret"
	token, _ := MakeJWT(id, tokenSecret, time.Minute)
	validatedId, err := ValidateJWT(token, tokenSecret)
	if validatedId != id {
		t.Errorf("error validating JWT: %v", err)
	}
}

func TestExpiredJWT(t *testing.T) {
	id := uuid.New()
	tokenSecret := "secret"
	// Create a token that's already expired
	token, err := MakeJWT(id, tokenSecret, -time.Minute)
	if err != nil {
		t.Fatalf("error creating JWT: %v", err)
	}

	validatedId, err := ValidateJWT(token, tokenSecret)
	if validatedId != uuid.Nil {
		t.Errorf("token should be invalid: %v", err)
	}
}

func TestInvalidSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "correct-secret"
	invalidSecret := "wrong-secret"
	expiresIn := time.Minute

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("failed to create JWT: %v", err)
	}

	_, err = ValidateJWT(token, invalidSecret)

	if err == nil {
		t.Fatal("expected error when validating with wrong secret, got none")
	}

	if !strings.Contains(err.Error(), "error parsing token") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAuthorizationNotPresent(t *testing.T) {
	req, _ := http.NewRequest("GET", "/app", nil)
	if req.Header.Get("Authorization") != "" {
		t.Fatal("New request should not have Authorization header by default")
	}
	_, err := GetBearerToken(req.Header)

	if err == nil {
		t.Fatal("GetBearerToken should return and error when Authorization header is not present.")
	}
}

func TestAuthorizationPresent(t *testing.T) {
	req, _ := http.NewRequest("GET", "/app", nil)
	if req.Header.Get("Authorization") != "" {
		t.Fatal("New request should not have Authorization header by default")
	}
	req.Header.Set("Authorization", "Bearer 1234")
	token, err := GetBearerToken(req.Header)

	if err != nil {
		t.Fatal("GetBearerToken should not return and error when Authorization header is present.")
	}
	if token != "1234" {
		t.Fatal("Token should be stripped of 'Bearer ' prefix.")
	}
}

func TestInvalidAuthorizationHeaders(t *testing.T) {
	cases := []string{"1234", "Bearer  1234", "Bearer ", "", " Bearer 1234", "bearer 1234"}

	req, _ := http.NewRequest("GET", "/app", nil)

	for _, testCase := range cases {
		req.Header.Set("Authorization", testCase)
		_, err := GetBearerToken(req.Header)
		if err == nil {
			t.Fatalf("Malformed bearer token `%s` should not be accepted", testCase)
		}

	}
}
