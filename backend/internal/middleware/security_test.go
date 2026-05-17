package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/o2ai/launch-assistant/backend/internal/middleware"
)

func noop(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }

func TestSecurityHeaders(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	middleware.SecurityHeaders(http.HandlerFunc(noop)).ServeHTTP(w, req)

	headers := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"Referrer-Policy":        "strict-origin-when-cross-origin",
		"Permissions-Policy":     "camera=(), microphone=(), geolocation=()",
	}
	for key, want := range headers {
		if got := w.Header().Get(key); got != want {
			t.Errorf("%s = %q, want %q", key, got, want)
		}
	}
}

func TestCORS_DefaultWildcard(t *testing.T) {
	os.Unsetenv("ALLOWED_ORIGIN")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	middleware.CORS(http.HandlerFunc(noop)).ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("ACAO = %q, want *", got)
	}
}

func TestCORS_RestrictedOrigin(t *testing.T) {
	os.Setenv("ALLOWED_ORIGIN", "https://launch.example.com")
	t.Cleanup(func() { os.Unsetenv("ALLOWED_ORIGIN") })

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	middleware.CORS(http.HandlerFunc(noop)).ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "https://launch.example.com" {
		t.Errorf("ACAO = %q, want https://launch.example.com", got)
	}
}

func TestCORS_Preflight(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	w := httptest.NewRecorder()

	middleware.CORS(http.HandlerFunc(noop)).ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("preflight status = %d, want 204", w.Code)
	}
}
