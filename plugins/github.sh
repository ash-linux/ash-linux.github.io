#!/bin/bash
QUERY="${1#gh:}"
echo "GitHub issues for: $QUERY"
gh issue list --search "$QUERY" --limit 10 || echo "Error running gh CLI (is it installed and authenticated?)"
