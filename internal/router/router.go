package router

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/zimmah/rss-aggregator/internal/database"
)

type ApiConfig struct {
	DB *database.Queries
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
	mux.Handle("POST /v1/users", logger(http.HandlerFunc(cfg.handlePostUsers)))

	fmt.Printf("Server listening on %s\n", port)
	log.Fatal(server.ListenAndServe())
}
