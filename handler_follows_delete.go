package main

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/ellielle/rss-aggregator/internal/database"
)

func (cfg *apiConfig) handlerFollowsDelete(w http.ResponseWriter, r *http.Request, user database.User) {
	defer r.Body.Close()
	type response struct {
		Message string `json:"message"`
	}
	// Get feed follow ID from query params
	follow_id_str := r.PathValue("feed_follow_id")
	if follow_id_str == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid follow ID")
		return
	}
	follow_id, err := uuid.Parse(follow_id_str)

	err = cfg.DB.DeleteFeedsFollows(context.Background(), database.DeleteFeedsFollowsParams{ID: follow_id, UserID: user.ID})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Feed does not exist")
		return
	}

	respondWithJSON(w, http.StatusOK, response{Message: "OK"})
}
