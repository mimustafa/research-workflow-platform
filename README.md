# Research Workflow Platform

A portfolio project for demonstrating a small, production-minded research workflow system.

The platform will accept a research request, coordinate asynchronous processing through
RabbitMQ, run a Python analysis worker, persist workflow state in PostgreSQL, and present
results through a Next.js web application. Go will provide the backend API and workflow
orchestration boundary.

## Planned stack

- **Go**: HTTP API, validation, workflow state transitions
- **Python**: asynchronous analysis worker
- **Next.js + React + TypeScript**: user-facing dashboard
- **RabbitMQ**: durable work queue and event transport
- **PostgreSQL**: workflow and result persistence
- **Docker Compose**: reproducible local development environment
- **Tests and observability**: health checks, structured logs, and focused automated tests

## Repository layout

```text
apps/
  api/       Go API and workflow orchestration
  worker/    Python RabbitMQ consumer and analysis
  web/       Next.js dashboard
infra/
  postgres/ Database initialization and local infrastructure notes
docs/        Architecture and demo walkthrough
```

## Development status

This repository is being built in small phases. The first phase establishes the structure
and local infrastructure. Application behavior will be added incrementally with a focused
commit and test for each phase.

## Quick start

```bash
cp .env.example .env
docker compose up -d postgres rabbitmq
docker compose ps
```

The initial infrastructure exposes:

- PostgreSQL on `localhost:5432`
- RabbitMQ AMQP on `localhost:5672`
- RabbitMQ management UI on `http://localhost:15672`

Default local credentials are in `.env.example` and are for development only.

To run the current Go API in Docker:

```bash
docker compose up --build
curl http://localhost:8080/health
```

Open the dashboard at `http://localhost:3000`.

For the complete interview walkthrough:

```bash
./scripts/demo.sh
```