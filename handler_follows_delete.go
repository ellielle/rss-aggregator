package main

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/ellielle/rss-aggregator/internal/database"
)

// Handler to delete a User's Followed Feed by feeds_follows ID
// Requires API key
func (cfg *apiConfig) handlerFollowsDelete(w http.ResponseWriter, r *http.Request, user database.User) {
	defer r.Body.Close()
	type response struct {
		Message string `json:"message"`
	}
	// Get feed follow ID from query params
	followIdStr := r.PathValue("feed_follow_id")
	if followIdStr == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid follow ID")
		return
	}
	followId, err := uuid.Parse(followIdStr)

	// Delete the follow relationship from the user
	err = cfg.DB.DeleteFeedsFollows(context.Background(), database.DeleteFeedsFollowsParams{ID: followId, UserID: user.ID})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Feed does not exist")
		return
	}

	respondWithJSON(w, http.StatusOK, response{Message: "OK"})
}
