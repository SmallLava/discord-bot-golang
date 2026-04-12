package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("warning, .env file not found")
	}

	var token string = os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("missing token in enviroment var")
	}

	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("create discord session fail: %v", err)
	}

	bot.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if interaction.Type != discordgo.InteractionApplicationCommand {
			return
		}

		// 根據指令名稱處理邏輯
		if interaction.ApplicationCommandData().Name == "ping" {
			err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Pong! 🏓",
				},
			})
			if err != nil {
				log.Printf("回應指令時發生錯誤: %v", err)
			}
		}
	})

	err = bot.Open()
	if err != nil {
		log.Fatalf("connect fail: %v", err)
	}
	defer bot.Close()

	command := &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "測試機器人是否存活並回覆 Pong!",
	}

	createdCommand, err := bot.ApplicationCommandCreate(bot.State.User.ID, "943463729163558933", command)
	if err != nil {
		log.Fatalf("cannot create slash command: %v", err)
	}
	fmt.Println("Bot online, press Ctrl+C to close")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	bot.ApplicationCommandDelete(bot.State.User.ID, "", createdCommand.ID)
}
