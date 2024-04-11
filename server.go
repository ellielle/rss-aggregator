package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	dotenv "github.com/joho/godotenv"
)

// Load the database connection URL from the environment
// Then create a database connection and return it
func getDatabase() *sql.DB {
	err := dotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	dbURL := os.Getenv("DB_CONNECTION")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("error connecting to database")
	}

	return db
}

// Create all the routes for the http mux
func (apiCfg *apiConfig) createRouter(mux *http.ServeMux) {
	// const filepathRoot = "."
	//	handlerFileserver := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	//	mux.Handle("GET /app/*", handlerFileserver)

	// Endpoints that respond with 200 OK and an error, respectively
	mux.HandleFunc("GET /v1/readiness", handlerReadiness)
	mux.HandleFunc("GET /v1/err", handlerError)
	// Authenticated endpoint for users to get their own information
	mux.HandleFunc("GET /v1/users", apiCfg.middlewareAuth(apiCfg.handlerUsers))
	// Creates a user
	mux.HandleFunc("POST /v1/users", apiCfg.handlerUsersCreate)
	// Creates a feed for an authenticated user
	mux.HandleFunc("POST /v1/feeds", apiCfg.middlewareAuth(apiCfg.handlerFeedsCreate))
	// Non-authenticated endpoint to retrieve all feeds
	mux.HandleFunc("GET /v1/feeds", apiCfg.handlerFeeds)
	// Authenticated endpoint for a User to subscribe to a Feed
	mux.HandleFunc("POST /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFollowsCreate))
	// Authenticated endpoint to unsubscribe a User from a Feed
	mux.HandleFunc("DELETE /v1/feed_follows/{feed_follow_id}", apiCfg.middlewareAuth(apiCfg.handlerFollowsDelete))
	// Authenticated endpoint for a User to view all Feeds they Follow
	mux.HandleFunc("GET /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFollowsAll))
	// Authenticated endpoint for a User to view Posts from their feeds, with an optional
	// limit query
	mux.HandleFunc("GET /v1/posts/{limit}", apiCfg.middlewareAuth(apiCfg.handlerPosts))
}
