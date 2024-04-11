package main

import (
	"context"
	"net/http"
)

// Handler function to view all available feeds
func (cfg *apiConfig) handlerFeeds(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type response struct {
		Feeds []Feed `json:"feeds"`
	}

	// Get list of feeds from the database
	dbFeeds, err := cfg.DB.ListFeeds(context.Background())
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No feeds found")
		return
	}

	// Map the database API to a JSON-friendly struct
	feeds := []Feed{}
	for _, feed := range dbFeeds {
		feeds = append(feeds, DatabaseFeedToFeed(feed))
	}

	respondWithJSON(w, http.StatusOK, response{Feeds: feeds})
}
