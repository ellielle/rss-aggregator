package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Handler to find a user via API key
func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Id        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		ApiKey    string    `json:"api_key"`
	}

	// Get API key from headers
	apiKey, found := strings.CutPrefix(r.Header.Get("Authorization"), "ApiKey ")
	if !found {
		respondWithError(w, http.StatusBadRequest, "Missing authorization header")
		return
	}

	// Look user up by API key
	user, err := cfg.DB.GetUserByApiKey(context.Background(), apiKey)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No user found")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Id:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		ApiKey:    apiKey,
	})
}
