.PHONY: infra-up infra-down infra-logs validate

infra-up:
	docker compose up -d postgres rabbitmq

infra-down:
	docker compose down

infra-logs:
	docker compose logs -f postgres rabbitmq

validate:
	docker compose config
