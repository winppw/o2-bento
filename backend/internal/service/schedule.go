package service

import "time"

const (
	openHour    = 11
	openMinute  = 0
	closeHour   = 16
	closeMinute = 0
	timezone    = "Asia/Bangkok"
)

type WindowStatus struct {
	IsOpen        bool      `json:"is_open"`
	OpenAt        time.Time `json:"open_at"`
	CloseAt       time.Time `json:"close_at"`
	IsFridayExtra bool      `json:"is_friday_extra"`
	Message       string    `json:"message"`
}

func GetWindowStatus(now time.Time) WindowStatus {
	loc, _ := time.LoadLocation(timezone)
	now = now.In(loc)

	open := time.Date(now.Year(), now.Month(), now.Day(), openHour, openMinute, 0, 0, loc)
	close := time.Date(now.Year(), now.Month(), now.Day(), closeHour, closeMinute, 0, 0, loc)

	weekday := now.Weekday()
	isWeekend := weekday == time.Saturday || weekday == time.Sunday
	isOpen := !isWeekend && !now.Before(open) && now.Before(close)
	isFriday := weekday == time.Friday

	msg := ""
	switch {
	case isWeekend || now.After(close) || now.Equal(close):
		msg = "Order window is closed for today."
	case now.Before(open):
		msg = "Order window opens at 11:00 AM."
	case isOpen && isFriday:
		msg = "Order window is open. You can also submit Monday's order today."
	case isOpen:
		msg = "Order window is open. Deadline is 4:00 PM."
	default:
		msg = "Order window is closed for today."
	}

	return WindowStatus{
		IsOpen:        isOpen,
		OpenAt:        open,
		CloseAt:       close,
		IsFridayExtra: isFriday,
		Message:       msg,
	}
}
