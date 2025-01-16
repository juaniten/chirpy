package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, req *http.Request) {
	type email struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(req.Body)
	params := email{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error decoding request: %s", err))
		return
	}

	dbUser, err := cfg.db.CreateUser(
		req.Context(),
		sql.NullString{
			String: params.Email,
			Valid:  true,
		})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error creating user: %s", err))
	}

	responseUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email.String,
	}
	err = respondWithJSON(w, 201, responseUser)
	if err != nil {
		log.Println("Error generating user creation response.")
	}
}
