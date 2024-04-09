package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/ellielle/rss-aggregator/internal/database"
)

// Handler to create a new feed using the given Name and Url
// Requires API key
func (cfg *apiConfig) handlerFeedsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	defer r.Body.Close()
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	type response struct {
		Id       uuid.UUID `json:"feed"`
		FollowId uuid.UUID `json:"feed_follow"`
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
		respondWithError(w, http.StatusNotFound, "Feed already exists")
		return
	}

	id = uuid.New()
	created = time.Now()
	updated = created

	follow, err := cfg.DB.CreateFeedsFollows(context.Background(), database.CreateFeedsFollowsParams{ID: id, UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "User or Feed does not exist")
		return
	}

	respondWithJSON(w, http.StatusCreated, response{Id: feed.ID, FollowId: follow.ID})
}
