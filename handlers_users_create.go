package main

import "net/http"

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type response struct {
		Name string `json:"name"`
	}

	cfg.DB.CreateUser(ctx context.Context, arg database.CreateUserParams)

	respondWithJSON(w, http.StatusCreated, response{Name: "dicks"})
}
