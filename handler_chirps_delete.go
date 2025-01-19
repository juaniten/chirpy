package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/juaniten/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, req *http.Request) {

	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "bearer token not present", err)
		return
	}

	userId, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error validating JWT on chirp creation: %v", err)
		respondWithError(w, http.StatusUnauthorized, "bearer token not valid", err)
		return
	}

	chirpId, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := cfg.db.GetChirp(req.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}
	if userId != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "User not authorized", err)
		return
	}

	err = cfg.db.DeleteChirp(req.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
