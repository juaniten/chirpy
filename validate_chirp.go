package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	const maxChirpLength = 140
	type chirp struct {
		Body string `json:"body"`
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	decoder := json.NewDecoder(req.Body)
	params := chirp{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error decoding chirp: %s", err))
		return
	}

	if len(params.Body) > maxChirpLength {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	err = respondWithJSON(w, 200,
		struct {
			Valid bool `json:"valid"`
		}{Valid: true})
	if err != nil {
		respondWithError(w, 500, "Error coding valid response :S")
		return
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}
