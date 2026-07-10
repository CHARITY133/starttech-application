# starttech-application

Application source + Kubernetes manifests + CI/CD for the StartTech full-stack
app (React frontend, Golang backend), deployed onto the `starttech-infra`
platform.

## Getting the source in

Pull the actual frontend/backend code from the existing app repo's
`feature/full-stack` branch into this repo's `frontend/` and `backend/`
folders:

```bash
git clone --branch feature/full-stack \
  https://github.com/Innocent9712/much-to-do.git /tmp/much-to-do

# Adjust these paths to match that repo's actual layout once you inspect it
cp -r /tmp/much-to-do/frontend/*  ./frontend/
cp -r /tmp/much-to-do/backend/*   ./backend/
```

Then adapt the two apps to this assessment's requirements (see below) before
committing.

## Frontend requirements checklist

- [ ] All API calls use **relative paths**, e.g. `fetch('/api/v1/health')`,
      never a hardcoded ALB or domain — this is what makes the unified
      CloudFront setup work (no mixed-content, no CORS headaches).
- [ ] `npm run build` outputs to `frontend/build/` (Create React App default;
      adjust `deploy-frontend.sh` / the workflow if using Vite's `dist/`).
- [ ] `package.json` has a working `build` script.

## Backend requirements checklist

- [ ] Exposes `GET /api/v1/health` (or `/health`) returning 200 with a small
      JSON body — used by k8s readiness/liveness probes, the ALB target group
      health check, and `scripts/health-check.sh`.
- [ ] Reads `REDIS_HOST` and `MONGO_URI` from environment variables (populated
      via the `backend-secrets` Kubernetes Secret referenced in
      `k8s/deployment.yaml` — create it separately, don't commit real values:
      ```bash
      kubectl create secret generic backend-secrets \
        --from-literal=redis-host=<elasticache-endpoint> \
        --from-literal=mongo-uri='<mongodb atlas connection string>'
      ```
      ).
- [ ] Listens on port `8080` (matches `k8s/deployment.yaml` /
      `k8s/service.yaml` and the ECR/EKS container port convention used by
      the infra repo).
- [ ] Logs structured JSON to stdout (e.g. via `log/slog` or `zerolog`), so
      Container Insights / FluentBit can parse it.
- [ ] `backend/Dockerfile` builds a small production image (multi-stage build
      recommended: `golang:1.22` builder -> `gcr.io/distroless/static` or
      `alpine` runtime).

## Deploying

CI/CD does this automatically on push to `main`, but locally:

```bash
# Frontend
S3_BUCKET=<bucket-name> \
CLOUDFRONT_DISTRIBUTION_ID=<dist-id> \
./scripts/deploy-frontend.sh

# Backend
ECR_REGISTRY=<account>.dkr.ecr.us-east-1.amazonaws.com \
./scripts/deploy-backend.sh

# Verify
./scripts/health-check.sh https://<your-cloudfront-domain>

# Roll back if needed
./scripts/rollback.sh
```

## Required GitHub Secrets

| Secret | Used by |
|---|---|
| `AWS_DEPLOY_ROLE_ARN` | both workflows (OIDC role assumption) |
| `S3_FRONTEND_BUCKET` | `frontend-ci-cd.yml` |
| `CLOUDFRONT_DISTRIBUTION_ID` | `frontend-ci-cd.yml` |

All come from `terraform output` in `starttech-infra`.
