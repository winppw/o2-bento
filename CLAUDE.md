# O2AI Launch Assistant

Notification system that reminds team members to submit their daily lunch order via a Google Form (organization email login required).

## Project Structure

```
o2ai-launch-assistant/
├── frontend/        # Next.js web app — order status, countdown, form link
├── backend/         # Go REST API — order state, schedule logic
├── discord-bot/     # Go Discord bot — scheduled reminders + slash commands
└── docker-compose.yml
```

## Business Rules

- **Submission window:** 11:00 AM – 4:00 PM (Thai time, UTC+7)
- **Deadline:** 4:00 PM daily — no submissions or edits after this
- **Friday special:** Friday window also accepts Monday's order
- **One order per person per day** — the Google Form enforces this
- **Edit policy:** Use the Google Form "Edit response" link; do NOT re-submit
- **Pickup:** Pantry, 11:00 AM – 12:00 PM the following day
- **Identification:** Each box is labeled with an order number; a paper sheet at the pantry maps names to numbers

## Google Form

- Restricted to organization email accounts (Google Workspace domain login required)
- The form link and edit link are configured via environment variables (`GOOGLE_FORM_URL`, `GOOGLE_FORM_EDIT_URL`)

## Notification Schedule (Discord bot)

| Time  | Message |
|-------|---------|
| 11:00 AM | "Lunch order window is now open! Submit by 4 PM." |
| 3:30 PM  | "30 minutes left to submit your lunch order." |
| 4:00 PM  | "Order window is now closed." |
| Friday 11:00 AM | Same open message + "You can also submit Monday's order." |

## Stack

| Layer | Technology |
|-------|-----------|
| Frontend | Next.js 14 (App Router), TypeScript, Tailwind CSS |
| Backend | Go 1.22, Chi router |
| Discord Bot | Go 1.22, discordgo |
| Containerization | Docker Compose |

## Environment Variables

### Backend (`backend/.env`)
```
PORT=8080
GOOGLE_FORM_URL=https://forms.gle/...
GOOGLE_FORM_EDIT_URL=https://docs.google.com/forms/d/.../edit
TZ=Asia/Bangkok
```

### Discord Bot (`discord-bot/.env`)
```
DISCORD_TOKEN=
DISCORD_CHANNEL_ID=
GOOGLE_FORM_URL=https://forms.gle/...
TZ=Asia/Bangkok
```

### Frontend (`frontend/.env.local`)
```
NEXT_PUBLIC_GOOGLE_FORM_URL=https://forms.gle/...
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Running Locally

```bash
docker compose up --build
```

Or per service:

```bash
# Frontend
cd frontend && npm install && npm run dev

# Backend
cd backend && go run ./cmd/server

# Discord bot
cd discord-bot && go run ./cmd/bot
```

## Key Files

- `backend/internal/service/schedule.go` — window open/close logic
- `discord-bot/internal/scheduler/reminder.go` — cron-based Discord notifications
- `frontend/src/components/OrderStatus.tsx` — live countdown + form link button
- `frontend/src/lib/schedule.ts` — client-side window detection (mirrors backend rules)
