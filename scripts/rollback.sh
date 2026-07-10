#!/usr/bin/env bash
# rollback.sh
# Rolls the backend-api Deployment back to its previous revision.
set -euo pipefail

: "${EKS_CLUSTER_NAME:=starttech-cluster}"
: "${AWS_REGION:=us-east-1}"

aws eks update-kubeconfig --name "${EKS_CLUSTER_NAME}" --region "${AWS_REGION}"

echo "==> Current rollout history"
kubectl rollout history deployment/backend-api

echo "==> Rolling back to previous revision"
kubectl rollout undo deployment/backend-api

echo "==> Waiting for rollout to stabilize"
kubectl rollout status deployment/backend-api

echo "Rollback complete."
