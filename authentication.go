package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/ellielle/rss-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

// Middleware that authenticates a request, gets the user and calls the next authed handler
func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get API key from headers
		apiKey, found := strings.CutPrefix(r.Header.Get("Authorization"), "ApiKey ")
		if !found {
			respondWithError(w, http.StatusBadRequest, "Missing authorization header")
			return
		}

		// Look user up by API key
		user, err := cfg.DB.GetUserByApiKey(context.Background(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Invalid API Key")
			return
		}
		// TODO: testing
		handler(w, r, user)
	})
}
