package main

import (
	"log"
	"net/http"
	"os"
	"time"

	dotenv "github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/ellielle/rss-aggregator/internal/database"
	"github.com/ellielle/rss-aggregator/internal/httpclient"
)

type apiConfig struct {
	DB     *database.Queries
	Client *httpclient.Client
}

func main() {
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

	// Create a new http Client with cache
	// The cache reap time is low so it can be tested
	client := httpclient.NewClient(5*time.Second, 20*time.Second)

	// Fill config struct
	apiCfg := apiConfig{DB: dbQueries, Client: &client}

	// Create a new request mux
	mux := http.NewServeMux()

	apiCfg.createRouter(mux)

	// Add CORS headers
	corsMux := middlewareCors(mux)

	// Configure server and start
	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	log.Printf("Serving on port: %v", port)
	log.Fatal(server.ListenAndServe())
}
