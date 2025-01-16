package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	const maxChirpLength = 140
	type chirp struct {
		Body string `json:"body"`
	}

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
			CleanedBody string `json:"cleaned_body"`
		}{
			CleanedBody: replaceBadWords(params.Body),
		})
	if err != nil {
		respondWithError(w, 500, "Error coding response")
		return
	}
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
