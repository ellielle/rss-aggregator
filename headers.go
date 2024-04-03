package main

import (
	"encoding/json"
	"net/http"
)

// Middleware handler to add basic (and open, which is necessary for the course to access the server) CORS headers
// Allows BootDev to access via browser
func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Wrapper for Errors to be returned in JSON
func respondWithError(w http.ResponseWriter, code int, message string) error {
	type returnError struct {
		Error string `json:"error"`
	}
	return respondWithJSON(w, code, returnError{Error: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}
