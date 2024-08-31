package router

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zimmah/rss-aggregator/internal/database"
)

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	FeedID    uuid.UUID `json:"feed_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func databaseFeedFollowToFeedFollow(feedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        feedFollow.ID,
		FeedID:    feedFollow.FeedID,
		UserID:    feedFollow.UserID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
	}
}

func (cfg *ApiConfig) handlePostFeedFollows(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	feed := parameters{}
	err := decoder.Decode(&feed)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	now := time.Now().UTC()
	params := database.FollowFeedParams{
		ID:        uuid.New(),
		FeedID:    feed.FeedID,
		UserID:    cfg.user.ID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	dbFeedFollow, err := cfg.DB.FollowFeed(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't follow feed")
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedFollowToFeedFollow(dbFeedFollow))
}

func (cfg *ApiConfig) handleDeleteFeedFollowByID(w http.ResponseWriter, r *http.Request) {
	feedFollowID, err := uuid.Parse(r.PathValue("feedFollowID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse feedFollowID")
		return
	}

	err = cfg.DB.DeleteFeedFollowByFeedFollowID(r.Context(), feedFollowID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete feedFollow")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *ApiConfig) handleGetFeedFollows(w http.ResponseWriter, r *http.Request) {
	dbFeedFollows, err := cfg.DB.GetFeedFollowsByUserApiKey(r.Context(), cfg.apiKey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get feedFollows")
		return
	}

	var feedFollows []FeedFollow
	for _, dbFeedFollow := range dbFeedFollows {
		feedFollows = append(feedFollows, databaseFeedFollowToFeedFollow(dbFeedFollow))
	}

	if len(feedFollows) == 0 {
		respondWithError(w, http.StatusNotFound, "No feedFollows found for user")
		return
	}
	respondWithJSON(w, http.StatusOK, feedFollows)
}
