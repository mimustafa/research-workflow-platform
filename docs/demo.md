# Interview demo

## Start-to-finish walkthrough

From the repository root:

```bash
cp .env.example .env
./scripts/demo.sh
```

Then open `http://localhost:3000`, submit another brief, and show the workflow moving
through the API, RabbitMQ, Python worker, and PostgreSQL-backed dashboard.

Useful commands while explaining the system:

```bash
docker compose ps
docker compose logs -f api worker
docker compose exec postgres psql -U research -d research_workflows \
  -c "SELECT id, status, result FROM workflows ORDER BY created_at DESC;"
```

## Five-minute explanation

1. The Next.js page submits a brief to the Go API.
2. The API validates and stores a `queued` workflow in PostgreSQL.
3. The API publishes the workflow JSON as a durable RabbitMQ message.
4. The Python worker acknowledges the message and records a deterministic result.
5. The dashboard reads the persisted workflow state back from the API.

The design deliberately uses one workflow type so the system is easy to understand while
still demonstrating a synchronous API, asynchronous processing, persistence, containers,
and a typed frontend.
