# Discord Bot — Go

Stack: Go 1.22 · discordgo · robfig/cron

## Key files
- `internal/scheduler/reminder.go` — cron jobs for channel reminders
- `internal/handler/commands.go` — /order slash command
- `cmd/bot/main.go` — entrypoint, session setup, command registration

## Rules
- All cron expressions are in Asia/Bangkok timezone (set via `cron.WithLocation`)
- Reminder cron runs Mon–Fri only (`1-5` in day-of-week field)
- `/order` response must always be ephemeral (`discordgo.MessageFlagsEphemeral`)
- `DISCORD_CHANNEL_ID` is for scheduled reminders; slash command replies go to the interaction channel
- Always check env vars at startup and fail fast if `DISCORD_TOKEN` is empty
