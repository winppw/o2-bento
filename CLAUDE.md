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

## Port Range — 5100–5199

All services in this project use ports in the **5100–5199** range to avoid conflicts with other local services.

| Port | Service | Mnemonic |
|------|---------|----------|
| 5111 | Frontend (Next.js) | open hour = 11:00 AM |
| 5160 | Backend (Go API) | deadline = 16:00 |
| 5101–5159 | Reserved for future services | — |

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
DISCORD_CHANNEL_IDS=id1,id2,id3   # comma-separated; falls back to DISCORD_CHANNEL_ID
DISCORD_DM_USER_ID=
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

## Branching Strategy

`main` is the only permanent branch. All other branches are temporary and must be deleted after merging.

### Naming convention
| Type | Pattern | Example |
|------|---------|---------|
| Feature | `feat/<slug>` | `feat/discord-slash-commands` |
| Bug fix | `fix/<slug>` | `fix/friday-timezone-edge` |
| Chore / infra | `chore/<slug>` | `chore/update-deps` |

### Workflow

**Always work on a branch — never commit directly to `main`.**

```bash
# 1. Start from main
git checkout main && git pull

# 2. Create a working branch
git checkout -b feat/your-feature

# 3. Work, commit, test
make test

# 4. Push the branch and open a PR for human review
git push -u origin feat/your-feature
gh pr create --title "feat: your-feature" --body "$(cat <<'EOF'
## Summary
- What changed and why

## Test plan
- [ ] Tested locally
EOF
)"
```

**Do NOT merge into `main` yourself.** Open the PR and wait for a human to review and merge it. Branch deletion happens after merge.

## Key Files

- `backend/internal/service/schedule.go` — window open/close logic
- `discord-bot/internal/scheduler/reminder.go` — cron-based Discord notifications
- `frontend/src/components/OrderStatus.tsx` — live countdown + form link button
- `frontend/src/lib/schedule.ts` — client-side window detection (mirrors backend rules)
