#!/usr/bin/env bash
# deploy-backend.sh
# Builds, tags, and pushes the backend image, then rolls it out to EKS.
# Expects env vars: ECR_REGISTRY, ECR_REPOSITORY, EKS_CLUSTER_NAME, AWS_REGION
set -euo pipefail

: "${ECR_REGISTRY:?e.g. 123456789012.dkr.ecr.us-east-1.amazonaws.com}"
: "${ECR_REPOSITORY:=starttech-backend-api}"
: "${EKS_CLUSTER_NAME:=starttech-cluster}"
: "${AWS_REGION:=us-east-1}"

GIT_SHA=$(git rev-parse --short HEAD)
IMAGE_URI="${ECR_REGISTRY}/${ECR_REPOSITORY}:${GIT_SHA}"

cd "$(dirname "$0")/../backend"

echo "==> Running Go tests"
go test ./...

echo "==> Building Docker image ${IMAGE_URI}"
docker build -t "${IMAGE_URI}" .

echo "==> Authenticating to ECR"
aws ecr get-login-password --region "${AWS_REGION}" \
  | docker login --username AWS --password-stdin "${ECR_REGISTRY}"

echo "==> Pushing image"
docker push "${IMAGE_URI}"

echo "==> Updating kubeconfig"
aws eks update-kubeconfig --name "${EKS_CLUSTER_NAME}" --region "${AWS_REGION}"

echo "==> Patching deployment image"
cd ../k8s
kubectl set image deployment/backend-api backend-api="${IMAGE_URI}"

echo "==> Applying manifests"
kubectl apply -f .

echo "==> Waiting for rollout"
kubectl rollout status deployment/backend-api

echo "Backend deployed: ${IMAGE_URI}"
