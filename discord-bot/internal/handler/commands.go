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

// BuildStatusMessage returns the /order reply text for a given point in time.
// Extracted for testability — HandleInteraction calls this with time.Now().
func BuildStatusMessage(now time.Time, formURL string) string {
	loc, _ := time.LoadLocation(timezone)
	now = now.In(loc)

	open := time.Date(now.Year(), now.Month(), now.Day(), 11, 0, 0, 0, loc)
	close := time.Date(now.Year(), now.Month(), now.Day(), 16, 0, 0, 0, loc)
	weekday := now.Weekday()
	isWeekend := weekday == time.Saturday || weekday == time.Sunday
	isOpen := !isWeekend && !now.Before(open) && now.Before(close)

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

	return fmt.Sprintf("**Lunch Order Status**\n%s%s\n\n%s", status, extra, formURL)
}

func HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
	if i.ApplicationCommandData().Name != "order" {
		return
	}

	msg := BuildStatusMessage(time.Now(), os.Getenv("GOOGLE_FORM_URL"))

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
