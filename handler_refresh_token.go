package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/juaniten/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, req *http.Request) {
	requestToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid refresh token", err)
		return
	}

	refreshToken, err := cfg.db.GetRefreshToken(req.Context(), requestToken)
	if err != nil {
		fmt.Printf("error getting refresh token from database: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "error processing refresh token", err)
		return
	}
	if !refreshToken.ExpiresAt.Valid || time.Now().After(refreshToken.ExpiresAt.Time) || refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "invalid refresh token", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(req.Context(), refreshToken.Token)
	if err != nil {
		fmt.Printf("error getting user from refresh token from database: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "unable to process user for the refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to process the refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, req *http.Request) {
	requestToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid refresh token", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(req.Context(), requestToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to revoke refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
