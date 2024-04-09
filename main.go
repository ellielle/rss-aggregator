package main

import (
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
	apiCfg := apiConfig{DB: dbQueries}

	// Create a new request mux
	mux := http.NewServeMux()

	apiCfg.createRouter(mux, db)

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
