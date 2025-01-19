package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/juaniten/chirpy/internal/auth"
	"github.com/juaniten/chirpy/internal/database"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, req *http.Request) {
	type userInput struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(req.Body)
	params := userInput{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request: %s", err)
		return
	}

	hashed, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

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
		respondWithError(w, http.StatusInternalServerError, "Error creating user: %s", err)
	}

	responseUser := User{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email.String,
		IsChirpyRed: dbUser.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusCreated, responseUser)
}
