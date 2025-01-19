package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/juaniten/chirpy/internal/auth"
	"github.com/juaniten/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, req *http.Request) {

	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid access token", err)
		return
	}
	userId, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid access token", err)
		return
	}

	type userInput struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(req.Body)
	params := userInput{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request: %s", err)
		return
	}

	hashed, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	dbUser, err := cfg.db.UpdateUser(
		req.Context(),
		database.UpdateUserParams{
			ID:             userId,
			HashedPassword: hashed,
			Email: sql.NullString{
				String: params.Email,
				Valid:  true,
			},
		})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating user: %s", err)
	}

	responseUser := User{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email.String,
		IsChirpyRed: dbUser.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusOK, responseUser)
}
