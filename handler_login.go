package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/h4r5h1l/Chirpy/internal/auth"
	"github.com/h4r5h1l/Chirpy/internal/database"
)

func handlerLogin(db *database.Queries, jwt_secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params UserRequest
		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		user, err := db.GetUserByEmail(r.Context(), params.Email)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}

		match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to check password hash")
			return
		}
		if !match {
			respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
			return
		}
		access_token, err := auth.MakeJWT(user.ID, jwt_secret, time.Hour)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to create JWT")
			return
		}
		refresh_token, err := auth.MakeRefreshToken()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to create refresh token")
			return
		}
		err = db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			Token:     refresh_token,
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(60 * 24 * time.Hour), // refresh token valid for 60 days
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to save refresh token")
			return
		}
		jsonResponse := struct {
			User
			Token        string `json:"token"`
			RefreshToken string `json:"refresh_token"`
		}{
			User: User{
				ID:        user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				Email:     user.Email,
				IsRed:     user.IsChirpyRed,
			},
			Token:        access_token,
			RefreshToken: refresh_token,
		}
		respondWithJSON(w, http.StatusOK, jsonResponse)
	}
}
