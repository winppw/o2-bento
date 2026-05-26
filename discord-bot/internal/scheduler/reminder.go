package scheduler

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
)

const sendTimeout = 30 * time.Second

const timezone = "Asia/Bangkok"

// BuildOpenMessage returns the 11 AM open reminder text.
// Extracted for testability.
func BuildOpenMessage(now time.Time, url string) string {
	msg := fmt.Sprintf(
		"🍱 **Lunch order window is now open!**\nSubmit your order before **4:00 PM**.\n\n%s",
		url,
	)
	if now.Weekday() == time.Friday {
		msg += "\n\n📅 Today you can also submit your **Monday** order."
	}
	return msg
}

func Start(dg *discordgo.Session) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Fatalf("load timezone: %v", err)
	}

	c := cron.New(cron.WithLocation(loc))

	// 11:00 AM — window opens
	c.AddFunc("0 11 * * 1-5", func() {
		sendMessage(dg, BuildOpenMessage(time.Now().In(loc), formURL()))
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
	dmUserID := os.Getenv("DISCORD_DM_USER_ID")

	if channelID == "" && dmUserID == "" {
		log.Println("no notification target set (DISCORD_CHANNEL_ID or DISCORD_DM_USER_ID)")
		return
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		if channelID != "" {
			if _, err := dg.ChannelMessageSend(channelID, content); err != nil {
				log.Printf("send channel message: %v", err)
			} else {
				log.Printf("channel message sent to %s", channelID)
			}
		}

		if dmUserID != "" {
			ch, err := dg.UserChannelCreate(dmUserID)
			if err != nil {
				log.Printf("open DM channel for user %s: %v", dmUserID, err)
				return
			}
			if _, err := dg.ChannelMessageSend(ch.ID, content); err != nil {
				log.Printf("send DM to user %s: %v", dmUserID, err)
			} else {
				log.Printf("DM sent to user %s", dmUserID)
			}
		}
	}()

	select {
	case <-done:
	case <-time.After(sendTimeout):
		log.Printf("sendMessage: timed out after %s (Discord session may be reconnecting)", sendTimeout)
	}
}

func formURL() string {
	return os.Getenv("GOOGLE_FORM_URL")
}
