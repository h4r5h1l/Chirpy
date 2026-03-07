package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/h4r5h1l/Chirpy/internal/auth"
	"github.com/h4r5h1l/Chirpy/internal/database"
	"net/http"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	IsRed     bool      `json:"is_chirpy_red"`
}
type UserRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	AccessToken string `json:"access_token"`
}

// handlerUsers is an HTTP handler that creates a new user in the database based on the email provided in the request body
func handlerUsers(db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var params UserRequest
		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		hashedPassword, err := auth.HashPassword(params.Password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
			return
		}
		user, err := db.CreateUser(r.Context(), database.CreateUserParams{
			Email:          params.Email,
			HashedPassword: hashedPassword,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to create user: "+err.Error())
			return
		}
		jsonResp := User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			IsRed:     user.IsChirpyRed,
		}
		respondWithJSON(w, http.StatusCreated, jsonResp)
	}
}

func handlerUpdateUser(db *database.Queries, jwt_secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Missing or invalid token")
			return
		}
		userId, err := auth.ValidateJWT(tokenString, jwt_secret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}
		user, err := db.GetUserByID(r.Context(), userId)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		var params UserRequest
		err = json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		hashedPassword, err := auth.HashPassword(params.Password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
			return
		}
		_, err = db.UpdatePassword(r.Context(), database.UpdatePasswordParams{
			HashedPassword: hashedPassword,
			ID:             user.ID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to update password")
			return
		}
		user, err = db.UpdateEmail(r.Context(), database.UpdateEmailParams{
			Email: params.Email,
			ID:    user.ID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to update email")
			return
		}

		jsonResp := User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			IsRed:     user.IsChirpyRed,
		}
		respondWithJSON(w, http.StatusOK, jsonResp)
	}
}
