package main

import (
	"context"
	"net/http"
)

func (cfg *apiConfig) handlerFeeds(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type response struct {
		Feeds []Feed `json:"feeds"`
	}

	db_feeds, err := cfg.DB.ListFeeds(context.Background())
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No feeds found")
		return
	}

	feeds := []Feed{}

	for _, feed := range db_feeds {
		feeds = append(feeds, DatabaseFeedToFeed(feed))
	}

	respondWithJSON(w, http.StatusOK, response{Feeds: feeds})
}
