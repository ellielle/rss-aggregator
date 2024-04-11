package main

import (
	"log"
	"net/http"
	"os"
	"time"

	dotenv "github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/ellielle/rss-aggregator/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// interval to retrieve RSS Feed data
	// interval is multiplied by time.Second
	const FEED_INTERVAL = 60

	// Retrive port from environment
	err := dotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	port := os.Getenv("PORT")

	// Get database connection
	db := getDatabase()

	// Shove database queries into a config struct
	dbQueries := database.New(db)

	// Fill config struct
	apiCfg := apiConfig{DB: dbQueries}

	// Create a new request mux
	mux := http.NewServeMux()

	// Attach endpoints to request mux
	apiCfg.createRouter(mux)

	// Add CORS headers
	corsMux := middlewareCors(mux)

	// Start Feed refreshing worker
	// Run in a goroutine so it can be concurrently processed
	// while handling http requests
	go updateFeedData(&apiCfg, FEED_INTERVAL*time.Second)

	// Configure server and start
	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	log.Printf("Serving on port: %v", port)
	log.Fatal(server.ListenAndServe())
}
