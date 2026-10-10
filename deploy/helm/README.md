# Helm deployment

This chart deploys the web application, Go API, Python worker, PostgreSQL, and RabbitMQ
to Kubernetes. It is designed for a local Kind cluster or interview demonstration, not
as a production database or message-broker topology.

## Prerequisites

1. Build the three application images with `docker compose build`.
2. Load them into a local Kind cluster, for example:

   ```bash
   kind load docker-image research-workflow-platform-api:latest \
     research-workflow-platform-worker:latest \
     research-workflow-platform-web:latest --name research-workflow
   ```

3. Override the demonstration passwords before installation.

## Install

```bash
helm upgrade --install research-workflow ./deploy/helm/research-workflow \
  --namespace research-workflow --create-namespace \
  --set credentials.postgresPassword='replace-this' \
  --set credentials.rabbitmqPassword='replace-this'
```

Access the services through local port forwarding:

```bash
kubectl -n research-workflow port-forward service/research-workflow-research-workflow-web 3000:3000
kubectl -n research-workflow port-forward service/research-workflow-research-workflow-api 8080:8080
```

The dashboard expects the API at `http://localhost:8080`, so keep both port-forward
commands running while using the local Kubernetes demo.

## OpenBao

The chart can create an **unsafe development-only** OpenBao server:

```bash
helm upgrade --install research-workflow ./deploy/helm/research-workflow \
  --namespace research-workflow --create-namespace \
  --set openbao.enabled=true
```

It exists to demonstrate the deployment surface only. It is not initialized, unsealed,
or configured for application secret injection; do not use this mode outside a local demo.
