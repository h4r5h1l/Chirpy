package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

// readyCheckHandler is a simple HTTP handler that responds with "OK" to indicate that the server is ready to handle requests
func readyCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type header to indicate that the response is plain text
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// Write the HTTP status code and response body to indicate that the server is healthy
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// respondWithError is a helper function that sends a JSON response with an error message and the specified HTTP status code
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	errorData, err := json.Marshal(map[string]string{"error": message})
	if err != nil {
		w.Write([]byte(`{"error": "Failed to marshal error JSON"}`))
		return
	}
	w.WriteHeader(statusCode)
	w.Write(errorData)
}

// respondWithJSON is a helper function that sends a JSON response with the specified payload and HTTP status code
func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to marshal JSON")
		return
	}
	w.WriteHeader(statusCode)
	w.Write(jsonData)
}

func replace_words(input string) string {
	bad := []string{"kerfuffle", "sharbert", "fornax"}
	parts := strings.Split(input, " ")
	for i, word := range parts {
		for _, bad_word := range bad {
			if strings.ToLower(word) == bad_word {
				parts[i] = "****"
			}
		}
	}
	return strings.Join(parts, " ")
}

// handlerValidateChirp is an HTTP handler that validates the length of a chirp and responds with a JSON object indicating whether it is valid or not
func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type ChirpRequest struct {
		Body string `json:"body"`
	}
	var chirpRequest ChirpRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&chirpRequest); err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}
	if len(chirpRequest.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	chirpRequest.Body = replace_words(chirpRequest.Body)
	respondWithJSON(w, http.StatusOK, map[string]string{"cleaned_body": chirpRequest.Body})
}
