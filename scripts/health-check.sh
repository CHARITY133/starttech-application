#!/usr/bin/env bash
# health-check.sh
# Curls the unified CloudFront domain's /api/health endpoint.
# Usage: ./health-check.sh https://d1234abcd.cloudfront.net
set -euo pipefail

BASE_URL="${1:?Usage: $0 <base_url e.g. https://d1234abcd.cloudfront.net>}"
ENDPOINT="${BASE_URL%/}/api/health"

echo "==> Checking ${ENDPOINT}"

STATUS=$(curl -s -o /tmp/health-check-body.json -w "%{http_code}" "${ENDPOINT}" || echo "000")

if [ "${STATUS}" = "200" ]; then
  echo "Healthy (HTTP 200)"
  cat /tmp/health-check-body.json
  exit 0
else
  echo "Unhealthy (HTTP ${STATUS})"
  cat /tmp/health-check-body.json 2>/dev/null || true
  exit 1
fi
