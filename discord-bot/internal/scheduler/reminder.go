package scheduler

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
)

const timezone = "Asia/Bangkok"

func Start(dg *discordgo.Session) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Fatalf("load timezone: %v", err)
	}

	c := cron.New(cron.WithLocation(loc))

	// 11:00 AM — window opens
	c.AddFunc("0 11 * * 1-5", func() {
		now := time.Now().In(loc)
		msg := fmt.Sprintf(
			"🍱 **Lunch order window is now open!**\nSubmit your order before **4:00 PM**.\n\n%s",
			formURL(),
		)
		// Friday: add Monday note
		if now.Weekday() == time.Friday {
			msg += "\n\n📅 Today you can also submit your **Monday** order."
		}
		sendMessage(dg, msg)
	})

	// 3:30 PM — 30-minute warning
	c.AddFunc("30 15 * * 1-5", func() {
		sendMessage(dg, "⏰ **30 minutes left** to submit your lunch order!\n\n"+formURL())
	})

	// 4:00 PM — window closes
	c.AddFunc("0 16 * * 1-5", func() {
		sendMessage(dg, "🔒 Lunch order window is now **closed** for today. See you tomorrow!")
	})

	c.Start()
	log.Println("reminder scheduler started (Asia/Bangkok)")
}

func sendMessage(dg *discordgo.Session, content string) {
	channelID := os.Getenv("DISCORD_CHANNEL_ID")
	if channelID == "" {
		log.Println("DISCORD_CHANNEL_ID not set")
		return
	}
	if _, err := dg.ChannelMessageSend(channelID, content); err != nil {
		log.Printf("send discord message: %v", err)
	}
}

func formURL() string {
	return os.Getenv("GOOGLE_FORM_URL")
}
