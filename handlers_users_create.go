package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/ellielle/rss-aggregator/internal/database"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type parameters struct {
		Name string `json:"name"`
	}
	type response struct {
		Name string `json:"name"`
	}

	var params = parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad JSON format")
		return
	}
	id := uuid.New()
	created := time.Now()
	updated := created

	user, err := cfg.DB.CreateUser(context.Background(), database.CreateUserParams{ID: id, Name: params.Name, CreatedAt: created, UpdatedAt: updated})
	if err != nil {
		respondWithError(w, http.StatusMethodNotAllowed, err.Error())
		return
	}
	log.Print(user)

	respondWithJSON(w, http.StatusCreated, response{Name: user.Name})
}
