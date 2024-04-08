package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/ellielle/rss-aggregator/internal/database"
)

func (cfg *apiConfig) handlerFeedsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	defer r.Body.Close()
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	type response struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		Url       string    `json:"url"`
		UserId    uuid.UUID `json:"user_id"`
	}

	// Get name from params, or throw an error if it doesn't exist
	var params = parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad JSON format")
		return
	}

	// Give each Feed a uuid for ID
	id := uuid.New()
	created := time.Now()
	updated := created

	// Create a feed owned by the User
	feed, err := cfg.DB.CreateFeed(context.Background(), database.CreateFeedParams{ID: id, CreatedAt: created, UpdatedAt: updated, Url: params.Url, Name: params.Name, UserID: user.ID})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, response{Id: feed.ID, CreatedAt: feed.CreatedAt, UpdatedAt: feed.UpdatedAt, Name: feed.Name, Url: feed.Url, UserId: feed.UserID})
}
