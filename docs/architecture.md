# Architecture

## Initial boundary

```text
Next.js dashboard
        |
        v
     Go API -------> PostgreSQL
        |
        v
    RabbitMQ ------> Python worker
        ^                |
        |                v
        +---------- PostgreSQL
```

The first implementation will use one workflow type: **research brief generation**.
Keeping one end-to-end path makes the system easy to explain in an interview while still
demonstrating synchronous APIs, asynchronous messaging, persistence, and a frontend.

## Phase boundaries

1. Repository and local infrastructure
2. Go API health endpoint and workflow model
3. RabbitMQ publisher and Python worker
4. Next.js dashboard
5. End-to-end workflow, tests, and observability
