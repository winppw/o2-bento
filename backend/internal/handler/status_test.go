package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/o2ai/launch-assistant/backend/internal/handler"
)

func TestGetStatus_ReturnsJSON(t *testing.T) {
	os.Setenv("GOOGLE_FORM_URL", "https://forms.gle/test")
	t.Cleanup(func() { os.Unsetenv("GOOGLE_FORM_URL") })

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	w := httptest.NewRecorder()

	handler.GetStatus(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}

	var body map[string]any
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	for _, field := range []string{"is_open", "open_at", "close_at", "form_url"} {
		if _, ok := body[field]; !ok {
			t.Errorf("response missing field %q", field)
		}
	}
	if body["form_url"] != "https://forms.gle/test" {
		t.Errorf("form_url = %v, want https://forms.gle/test", body["form_url"])
	}
}

func TestGetStatus_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/status", nil)
	w := httptest.NewRecorder()

	// Chi handles method routing; here we test the handler itself just returns a valid response
	handler.GetStatus(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("handler should not enforce method itself, chi router does")
	}
}
