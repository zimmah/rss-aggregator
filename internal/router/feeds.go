package router

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zimmah/rss-aggregator/internal/database"
)

type Feed struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	LastFetchedAt *time.Time `json:"last_fetched_at"`
	Name          string     `json:"name"`
	URL           string     `json:"url"`
}

type FeedAndFeedFollow struct {
	Feed       Feed       `json:"feed"`
	FeedFollow FeedFollow `json:"feed_follow"`
}

func databaseFeedToFeed(dbFeed database.Feed) Feed {
	feed := Feed{
		ID:        dbFeed.ID,
		CreatedAt: dbFeed.CreatedAt,
		UpdatedAt: dbFeed.UpdatedAt,
		Name:      dbFeed.Name,
		UserID:    dbFeed.UserID,
		URL:       dbFeed.Url,
	}
	feed.SetLastFetchedAt(dbFeed.LastFetchedAt)
	return feed
}

func (f *Feed) SetLastFetchedAt(nt sql.NullTime) {
	if nt.Valid {
		f.LastFetchedAt = &nt.Time
	} else {
		f.LastFetchedAt = nil
	}
}

func (cfg *ApiConfig) handlePostFeeds(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	feed := parameters{}
	err := decoder.Decode(&feed)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	now := time.Now().UTC()
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      feed.Name,
		Url:       feed.URL,
		UserID:    cfg.user.ID,
	}
	ffParams := database.FollowFeedParams{
		ID:        uuid.New(),
		FeedID:    params.ID,
		UserID:    cfg.user.ID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	dbFeed, err := cfg.DB.CreateFeed(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed")
		return
	}
	dbFeedFollow, err := cfg.DB.FollowFeed(r.Context(), ffParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't follow feed")
		return
	}

	feedFeedFollow := FeedAndFeedFollow{
		Feed:       databaseFeedToFeed(dbFeed),
		FeedFollow: databaseFeedFollowToFeedFollow(dbFeedFollow),
	}
	respondWithJSON(w, http.StatusCreated, feedFeedFollow)
}

func (cfg *ApiConfig) handleGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get feeds")
		return
	}

	if len(feeds) == 0 {
		respondWithError(w, http.StatusNotFound, "No feeds found for this user")
		return
	}
	respondWithJSON(w, http.StatusOK, feeds)
}
