.PHONY: test test-backend test-discord test-frontend security-check setup-hooks

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
