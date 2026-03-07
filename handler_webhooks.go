package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/h4r5h1l/Chirpy/internal/database"
)

type WebhookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func handlerPolkaWebhooks(db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var params WebhookRequest
		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		if params.Event != "user.upgraded" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Update the user's chirpy red status in the database
		_, err = db.UpdateChirpyRedStatus(r.Context(), params.Data.UserID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondWithError(w, http.StatusNotFound, "User not found")
				return
			}
			respondWithError(w, http.StatusInternalServerError, "Failed to update user status")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
