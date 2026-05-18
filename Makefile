include .env

SHELL = /bin/sh
UID := $(shell id -u)
COMPOSE = docker compose -p recru_backend -f docker-compose.local.yaml

.PHONY: docker-up docker-down docker-restart docker-stop \
        go-bash db-bash db-console redis-bash redis-cli \
        composer-install composer-update composer-require composer-require-d composer-remove \
        laravel-setup project-installation \
        db-migrate db-migrate-force db-reset db-rollback db-fresh db-fresh-seed db-make-migration db-seed db-seed-class \
        logs logs-container logs-laravel containers-status \
        help test

# === DOCKER OPERATIONS ===
up:
	@env UID=${UID} $(COMPOSE) up -d --remove-orphans

down:
	@env UID=${UID} $(COMPOSE) down

restart: docker-down docker-up

stop:
	@env UID=${UID} $(COMPOSE) stop

# === CONTAINER ACCESS ===
go:
	@env UID=${UID} $(COMPOSE) exec app bash

db:
	@env UID=${UID} $(COMPOSE) exec db bash

# usage: make dbc username=YOUR_USERNAME
db-c:
	@env UID=${UID} $(COMPOSE) exec db psql -U $(user) -d db

redis:
	@env UID=${UID} $(COMPOSE) exec redis bash

redis-c:
	@env UID=${UID} $(COMPOSE) exec redis redis-cli

# === HELP ===
help:
	@echo "Makefile Commands:"
	@echo ""
	@echo "  🚀 Docker Operations:"
	@echo "    up          - Start Docker containers"
	@echo "    down        - Stop and remove Docker containers"
	@echo "    restart     - Restart Docker containers"
	@echo "    stop        - Stop Docker containers"
	@echo ""
	@echo "  🐳 Container Access:"
	@echo "    go         - Access Go container bash"
	@echo "    db         - Access database container bash"
	@echo "    db-c       - Access PostgreSQL console"
	@echo "    redis      - Access Redis container bash"
	@echo "    redis-c    - Access Redis CLI"