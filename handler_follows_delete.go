package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/ellielle/rss-aggregator/internal/database"
)

func (cfg *apiConfig) handlerFollowsDelete(w http.ResponseWriter, r *http.Request, user database.User) {
	defer r.Body.Close()
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
