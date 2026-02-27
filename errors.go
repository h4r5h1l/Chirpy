package main

import (
	"encoding/json"
	"net/http"
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
