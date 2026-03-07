package main

import (
	"net/http"

	"github.com/h4r5h1l/Chirpy/internal/auth"
	"github.com/h4r5h1l/Chirpy/internal/database"
)

func handlerRevokeToken(db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid or missing authorization header")
			return
		}
		_, err = db.GetRefreshToken(r.Context(), token)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid or expired refresh token")
			return
		}
		_, err = db.RevokeRefreshToken(r.Context(), token)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to revoke refresh token")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
