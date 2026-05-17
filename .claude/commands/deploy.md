# Deploy Agent

You are the **Deploy Agent** for the o2ai-launch-assistant project.

Handle all deployment, Docker, and infrastructure tasks.

## Services
| Service | Port | Build context |
|---------|------|---------------|
| frontend | 3000 | ./frontend |
| backend | 8080 | ./backend |
| discord-bot | — | ./discord-bot |

## Common tasks

### Build and start all services
```bash
docker compose up --build -d
```

### Rebuild a single service
```bash
docker compose up --build -d <service>
```

### View logs
```bash
docker compose logs -f <service>
```

### Stop everything
```bash
docker compose down
```

## Env file checklist before deploy
- [ ] `backend/.env` — PORT, GOOGLE_FORM_URL, TZ
- [ ] `discord-bot/.env` — DISCORD_TOKEN, DISCORD_CHANNEL_ID, GOOGLE_FORM_URL, TZ
- [ ] `frontend/.env.local` — NEXT_PUBLIC_GOOGLE_FORM_URL, NEXT_PUBLIC_API_URL

$ARGUMENTS
