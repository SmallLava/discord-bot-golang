package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	botinit "github.com/SmallLava/discord-leaving-bot/internal/bot-init"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	_ = godotenv.Load()

	bot, cleanup := botinit.BotInit()

	err := bot.Open()
	if err != nil {
		log.Fatalf("connect fail: %v", err)
	}
	defer bot.Close()

	fmt.Println("Bot online, press Ctrl+C to close")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	cleanup()
}
