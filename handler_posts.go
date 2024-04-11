package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ellielle/rss-aggregator/internal/database"
)

// Handler function that retrieves Posts a User is subscribed to
// An optional limit can be added as a query parameter, with a
// default of 10.
func (cfg *apiConfig) handlerPosts(w http.ResponseWriter, r *http.Request, user database.User) {
	defer r.Body.Close()
	type response struct {
		Posts []Post `json:"posts"`
	}

	// Set a default post retrieval limit of 10
	var postLimit int32 = 10

	if r.PathValue("limit") != "" {
		limit, err := strconv.Atoi(r.PathValue("limit"))

		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid post limit")
			return
		}
		postLimit = int32(limit)
	}

	posts, err := cfg.DB.GetPostsByUser(context.Background(), database.GetPostsByUserParams{UserID: user.ID, Limit: postLimit})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No posts found")
		return
	}

	// Create a new slice to hold the []Post
	postSlice := []Post{}
	// Convert from database.Posts to JSON-friendly Posts
	for _, post := range posts {
		postSlice = append(postSlice, DatabasePostsToPosts(post))
	}

	respondWithJSON(w, http.StatusOK, response{postSlice})
}
