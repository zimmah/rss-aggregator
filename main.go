package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/zimmah/rss-aggregator/internal/database"
	"github.com/zimmah/rss-aggregator/internal/router"
	"github.com/zimmah/rss-aggregator/internal/scraper"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbConnectionString := os.Getenv("CONNECTION_STRING")
	if dbConnectionString == "" {
		log.Fatalf("CONNECTION_STRING environment variable not set")
	}

	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		log.Fatalf("couldn't establish database connection: %v", err)
		return
	}
	defer db.Close()

	queries := database.New(db)
	config := router.ApiConfig{
		DB: queries,
	}

	const collectionConcurrency = 10
	const collectionInterval = time.Minute
	go scraper.StartWorker(queries, collectionConcurrency, collectionInterval)

	router.Router(&config)
}
