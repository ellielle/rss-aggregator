package main

import (
	"context"
	"net/http"

	"github.com/ellielle/rss-aggregator/internal/database"
)

func (cfg *apiConfig) handlerFeeds(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Feeds []database.Feed `json:"feeds"`
	}

	feeds, err := cfg.DB.GetAllFeeds(context.Background())
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No feeds found")
		return
	}

	respondWithJSON(w, http.StatusOK, response{Feeds: feeds})
}
