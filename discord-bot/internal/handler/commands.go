package handler

import (
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

const timezone = "Asia/Bangkok"

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "order",
		Description: "Show lunch order status and form link",
	},
}

func HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
	if i.ApplicationCommandData().Name != "order" {
		return
	}

	loc, _ := time.LoadLocation(timezone)
	now := time.Now().In(loc)

	open := time.Date(now.Year(), now.Month(), now.Day(), 11, 0, 0, 0, loc)
	close := time.Date(now.Year(), now.Month(), now.Day(), 16, 0, 0, 0, loc)
	isOpen := now.After(open) && now.Before(close)

	var status, extra string
	if isOpen {
		remaining := close.Sub(now).Round(time.Minute)
		status = fmt.Sprintf("✅ **Open** — closes in %s", remaining)
	} else if now.Before(open) {
		remaining := open.Sub(now).Round(time.Minute)
		status = fmt.Sprintf("🔒 **Closed** — opens in %s", remaining)
	} else {
		status = "🔒 **Closed** — see you tomorrow!"
	}

	if now.Weekday() == time.Friday && isOpen {
		extra = "\n📅 You can also submit your **Monday** order today."
	}

	msg := fmt.Sprintf("**Lunch Order Status**\n%s%s\n\n%s", status, extra, os.Getenv("GOOGLE_FORM_URL"))

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
