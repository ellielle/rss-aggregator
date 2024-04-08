package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	dotenv "github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/ellielle/rss-aggregator/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	const filepathRoot = "."
	err := dotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_CONNECTION")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("error connecting to database")
	}

	dbQueries := database.New(db)
	apiCfg := apiConfig{DB: dbQueries}

	mux := http.NewServeMux()

	//	handlerFileserver := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	//	mux.Handle("GET /app/*", handlerFileserver)
	// Endpoints that respond with 200OK and an error, respectively
	mux.HandleFunc("GET /v1/readiness", handlerReadiness)
	mux.HandleFunc("GET /v1/err", handlerError)
	// Endpoint for users to get their own information, requires
	// their api_key
	mux.HandleFunc("GET /v1/users", apiCfg.handlerUsers)
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

	corsMux := middlewareCors(mux)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	log.Printf("Serving on port: %v", port)
	log.Fatal(server.ListenAndServe())
}
