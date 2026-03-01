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
		if params.ExpiresInSecs <= 0 || params.ExpiresInSecs > 3600 {
			params.ExpiresInSecs = 3600 // default to 1 hour if not provided or exceeds max value
		}
		token, err := auth.MakeJWT(user.ID, jwt_secret, time.Duration(params.ExpiresInSecs)*time.Second)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to create JWT")
			return
		}
		jsonResponse := struct {
			User  
			Token string `json:"token"`
		}{
			User: User{
				ID:        user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				Email:     user.Email,
			},
			Token: token,
		}
		respondWithJSON(w, http.StatusOK, jsonResponse)
	}
}
