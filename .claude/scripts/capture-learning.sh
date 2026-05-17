#!/bin/bash
# Appends a timestamped entry to learning-log.md when Claude stops a session.
# Reads CLAUDE_SUMMARY env var if set (populated by Stop hook input).

LOG="$(dirname "$0")/../learning-log.md"
DATE=$(TZ=Asia/Bangkok date '+%Y-%m-%d %H:%M')
SUMMARY="${CLAUDE_SUMMARY:-Session ended (no summary captured).}"

if [ ! -f "$LOG" ]; then
  echo "# Learning Log\n\nDecisions, fixes, and patterns discovered during development.\n" > "$LOG"
fi

printf "\n## %s\n%s\n" "$DATE" "$SUMMARY" >> "$LOG"
