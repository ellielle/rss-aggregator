package main

import (
	"net/http"

	"github.com/ellielle/rss-aggregator/internal/database"
)

func (cfg *apiConfig) handlerFeedsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithError(w, http.StatusNotImplemented, "Not implemented")
}
