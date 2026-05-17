package service_test

import (
	_ "time/tzdata"
	"testing"
	"time"

	"github.com/o2ai/launch-assistant/backend/internal/service"
)

// bkk constructs a time in Asia/Bangkok. Use real calendar dates.
// Jan 2024: Mon=8, Tue=9, Wed=10, Thu=11, Fri=12, Sat=13, Sun=14
func bkk(year int, month time.Month, day, hour, min int) time.Time {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return time.Date(year, month, day, hour, min, 0, 0, loc)
}

func TestGetWindowStatus(t *testing.T) {
	tests := []struct {
		name          string
		now           time.Time
		wantOpen      bool
		wantFriday    bool
		wantMsgSubstr string
	}{
		{
			name:          "before open on Tuesday",
			now:           bkk(2024, time.January, 9, 10, 59), // Tue
			wantOpen:      false,
			wantFriday:    false,
			wantMsgSubstr: "opens at 11:00 AM",
		},
		{
			name:          "exactly at open time",
			now:           bkk(2024, time.January, 9, 11, 0), // Tue
			wantOpen:      true,
			wantFriday:    false,
			wantMsgSubstr: "Deadline is 4:00 PM",
		},
		{
			name:          "mid-window on Wednesday",
			now:           bkk(2024, time.January, 10, 14, 30), // Wed
			wantOpen:      true,
			wantFriday:    false,
			wantMsgSubstr: "Deadline is 4:00 PM",
		},
		{
			name:          "exactly at deadline — closed",
			now:           bkk(2024, time.January, 11, 16, 0), // Thu
			wantOpen:      false,
			wantFriday:    false,
			wantMsgSubstr: "closed for today",
		},
		{
			name:          "after deadline on Monday",
			now:           bkk(2024, time.January, 8, 17, 0), // Mon
			wantOpen:      false,
			wantFriday:    false,
			wantMsgSubstr: "closed for today",
		},
		{
			name:          "Friday open — Friday extra flag set",
			now:           bkk(2024, time.January, 12, 13, 0), // Fri
			wantOpen:      true,
			wantFriday:    true,
			wantMsgSubstr: "Monday",
		},
		{
			name:          "Friday before open — not open, still Friday",
			now:           bkk(2024, time.January, 12, 9, 0), // Fri
			wantOpen:      false,
			wantFriday:    true,
			wantMsgSubstr: "opens at 11:00 AM",
		},
		{
			name:          "Saturday during would-be window hours",
			now:           bkk(2024, time.January, 13, 14, 0), // Sat
			wantOpen:      false,
			wantFriday:    false,
			wantMsgSubstr: "closed for today",
		},
		{
			name:          "Sunday during would-be window hours",
			now:           bkk(2024, time.January, 14, 12, 0), // Sun
			wantOpen:      false,
			wantFriday:    false,
			wantMsgSubstr: "closed for today",
		},
		{
			name:          "one second before deadline — still open",
			now:           bkk(2024, time.January, 9, 15, 59), // Tue
			wantOpen:      true,
			wantFriday:    false,
			wantMsgSubstr: "Deadline is 4:00 PM",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := service.GetWindowStatus(tc.now)

			if got.IsOpen != tc.wantOpen {
				t.Errorf("IsOpen = %v, want %v", got.IsOpen, tc.wantOpen)
			}
			if got.IsFridayExtra != tc.wantFriday {
				t.Errorf("IsFridayExtra = %v, want %v", got.IsFridayExtra, tc.wantFriday)
			}
			if tc.wantMsgSubstr != "" && !containsStr(got.Message, tc.wantMsgSubstr) {
				t.Errorf("Message = %q, want substring %q", got.Message, tc.wantMsgSubstr)
			}
			if got.OpenAt.Hour() != 11 || got.OpenAt.Minute() != 0 {
				t.Errorf("OpenAt = %v, want 11:00", got.OpenAt)
			}
			if got.CloseAt.Hour() != 16 || got.CloseAt.Minute() != 0 {
				t.Errorf("CloseAt = %v, want 16:00", got.CloseAt)
			}
		})
	}
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
