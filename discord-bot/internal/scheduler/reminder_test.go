package scheduler_test

import (
	_ "time/tzdata"
	"strings"
	"testing"
	"time"

	"github.com/o2ai/launch-assistant/discord-bot/internal/scheduler"
)

// Jan 2024: Mon=8, Tue=9, Wed=10, Thu=11, Fri=12
func bkk(year int, month time.Month, day, hour, min int) time.Time {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return time.Date(year, month, day, hour, min, 0, 0, loc)
}

func TestBuildOpenMessage_Weekday(t *testing.T) {
	const url = "https://forms.gle/test"
	now := bkk(2024, time.January, 9, 11, 0) // Tuesday
	msg := scheduler.BuildOpenMessage(now, url)

	if !strings.Contains(msg, "open") {
		t.Errorf("expected 'open' in message, got: %s", msg)
	}
	if !strings.Contains(msg, url) {
		t.Errorf("expected form URL in message, got: %s", msg)
	}
	if strings.Contains(msg, "Monday") {
		t.Errorf("Tuesday message should not mention Monday, got: %s", msg)
	}
}

func TestBuildOpenMessage_Friday(t *testing.T) {
	const url = "https://forms.gle/test"
	now := bkk(2024, time.January, 12, 11, 0) // Friday
	msg := scheduler.BuildOpenMessage(now, url)

	if !strings.Contains(msg, "Monday") {
		t.Errorf("Friday message must mention Monday, got: %s", msg)
	}
}
