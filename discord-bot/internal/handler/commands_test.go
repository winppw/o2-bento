package handler_test

import (
	_ "time/tzdata"
	"strings"
	"testing"
	"time"

	"github.com/o2ai/launch-assistant/discord-bot/internal/handler"
)

// Jan 2024: Mon=8, Tue=9, Wed=10, Thu=11, Fri=12, Sat=13, Sun=14
func bkk(year int, month time.Month, day, hour, min int) time.Time {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return time.Date(year, month, day, hour, min, 0, 0, loc)
}

func TestBuildStatusMessage(t *testing.T) {
	const url = "https://forms.gle/test"

	tests := []struct {
		name        string
		now         time.Time
		wantContain []string
		wantAbsent  []string
	}{
		{
			name:        "open on Tuesday",
			now:         bkk(2024, time.January, 9, 13, 0),
			wantContain: []string{"Open", "Lunch Order Status", url},
			wantAbsent:  []string{"Monday"},
		},
		{
			name:        "closed before open on Monday",
			now:         bkk(2024, time.January, 8, 9, 0),
			wantContain: []string{"Closed", "opens in"},
			wantAbsent:  []string{"Monday"},
		},
		{
			name:        "closed after deadline on Wednesday",
			now:         bkk(2024, time.January, 10, 17, 0),
			wantContain: []string{"Closed", "see you tomorrow"},
		},
		{
			name:        "Friday open — Monday note present",
			now:         bkk(2024, time.January, 12, 14, 0),
			wantContain: []string{"Open", "Monday"},
		},
		{
			name:        "Friday before open — no Monday note",
			now:         bkk(2024, time.January, 12, 8, 0),
			wantContain: []string{"Closed"},
			wantAbsent:  []string{"Monday"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msg := handler.BuildStatusMessage(tc.now, url)
			for _, want := range tc.wantContain {
				if !strings.Contains(msg, want) {
					t.Errorf("message missing %q\nfull message: %s", want, msg)
				}
			}
			for _, absent := range tc.wantAbsent {
				if strings.Contains(msg, absent) {
					t.Errorf("message should not contain %q\nfull message: %s", absent, msg)
				}
			}
		})
	}
}
