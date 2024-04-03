package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type response struct {
		Status string `json:"status"`
	}
	respondWithJSON(w, http.StatusOK, response{Status: "ok"})
}

func handlerError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
