package main

import (
	"encoding/json"
	"net/http"

	"github.com/h4r5h1l/Chirpy/internal/database"
)

// handlerUsers is an HTTP handler that creates a new user in the database based on the email provided in the request body
func handlerUsers(db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type UserRequest struct {
			Email string `json:"email"`
		}
		var userRequest UserRequest
		err := json.NewDecoder(r.Body).Decode(&userRequest)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		user, err := db.CreateUser(r.Context(), userRequest.Email)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to create user")
			return
		}
		respondWithJSON(w, http.StatusOK, user)
	}
}
