package main

import (
	"context"
	"net/http"

	"github.com/ellielle/rss-aggregator/internal/database"
)

func (cfg *apiConfig) handlerFollowsAll(w http.ResponseWriter, r *http.Request, user database.User) {
	defer r.Body.Close()
	type response struct {
		Follows []database.FeedsFollow `json:"feeds"`
	}

	// Get the User's list of Followed Feeds
	followed, err := cfg.DB.ListFeedsFollows(context.Background(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Feed not found")
		return
	}

	respondWithJSON(w, http.StatusOK, response{Follows: followed})
}
