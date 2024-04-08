package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/ellielle/rss-aggregator/internal/database"
)

// User only needs to supply a name, this is a simple project
// focused on SQL use in Go
func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type parameters struct {
		Name string `json:"name"`
	}
	type response struct {
		Name   string `json:"name"`
		ApiKey string `json:"api_key"`
	}

	// Get name from params, or throw an error if it doesn't exist
	var params = parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad JSON format")
		return
	}

	// Give each user a uuid for ID
	id := uuid.New()
	created := time.Now()
	updated := created

	user, err := cfg.DB.CreateUser(context.Background(), database.CreateUserParams{ID: id, Name: params.Name, CreatedAt: created, UpdatedAt: updated})
	if err != nil {
		respondWithError(w, http.StatusMethodNotAllowed, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, response{Name: user.Name, ApiKey: user.ApiKey})
}
