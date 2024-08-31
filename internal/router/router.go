package router

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/zimmah/rss-aggregator/internal/database"
)

type ApiConfig struct {
	DB     *database.Queries
	apiKey string
	user   User
}

func Router(cfg ApiConfig) {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /v1/err", handleErr)

	mux.Handle("GET /v1/users", logger(cfg.auth(http.HandlerFunc(cfg.handleGetUsers))))
	mux.Handle("POST /v1/users", logger(http.HandlerFunc(cfg.handlePostUsers)))

	mux.Handle("GET /v1/feeds", logger(http.HandlerFunc(cfg.handleGetFeeds)))
	mux.Handle("POST /v1/feeds", logger(cfg.auth(http.HandlerFunc(cfg.handlePostFeeds))))

	mux.Handle("GET /v1/feed_follows", logger(cfg.auth(http.HandlerFunc(cfg.handleGetFeedFollows))))
	mux.Handle("POST /v1/feed_follows", logger(cfg.auth(http.HandlerFunc(cfg.handlePostFeedFollows))))
	mux.Handle("DELETE /v1/feed_follows/{feedFollowID}", logger(cfg.auth(http.HandlerFunc(cfg.handleDeleteFeedFollowByID))))

	fmt.Printf("Server listening on %s\n", port)
	log.Fatal(server.ListenAndServe())
}
