package main

import (
	"log"
	"os"

	"github.com/SmallLava/discord-leaving-bot/internal/cli"
	"github.com/SmallLava/discord-leaving-bot/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables if .env exists
	_ = godotenv.Load()

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/bot.db"
	}

	_, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	cli.Execute()
}
