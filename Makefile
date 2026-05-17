.PHONY: test test-backend test-discord test-frontend security-check setup-hooks clean-branches check-ports up dev

## Run
check-ports:
	@bash .claude/scripts/check-ports.sh

up: check-ports
	docker compose up --build

dev: check-ports
	@echo "Starting services in dev mode..."
	@echo "  frontend → http://localhost:3000"
	@echo "  backend  → http://localhost:8080"
	docker compose up --build --watch

## Test
test: test-backend test-discord test-frontend

test-backend:
	cd backend && go test ./... -race -count=1

test-discord:
	cd discord-bot && CGO_ENABLED=0 go test ./... -count=1

test-frontend:
	cd frontend && npm test -- --passWithNoTests

test-ci:
	cd backend && go test ./... -race -count=1
	cd discord-bot && go test ./... -race -count=1
	cd frontend && npm run test:ci

security-check:
	@echo "==> go vet (backend)"
	cd backend && go vet ./...
	@echo "==> go vet (discord-bot)"
	cd discord-bot && go vet ./...
	@echo "==> govulncheck (backend)"
	cd backend && govulncheck ./... 2>/dev/null || echo "(install: go install golang.org/x/vuln/cmd/govulncheck@latest)"
	@echo "==> govulncheck (discord-bot)"
	cd discord-bot && govulncheck ./... 2>/dev/null || echo "(install: go install golang.org/x/vuln/cmd/govulncheck@latest)"
	@echo "==> npm audit (frontend)"
	cd frontend && npm audit --audit-level=high

setup-hooks:
	git config core.hooksPath .githooks
	chmod +x .githooks/pre-commit
	@echo "Git hooks configured — pre-commit guard is active."

clean-branches:
	@echo "==> Branches already merged into main (will be deleted):"
	@git branch --merged main | grep -v '^\*\|main' || echo "  (none)"
	@git branch --merged main | grep -v '^\*\|main' | xargs -r git branch -d
	@echo "==> Pruning stale remote-tracking refs..."
	@git fetch --prune 2>/dev/null || true
	@echo "Done. Run 'git branch -a' to verify."
