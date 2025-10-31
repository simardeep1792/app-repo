# Progressive Delivery Demo App

HTTP service demonstrating canary deployments with automated rollback capabilities.

## Overview

This Go application provides:
- HTTP endpoints: `/`, `/healthz`, `/ready`
- Prometheus metrics for request count and latency
- Configurable failure injection for testing rollbacks
- Multi-stage Docker build for minimal image size

## Local Development

### Run with Docker

```bash
# Build and run
docker build -t app .
docker run -p 8080:8080 app

# Access endpoints
curl http://localhost:8080/
curl http://localhost:8080/healthz
curl http://localhost:8080/metrics
```

### Run with Go

```bash
go mod download
go run src/main.go
```

## Failure Injection

The app can simulate failures for testing rollback scenarios:

```bash
# Enable 10% failure rate on / endpoint
make toggle-bad

# Disable failure injection
make toggle-good
```

When `INJECT_FAILURE=true`, approximately 10% of requests to `/` return HTTP 500.

## CI/CD Pipeline

### Build Process

On push to main:
1. Build Docker image with caching
2. Scan with Trivy (fail on critical vulnerabilities)
3. Push to GHCR with tags:
   - `main` (latest from main branch)
   - `main-<gitsha>` (specific commit)

### Environment Updates

After successful build, CI automatically:
1. Checks out env-repo
2. Updates `envs/dev/rollout.yaml` with new image tag
3. Opens PR titled "chore: bump app image to <gitsha>"

### Required Secrets

- `GITHUB_TOKEN`: Automatically provided, needs packages:write
- `ENV_REPO_TOKEN`: PAT with repo access to env-repo (if different owner)

## Makefile Targets

- `build`: Build Docker image locally
- `push`: Build and push to registry
- `toggle-bad`: Enable failure injection
- `toggle-good`: Disable failure injection  
- `bump-dev`: Update local env-repo with new image (for testing)

## Metrics

The application exposes Prometheus metrics at `/metrics`:

- `http_requests_total{status,route}`: Request counter by status and route
- `http_request_duration_seconds{route}`: Request duration histogram

These metrics are used by Argo Rollouts AnalysisTemplates to determine deployment health.

## Environment Variables

- `PORT`: HTTP server port (default: 8080)
- `VERSION`: Version string shown on home page (default: 1.0.0)
- `INJECT_FAILURE`: Enable failure injection (default: false)
- `IMAGE_REGISTRY`: Docker registry for push (default: ghcr.io)
- `IMAGE_REPO`: Repository name (default: from GITHUB_REPOSITORY)
- `IMAGE_TAG`: Image tag (default: git short SHA)