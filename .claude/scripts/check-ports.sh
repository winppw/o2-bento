#!/bin/bash
# Verify required ports are free before starting services.
# Usage: bash .claude/scripts/check-ports.sh
#        PORT_FRONTEND=3001 PORT_BACKEND=9090 bash .claude/scripts/check-ports.sh

FRONTEND_PORT="${PORT_FRONTEND:-3000}"
BACKEND_PORT="${PORT_BACKEND:-8080}"

FAILED=0

check_port() {
  local port="$1"
  local service="$2"

  if lsof -iTCP:"$port" -sTCP:LISTEN -n -P &>/dev/null 2>&1; then
    local owner
    owner=$(lsof -iTCP:"$port" -sTCP:LISTEN -n -P 2>/dev/null | awk 'NR==2 {print $1" (PID "$2")"}')
    echo "  BUSY  :$port  $service — already used by $owner"
    echo "         kill it: kill \$(lsof -ti:$port)"
    FAILED=1
  else
    echo "  FREE  :$port  $service"
  fi
}

echo "Checking ports..."
check_port "$FRONTEND_PORT" "frontend (Next.js)"
check_port "$BACKEND_PORT"  "backend  (Go API) "

if [ "$FAILED" -ne 0 ]; then
  echo ""
  echo "One or more ports are in use. Free them before starting services."
  exit 1
fi

echo "All ports are free."
