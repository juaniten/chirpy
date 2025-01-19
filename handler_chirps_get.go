package main

import (
	"log"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/juaniten/chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, req *http.Request) {
	authorId := req.URL.Query().Get("author_id")
	sortParameter := req.URL.Query().Get("sort")

	var chirps []database.Chirp
	var err error

	if authorId == "" {
		chirps, err = cfg.db.GetChirps(req.Context())
	} else {
		var userId uuid.UUID
		userId, err = uuid.Parse(authorId)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
			return
		}

		chirps, err = cfg.db.GetChirpsByUser(req.Context(), userId)
	}

	if err != nil {
		log.Printf("Error getting chirps from database: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error fetching chirps", err)
		return
	}

	createdChirps := make([]Chirp, len(chirps))
	for i, chirp := range chirps {
		createdChirps[i] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	if sortParameter == "desc" {
		sort.Slice(createdChirps, func(i, j int) bool {
			return createdChirps[i].CreatedAt.After(createdChirps[j].CreatedAt)
		})
	}
	respondWithJSON(w, http.StatusOK, createdChirps)
}

func (cfg *apiConfig) handlerChirpGet(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("chirpID")
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}
	chirp, err := cfg.db.GetChirp(req.Context(), parsedUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}
	createdChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, createdChirp)
}
