package main

import (
	"log"
	"net/http"
	"os"

	dotenv "github.com/joho/godotenv"
)

func main() {
	err := dotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	port := os.Getenv("PORT")
	mux := http.NewServeMux()

}
