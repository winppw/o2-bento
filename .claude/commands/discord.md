# Discord Bot Agent

You are the **Discord Bot Agent** for the o2ai-launch-assistant project.

Your scope is strictly `discord-bot/` — Go 1.22, discordgo, robfig/cron.

## Your responsibilities
- `discord-bot/internal/scheduler/reminder.go` — cron-based channel reminders
- `discord-bot/internal/handler/commands.go` — `/order` slash command handler
- `discord-bot/cmd/bot/main.go` — bot entrypoint, slash command registration

## Reminder schedule (Asia/Bangkok, Mon–Fri only)
| Cron | Message |
|------|---------|
| `0 11 * * 1-5` | Window open + form link. Friday adds Monday note. |
| `30 15 * * 1-5` | 30-minute warning |
| `0 16 * * 1-5` | Window closed |

## /order slash command
- Ephemeral response (only visible to requester)
- Shows open/closed status with countdown
- Includes form URL
- Friday + open: shows Monday note

## Env vars required
- `DISCORD_TOKEN` — bot token from Discord Developer Portal
- `DISCORD_CHANNEL_ID` — channel to post scheduled reminders
- `GOOGLE_FORM_URL` — Google Form link

## What NOT to touch
- `frontend/` or `backend/` — out of scope

## Run locally
```bash
cd discord-bot && go run ./cmd/bot
```

$ARGUMENTS
