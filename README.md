# O2 Bento

Lunch order notification system for O2AI. Reminds team members to submit their daily lunch order via Google Form, tracks the submission window, and posts scheduled reminders to Discord.

## Services

| Service | Technology | Port |
|---------|-----------|------|
| Frontend | Next.js 14 (App Router), TypeScript, Tailwind CSS | 5111 |
| Backend | Go 1.23, Chi router | 5160 |
| Discord Bot | Go 1.22, discordgo | — |

## How It Works

- **Submission window:** 11:00 AM – 4:00 PM (Thai time, UTC+7), weekdays only
- **Friday special:** Friday's window also accepts Monday's order
- **One order per person per day** — enforced by the Google Form
- **Edit policy:** Use the Google Form "Edit response" link; do not re-submit
- **Pickup:** Pantry, 11:00 AM – 12:00 PM the following day

### Discord Notification Schedule

| Time | Message |
|------|---------|
| 11:00 AM | Order window is now open |
| 11:00 AM (Friday) | Open + reminder that Monday's order can be submitted |
| 3:30 PM | 30-minute warning |
| 4:00 PM | Order window is now closed |

## Prerequisites

- Docker & Docker Compose
- Discord bot token and channel ID (for the bot service)
- Google Form URL

## Setup

1. Copy and fill in the environment files:

```bash
# Backend
cp backend/.env.example backend/.env

# Discord bot
cp discord-bot/.env.example discord-bot/.env

# Frontend
cp frontend/.env.local.example frontend/.env.local
```

2. Set the required variables:

**`backend/.env`**
```
PORT=5160
GOOGLE_FORM_URL=https://forms.gle/...
GOOGLE_FORM_EDIT_URL=https://docs.google.com/forms/d/.../edit
TZ=Asia/Bangkok
```

**`discord-bot/.env`**
```
DISCORD_TOKEN=
DISCORD_CHANNEL_ID=
GOOGLE_FORM_URL=https://forms.gle/...
TZ=Asia/Bangkok
```

**`frontend/.env.local`**
```
NEXT_PUBLIC_GOOGLE_FORM_URL=https://forms.gle/...
NEXT_PUBLIC_API_URL=http://localhost:5160
PORT=5111
```

## Running

```bash
# Start all services
make up

# Run tests
make test

# Security audit
make security-check
```

Or run services individually:

```bash
# Frontend
cd frontend && npm install && npm run dev

# Backend
cd backend && go run ./cmd/server

# Discord bot
cd discord-bot && go run ./cmd/bot
```

## Project Structure

```
o2-bento/
├── frontend/                   # Next.js web app
│   └── src/
│       ├── app/                # App Router pages
│       ├── components/         # OrderStatus, countdown, form button
│       └── lib/schedule.ts     # Client-side window detection
├── backend/                    # Go REST API
│   └── internal/
│       └── service/schedule.go # Window open/close logic
├── discord-bot/                # Go Discord bot
│   └── internal/
│       └── scheduler/          # Cron-based reminder scheduler
└── docker-compose.yml
```
