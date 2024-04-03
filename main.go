package main

import (
	"log"
	"net/http"
	"os"

	dotenv "github.com/joho/godotenv"
)

func main() {
	const filepathRoot = "."
	err := dotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	port := os.Getenv("PORT")
	mux := http.NewServeMux()

	handlerFileserver := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/*", handlerFileserver)
	mux.HandleFunc("/v1/readiness", handlerReadiness)
	mux.HandleFunc("/v1/err", handlerError)

	corsMux := middlewareCors(mux)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	log.Printf("Serving on port: %v", port)
	log.Fatal(server.ListenAndServe())
}
