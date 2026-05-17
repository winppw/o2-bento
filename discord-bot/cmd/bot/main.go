package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"github.com/o2ai/launch-assistant/discord-bot/internal/handler"
	"github.com/o2ai/launch-assistant/discord-bot/internal/scheduler"
)

func main() {
	_ = godotenv.Load()

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN is required")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("create discord session: %v", err)
	}

	dg.AddHandler(handler.HandleInteraction)

	if err = dg.Open(); err != nil {
		log.Fatalf("open discord connection: %v", err)
	}
	defer dg.Close()

	appID := dg.State.User.ID

	// Set presence — shows "Watching lunch orders 🍱" under the bot name
	if err = dg.UpdateStatusComplex(discordgo.UpdateStatusData{
		Status: "online",
		Activities: []*discordgo.Activity{
			{
				Name: "lunch orders 🍱",
				Type: discordgo.ActivityTypeWatching,
			},
		},
	}); err != nil {
		log.Printf("set presence: %v", err)
	}

	// Register slash commands
	for _, cmd := range handler.Commands {
		if _, err := dg.ApplicationCommandCreate(appID, "", cmd); err != nil {
			log.Printf("register command %s: %v", cmd.Name, err)
		}
	}

	scheduler.Start(dg)

	log.Println("o2Bento is running. Ctrl+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}
