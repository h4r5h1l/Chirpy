package main

import (
	"net/http"
	"time"

	"github.com/h4r5h1l/Chirpy/internal/auth"
	"github.com/h4r5h1l/Chirpy/internal/database"
)

func handlerRefreshToken(db *database.Queries, jwt_secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid or missing authorization header")
			return
		}
		tokenData, err := db.GetRefreshToken(r.Context(), token)
		if err != nil || tokenData.ExpiresAt.Before(time.Now()) || tokenData.RevokedAt.Valid {
			respondWithError(w, http.StatusUnauthorized, "Invalid or expired refresh token")
			return
		}
		user_id := tokenData.UserID
		access_token, err := auth.MakeJWT(user_id, jwt_secret, time.Hour)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to create JWT")
			return
		}
		jsonResponse := struct {
			Token string `json:"token"`
		}{
			Token: access_token,
		}
		respondWithJSON(w, http.StatusOK, jsonResponse)
	}
}
