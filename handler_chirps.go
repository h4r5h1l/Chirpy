package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/h4r5h1l/Chirpy/internal/auth"
	"github.com/h4r5h1l/Chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}
type ChirpRequest struct {
	Body string `json:"body"`
}

// handlerChirps is an HTTP handler that validates the length of a chirp and responds with a JSON object indicating whether it is valid or not
func handlerChirps(db *database.Queries, jwt_secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params ChirpRequest
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			respondWithError(w, http.StatusBadRequest, "Something went wrong")
			return
		}
		if len(params.Body) > 140 {
			respondWithError(w, http.StatusBadRequest, "Chirp is too long")
			return
		}
		params.Body = replace_words(params.Body)
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Missing or invalid token")
			return
		}
		claims, err := auth.ValidateJWT(token, jwt_secret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}
		chirp, err := db.CreateChirp(r.Context(), database.CreateChirpParams{
			Body:   params.Body,
			UserID: claims,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to create chirp")
			return
		}
		jsonResp := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    claims,
		}
		respondWithJSON(w, http.StatusCreated, jsonResp)
	}
}
func handlerGetChirps(db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chirps, err := db.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to get chirps")
			return
		}
		var jsonResp []Chirp
		for _, chirp := range chirps {
			jsonResp = append(jsonResp, Chirp{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
			})
		}
		respondWithJSON(w, http.StatusOK, jsonResp)
	}
}
func handlerGetChirpById(db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("chirpID")
		if idStr == "" {
			respondWithError(w, http.StatusBadRequest, "Missing chirp ID")
			return
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
			return
		}
		chirp, err := db.GetChirpByChirpID(r.Context(), id)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Chirp not found")
			return
		}
		jsonResp := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		respondWithJSON(w, http.StatusOK, jsonResp)
	}
}

func handlerDeleteChirp(db *database.Queries, jwt_secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("chirpID")
		if idStr == "" {
			respondWithError(w, http.StatusBadRequest, "Missing chirp ID")
			return
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
			return
		}
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Missing or invalid token")
			return
		}
		claims, err := auth.ValidateJWT(token, jwt_secret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}
		chirp, err := db.GetChirpByChirpID(r.Context(), id)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Chirp not found")
			return
		}
		if chirp.UserID != claims {
			respondWithError(w, http.StatusForbidden, "You are not the owner of this chirp")
			return
		}
		err = db.DeleteChirpByChirpID(r.Context(), id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to delete chirp")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
