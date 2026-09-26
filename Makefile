COMPOSE ?= docker compose

.PHONY: env up down logs ps migrate seed test

# create local config from the template if it doesn't exist yet
env:
	@test -f .env || (cp .env.example .env && echo "created .env from .env.example")

up: env
	$(COMPOSE) up --build -d

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f

ps:
	$(COMPOSE) ps

# admin processes run as one-off containers from the same release image and config
migrate: env
	$(COMPOSE) run --rm scheduler-postgres-migrate

seed: env
	$(COMPOSE) run --rm meeting-scheduler ./seed

test:
	go test ./...