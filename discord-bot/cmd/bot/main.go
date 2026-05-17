package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"github.com/o2ai/launch-assistant/discord-bot/internal/handler"
	"github.com/o2ai/launch-assistant/discord-bot/internal/scheduler"
)

type appPatch struct {
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

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

	// Set bot presence — visible under the bot's name in the member list
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

	// Update application description and tags — visible on the bot's profile card
	patch := appPatch{
		Description: "Daily lunch order reminders for the O2AI team. Submit your bento before 4 PM! 🍱 Use /order to check the window status anytime.",
		Tags:        []string{"lunch", "order", "reminder", "o2ai", "bento"},
	}
	patchBytes, _ := json.Marshal(patch)
	if _, err = dg.RequestWithBucketID("PATCH", discordgo.EndpointApplications+"@me", patchBytes, discordgo.EndpointApplications+"@me"); err != nil {
		log.Printf("update application metadata: %v", err)
	} else {
		log.Println("application description and tags updated")
	}

	// Register slash commands
	appID := dg.State.User.ID
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
