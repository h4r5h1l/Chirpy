package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/h4r5h1l/Chirpy/internal/database"
	"net/http"
	"time"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}
type ChirpRequest struct {
	Body   string `json:"body"`
	UserID string `json:"user_id"`
}

// handlerChirps is an HTTP handler that validates the length of a chirp and responds with a JSON object indicating whether it is valid or not
func handlerChirps(db *database.Queries) http.HandlerFunc {
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
		uid, err := uuid.Parse(params.UserID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}
		chirp, err := db.CreateChirp(r.Context(), database.CreateChirpParams{
			Body:   params.Body,
			UserID: uid,
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
			UserID:    chirp.UserID,
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
