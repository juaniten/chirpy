package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/juaniten/chirpy/internal/auth"
	"github.com/juaniten/chirpy/internal/database"
)

type ChirpRequest struct {
	Body   string    `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, req *http.Request) {

	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "bearer token not present")
		return
	}

	userId, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error validating JWT on chirp creation: %v", err)
		respondWithError(w, http.StatusUnauthorized, "bearer token not valid")
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := ChirpRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error decoding request: %s", err))
		return
	}
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   replaceBadWords(params.Body),
		UserID: userId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating chirp: %s", err))
	}
	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    userId,
	})
}

func replaceBadWords(input string) string {
	words := strings.Split(input, " ")
	cleanWords := make([]string, len(words))
	badWords := make(map[string]struct{})
	badWords["kerfuffle"] = struct{}{}
	badWords["sharbert"] = struct{}{}
	badWords["fornax"] = struct{}{}

	for i, word := range words {
		if _, exists := badWords[strings.ToLower(word)]; exists {
			cleanWords[i] = "****"
		} else {
			cleanWords[i] = word
		}
	}
	return strings.Join(cleanWords, " ")
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, req *http.Request) {

	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "bearer token not present")
		return
	}

	userId, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error validating JWT on chirp creation: %v", err)
		respondWithError(w, http.StatusUnauthorized, "bearer token not valid")
		return
	}

	chirpId, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	chirp, err := cfg.db.GetChirp(req.Context(), chirpId)
	if err != nil {
		log.Printf("Error retrieving chirp from database: %v", err)
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}
	if userId != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "User not authorized")
		return
	}

	err = cfg.db.DeleteChirp(req.Context(), chirpId)
	if err != nil {
		log.Printf("Error deleting chirp from database: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Unable to delete chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
