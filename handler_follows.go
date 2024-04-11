package main

import (
	"context"
	"net/http"

	"github.com/ellielle/rss-aggregator/internal/database"
)

// Handler to retrieve all feeds followed by a User
// Requires API key
func (cfg *apiConfig) handlerFollowsAll(w http.ResponseWriter, r *http.Request, user database.User) {
	defer r.Body.Close()
	type response struct {
		Follows []FeedsFollow `json:"feeds"`
	}

	// Get the User's list of Followed Feeds
	followed, err := cfg.DB.ListFeedsFollows(context.Background(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Feed not found")
		return
	}

	// Map database struct to a JSON friend struct
	feedsFollow := []FeedsFollow{}
	for _, follow := range followed {
		feedsFollow = append(feedsFollow, DatabaseFeedFollowsToFeedFollows(follow))
	}

	respondWithJSON(w, http.StatusOK, response{Follows: feedsFollow})
}
