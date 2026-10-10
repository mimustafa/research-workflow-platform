# OpenBao local demonstration

OpenBao is included as an **optional local development service** to demonstrate secret
management without changing the simple Docker Compose startup path.

## Start the development server

```bash
docker compose -f compose.yaml -f compose.openbao.yaml up --build -d
export BAO_ADDR=http://127.0.0.1:8200
export BAO_TOKEN=local-root-token
```

Store an example platform credential:

```bash
docker compose -f compose.yaml -f compose.openbao.yaml exec openbao \
  bao kv put secret/research-workflow \
  postgres_password='replace-this' rabbitmq_password='replace-this'
```

Read it:

```bash
docker compose -f compose.yaml -f compose.openbao.yaml exec openbao \
  bao kv get secret/research-workflow
```

## Security boundary

The provided server runs in OpenBao development mode with a known root token, no durable
storage, and no application injection. It is deliberately limited to a local demonstration.

For a real deployment, use initialized and unsealed OpenBao, Kubernetes authentication,
short-lived application policies, TLS, durable storage, and an injector or agent to deliver
dynamic secrets to the API and worker.
