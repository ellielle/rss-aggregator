package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/ellielle/rss-aggregator/internal/database"
)

func (cfg *apiConfig) handlerFollowsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	defer r.Body.Close()
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	type response struct {
		Id        uuid.UUID `json:"id"`
		FeedId    uuid.UUID `json:"feed_id"`
		UserId    uuid.UUID `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	// Get feed ID from params, return an error if it doesn't exist
	var params = parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad feed ID")
		return
	}

	// Give each Feed Follow a uuid for ID
	id := uuid.New()
	created := time.Now()
	updated := created

	ff, err := cfg.DB.CreateFeedsFollows(context.Background(), database.CreateFeedsFollowsParams{
		ID:        id,
		CreatedAt: created,
		UpdatedAt: updated,
		FeedID:    params.FeedId,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error saving feed")
		return
	}

	respondWithJSON(w, http.StatusCreated, response{Id: ff.ID, CreatedAt: ff.CreatedAt, UpdatedAt: ff.UpdatedAt, FeedId: ff.FeedID, UserId: ff.UserID})
}
