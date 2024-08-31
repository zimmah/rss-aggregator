package router

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zimmah/rss-aggregator/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Apikey    string    `json:"api_key"`
}

func databaseUserToUser(user database.User) User {
	return User{
		ID:        user.ID,
		Apikey:    user.Apikey,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
	}
}

func (cfg *ApiConfig) handlePostUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	user := parameters{}
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	now := time.Now().UTC()
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      user.Name,
	}
	dbUser, err := cfg.DB.CreateUser(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseUserToUser(dbUser))
}

func (cfg *ApiConfig) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, cfg.user)
}
