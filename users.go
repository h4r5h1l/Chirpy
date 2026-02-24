package main

import (
	"encoding/json"
	"net/http"

	"github.com/h4r5h1l/Chirpy/internal/database"
)

type User struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Email     string `json:"email"`
}

// handlerUsers is an HTTP handler that creates a new user in the database based on the email provided in the request body
func handlerUsers(db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type UserRequest struct {
			Email string `json:"email"`
		}
		var params UserRequest
		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		user, err := db.CreateUser(r.Context(), params.Email)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to create user")
			return
		}
		u := User{
			ID:        user.ID.String(),
			CreatedAt: user.CreatedAt.String(),
			UpdatedAt: user.UpdatedAt.String(),
			Email:     user.Email,
		}
		respondWithJSON(w, http.StatusCreated, u)
	}
}
