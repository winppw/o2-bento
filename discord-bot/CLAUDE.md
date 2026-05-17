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
- `DISCORD_CHANNEL_ID` posts reminders to a server channel; `DISCORD_DM_USER_ID` DMs a specific user — both are optional and independent, but at least one must be set for reminders to be delivered
- Slash command replies go to the interaction channel, unaffected by the above env vars
- Always check env vars at startup and fail fast if `DISCORD_TOKEN` is empty
