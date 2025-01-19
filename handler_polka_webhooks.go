package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/juaniten/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, req *http.Request) {
	log.Printf("Headers received: %v", req.Header)
	requestApiKey, err := auth.GetAPIKey(req.Header)
	log.Println("API KEY REQUEST: ", requestApiKey)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusUnauthorized, "api key not present")
		return
	}
	if requestApiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "invalid api key")
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
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding request: %s", err))
		return
	}

	// Check if the event is "user.upgraded"
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Parse the user ID from the event data
	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	log.Printf("Received webhook for user ID: %s", params.Data.UserID)
	// Update the user in the database to mark as Chirpy Red member
	dbUser, err := cfg.db.UpgradeUser(req.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error updating user: %s", err))
		return
	}
	log.Printf("Upgrade result: %+v", dbUser)

	// Respond with 204 No Content if the update is successful
	w.WriteHeader(http.StatusNoContent)
}
