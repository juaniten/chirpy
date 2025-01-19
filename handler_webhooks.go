package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/juaniten/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, req *http.Request) {
	requestApiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "api key not present", err)
		return
	}
	if requestApiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "invalid api key", err)
		return
	}

	type eventInput struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(req.Body)
	params := eventInput{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request: %s", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID format", err)
		return
	}

	_, err = cfg.db.UpgradeUser(req.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error updating user: %s", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
