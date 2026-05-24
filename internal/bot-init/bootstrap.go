package botinit

import (
	"log"
	"os"

	"github.com/SmallLava/discord-leaving-bot/internal/database"
	"github.com/bwmarrin/discordgo"
)

func BotInit() (*discordgo.Session, func()) {
	var token string = os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("missing token")
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/bot.db"
	}

	_, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("database init fail: %v", err)
	}

	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("create discord session fail: %v", err)
	}
	cleanupSlash := slashCommand(bot)

	cleanup := func() {
		cleanupSlash()
		// GORM DB doesn't strictly need a manual Close() for SQLite in most cases, 
		// but if we use a pool, we would close it here.
	}

	return bot, cleanup
}
