package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	const maxChirpLength = 140
	type chirp struct {
		Body string `json:"body"`
	}
	type validationError struct {
		Error string `json:"error"`
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	decoder := json.NewDecoder(req.Body)
	params := chirp{}
	err := decoder.Decode(&params)
	if err != nil {
		errorMessage, _ := json.Marshal(validationError{
			Error: fmt.Sprintf("Error decoding chirp: %s", err),
		})

		w.WriteHeader(500)
		w.Write(errorMessage)
		return
	}

	if len(params.Body) > maxChirpLength {
		errorMessage, _ := json.Marshal(validationError{
			Error: fmt.Sprint("Chirp is too long"),
		})

		w.WriteHeader(400)
		w.Write(errorMessage)
		return

	}

	data, err := json.Marshal(struct {
		Valid bool `json:"valid"`
	}{
		Valid: true,
	})
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(data))
}
