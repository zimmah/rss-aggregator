package router

import (
	"github.com/zimmah/rss-aggregator/internal/database"
	"net/http"
)

func (cfg *ApiConfig) handleGetPosts(w http.ResponseWriter, r *http.Request) {
	params := database.GetPostsByUserParams{
		UserID: cfg.user.ID,
		Limit:  20,
	}
	posts, err := cfg.DB.GetPostsByUser(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could't retrieve posts")
		return
	}

	if len(posts) == 0 {
		respondWithError(w, http.StatusNotFound, "No posts found for this user")
		return
	}

	respondWithJSON(w, http.StatusOK, posts)
}
