package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/juaniten/chirpy/internal/auth"
	"github.com/juaniten/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	type userInput struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(req.Body)
	params := userInput{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding request: %s", err))
		return
	}

	hashed, _ := auth.HashPassword(params.Password)
	dbUser, err := cfg.db.CreateUser(
		req.Context(),
		database.CreateUserParams{
			HashedPassword: hashed,
			Email: sql.NullString{
				String: params.Email,
				Valid:  true,
			},
		})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating user: %s", err))
	}

	responseUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email.String,
	}
	err = respondWithJSON(w, http.StatusCreated, responseUser)
	if err != nil {
		log.Println("Error generating user creation response.")
	}
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, req *http.Request) {

	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		fmt.Printf("error getting refresh token from header: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "invalid access token")
		return
	}
	userId, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid access token")
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
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding request: %s", err))
		return
	}

	hashed, _ := auth.HashPassword(params.Password)
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
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error updating user: %s", err))
	}

	responseUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email.String,
	}
	err = respondWithJSON(w, http.StatusOK, responseUser)
	if err != nil {
		log.Println("Error generating user creation response.")
	}
}
