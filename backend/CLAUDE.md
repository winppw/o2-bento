# Backend — Go API

Stack: Go 1.22 · Chi router · godotenv

## Key files
- `internal/service/schedule.go` — window logic (source of truth)
- `internal/handler/status.go` — GET /status handler
- `cmd/server/main.go` — entrypoint + CORS middleware

## Rules
- All time calculations must `time.LoadLocation("Asia/Bangkok")` — never UTC
- `schedule.GetWindowStatus(time.Now())` is the only time entry point
- CORS is open (`*`) for now — restrict when deploying to production
- Env vars loaded via godotenv at startup; never hard-code form URLs
