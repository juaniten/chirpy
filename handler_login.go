package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/juaniten/chirpy/internal/auth"
	"github.com/juaniten/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type userInput struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(req.Body)
	params := userInput{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error decoding request: %s", err))
		return
	}

	dbUser, err := cfg.db.Login(
		req.Context(), sql.NullString{String: params.Email, Valid: true})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error login user: %s", err))
		return
	}

	err = auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	expirationTime := 3600
	token, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, time.Duration(expirationTime)*time.Second)
	if err != nil {
		respondWithError(w, 500, "Unable to create JWT")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 500, "Unable to create refresh token")
		return
	}
	_, err = cfg.db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken, UserID: dbUser.ID,
	})
	if err != nil {
		fmt.Println("error creating refresh token")
		respondWithError(w, http.StatusInternalServerError, "unable to create refesh token")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email.String,
		},
		Token:        token,
		RefreshToken: refreshToken,
	})
}
