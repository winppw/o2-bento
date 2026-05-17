package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/o2ai/launch-assistant/backend/internal/service"
)

type StatusResponse struct {
	service.WindowStatus
	FormURL string `json:"form_url"`
}

func GetStatus(w http.ResponseWriter, r *http.Request) {
	status := service.GetWindowStatus(time.Now())
	resp := StatusResponse{
		WindowStatus: status,
		FormURL:      os.Getenv("GOOGLE_FORM_URL"),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
