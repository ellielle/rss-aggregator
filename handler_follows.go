package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ellielle/rss-aggregator/internal/database"
)

func (cfg *apiConfig) handlerFollowsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
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

	ff, err := cfg.DB.CreateFeedsFollow(context.Background(), database.CreateFeedsFollowParams{
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

func (cfg *apiConfig) handlerFollowsDelete(w http.ResponseWriter, r *http.Request, user database.User) {
	type response struct {
		UserId   uuid.UUID `json:"user_id"`
		FollowId uuid.UUID `json:"deleted_feed"`
	}
	// Get feed follow ID from query params
	follow_id_str := r.PathValue("feed_follow_id")
	if follow_id_str == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid follow ID")
		return
	}
	follow_id, err := uuid.Parse(follow_id_str)

	// Get API key from headers
	apiKey, found := strings.CutPrefix(r.Header.Get("Authorization"), "ApiKey ")
	if !found {
		respondWithError(w, http.StatusBadRequest, "Missing authorization header")
		return
	}

	// Look user up by API key
	user, err = cfg.DB.GetUserByApiKey(context.Background(), apiKey)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No user found")
		return
	}

	deleted_id, err := cfg.DB.DeleteFeedsFollow(context.Background(), follow_id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusNotImplemented, response{UserId: user.ID, FollowId: deleted_id.ID})
}
