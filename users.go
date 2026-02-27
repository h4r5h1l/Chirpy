package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/h4r5h1l/Chirpy/internal/database"
	"net/http"
	"time"
)

// handlerUsers is an HTTP handler that creates a new user in the database based on the email provided in the request body
func handlerUsers(db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type User struct {
			ID        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Email     string    `json:"email"`
		}
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
		jsonResp := User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		}
		respondWithJSON(w, http.StatusCreated, jsonResp)
	}
}
