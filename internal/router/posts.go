package router

import (
	"net/http"
	"strconv"

	"github.com/zimmah/rss-aggregator/internal/database"
)

func (cfg *ApiConfig) handleGetPosts(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing query params")
	}
	if limit == 0 {
		limit = 20
	}

	params := database.GetPostsByUserParams{
		UserID: cfg.user.ID,
		Limit:  int32(limit),
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
