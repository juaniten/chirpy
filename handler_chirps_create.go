package main

import (
	"encoding/json"
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

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, req *http.Request) {

	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "bearer token not present", err)
		return
	}

	userId, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "bearer token not valid", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := ChirpRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request: %s", err)
		return
	}
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", err)
		return
	}

	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   replaceBadWords(params.Body),
		UserID: userId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating chirp: %s", err)
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
