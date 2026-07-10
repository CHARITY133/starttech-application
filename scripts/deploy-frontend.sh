#!/usr/bin/env bash
# deploy-frontend.sh
# Builds the React app and syncs it to S3, then invalidates CloudFront.
# Expects env vars: S3_BUCKET, CLOUDFRONT_DISTRIBUTION_ID
set -euo pipefail

: "${S3_BUCKET:?Set S3_BUCKET to the frontend bucket name}"
: "${CLOUDFRONT_DISTRIBUTION_ID:?Set CLOUDFRONT_DISTRIBUTION_ID}"

cd "$(dirname "$0")/../frontend"

echo "==> Installing dependencies"
npm ci

echo "==> Running npm audit (non-fatal, informational)"
npm audit --audit-level=high || true

echo "==> Building production bundle"
npm run build

echo "==> Syncing build/ to s3://${S3_BUCKET}"
aws s3 sync build/ "s3://${S3_BUCKET}" --delete

echo "==> Invalidating CloudFront cache"
aws cloudfront create-invalidation \
  --distribution-id "${CLOUDFRONT_DISTRIBUTION_ID}" \
  --paths "/*"

echo "Frontend deployed."
