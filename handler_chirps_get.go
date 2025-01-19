package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {

	authorId := req.URL.Query().Get("author_id")
	if authorId == "" {

		chirps, err := cfg.db.GetChirps(req.Context())
		if err != nil {
			log.Printf("Error getting chirps from database: %s", err)
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

		err = respondWithJSON(w, http.StatusOK, createdChirps)
		if err != nil {
			log.Printf("error creating response for getting chirps: %s", err)
			return
		}
	}

	userId, err := uuid.Parse(authorId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid author ID")
		return
	}

	chirps, err := cfg.db.GetChirpsByUser(req.Context(), userId)
	if err != nil {
		log.Printf("Error getting chirps from database: %s", err)
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

	err = respondWithJSON(w, http.StatusOK, createdChirps)
	if err != nil {
		log.Printf("error creating response for getting chirps: %s", err)
		return
	}

}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		fmt.Println("Error parsing UUID:", err)
		return
	}
	chirp, err := cfg.db.GetChirp(req.Context(), parsedUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}
	createdChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	err = respondWithJSON(w, http.StatusOK, createdChirp)
	if err != nil {
		log.Printf("error creating response for getting chirps: %s", err)
	}
}
