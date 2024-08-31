package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/zimmah/rss-aggregator/internal/database"
	"github.com/zimmah/rss-aggregator/internal/router"
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

	router.Router(config)
}
