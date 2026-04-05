COMPOSE ?= docker compose

POSTGRES_DB ?= scheduler_db
POSTGRES_USER ?= psg_user
POSTGRES_PASSWORD ?= psg_pass
POSTGRES_HOST ?= localhost
POSTGRES_PORT ?= 5432
POSTGRES_SSL_MODE ?= disable

.PHONY: up down logs ps seed

up:
	$(COMPOSE) up --build -d

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f

ps:
	$(COMPOSE) ps

seed:
	POSTGRES_DB=$(POSTGRES_DB) POSTGRES_USER=$(POSTGRES_USER) POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) POSTGRES_HOST=$(POSTGRES_HOST) POSTGRES_PORT=$(POSTGRES_PORT) POSTGRES_SSL_MODE=$(POSTGRES_SSL_MODE) go run ./cmd/seed

