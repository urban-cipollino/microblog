# Makefile
ifneq (,$(wildcard .env))
    include .env
    export
endif

DC             = docker compose
MIGRATIONS_DIR = migrations
DB_URL         = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
DB_TEST_NAME   = $(DB_NAME)_test

## —— Makefile  —————————————————————————————————————————————————————
help: ## Display description of make (cmd by default)
	@grep -E '(^[a-zA-Z0-9\./_-]+:.*?##.*$$)|(^##)' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}{printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | sed -e 's/\[32m##/[33m/'

## —— Linter Golanggci —————————————————————————————————————————————————————
lint: ## Run linter
	@golangci-lint run
cs: ## Run linter with auto fix
	@golangci-lint run --fix ./...

## —— Wire —————————————————————————————————————————————————————
wire: ## Update DI
	@wire ./internal/app

## —— Tests —————————————————————————————————————————————————————
race: ## Run race tests
	@ENV=test go test -race ./...

test: ## Run tests
	@ENV=test go test ./...

test.db.add: ## Create test database
	@$(DC) exec -T db psql -U $(DB_USER) -d postgres -c "CREATE DATABASE $(DB_TEST_NAME);"
	@echo "Test database $(DB_TEST_NAME) created"

test.db.drop: ## Drop test database
	@$(DC) exec -T db psql -U $(DB_USER) -d postgres -c "DROP DATABASE IF EXISTS $(DB_TEST_NAME);"
	@echo "Test database $(DB_TEST_NAME) dropped"

test.db.recreate: test.db.drop test.db.add  ## Drop and create test database

bench: ## Run benchmarks
	@ENV=test go test -bench=. -benchmem ./...

check: ## Runs Linter and tests
	@out="$$( $(MAKE) cs --no-print-directory 2>&1 )" || { echo "$$out"; exit 1; }
	@out="$$( $(MAKE) race --no-print-directory 2>&1 )" || { echo "$$out"; exit 1; }
	@$(MAKE) bench
	@echo "success"

### —— Docker 🐳 ————————————————————————————————————————————————————————————————
dc.build: ## Build images
	@$(DC) build --pull --no-cache

dc.start: ## Start stopped containers
	@$(DC) start
dc.stop: ## Stop containers
	@$(DC) stop

dc.up: ## Create and up containers with orphans removing
	@$(DC) up -d --remove-orphans
dc.down: ## Delete containers
	@$(DC) down --remove-orphans

dc.restart: dc.down dc.up ## Restart containers

dc.rebuild: dc.down dc.build dc.up ## Rebuild containers with up

dc.logs: ## Show logs
	@$(DC) logs --tail=0 --follow

dc.drop: ## Delete containers and volumes
	@echo "WARNING: This command will purge all containers and volumes! Proceed? (y/n)"
	@read answer && [ $$answer = y ] && $(DC) down -v --remove-orphans || echo "Cancelled."
dc.full_drop: ## Delete containers, images and volumes
	@echo "WARNING: This command will purge all containers, images and volumes! Proceed? (y/n)"
	@read answer && [ $$answer = y ] && $(DC) down -v --rmi=all --remove-orphans || echo "Cancelled."

## Goose 🪿 —————————————————————————————————————————————————————
mg: ## Create migration
	@goose -dir $(MIGRATIONS_DIR) create $(name) sql

mg.up: ## Run migration
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

mg.down: ## Rollback migration
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

mg.status: ## Migration status
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" status

mg.reset: ## Rollback all migration
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" reset

mg.redo: ## Rollback last migration and run it again
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" redo
