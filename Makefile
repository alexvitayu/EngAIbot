ifneq (,$(wildcard .env.dev))
    include .env.dev ## это загружает переменные окружения в Makefile
    export ## это делает переменные окружения из .env.dev доступными в shell командах
endif
# ====================
# BUILD
# ====================
build:
	go build -o bin/tgbot ./cmd

# ====================
# INSTALL binary file to GOPATH or GOBIN
# ====================
install: build
	go install ./cmd

# ====================
# Launch golangci linter on base of .golangci.yaml
# ====================
lint:
	golangci-lint run

# ====================
# Run tests
# ====================
test:
	go test ./... -v -race -cover

# ====================
# Работа только с БД из docker-compose-dev.yaml
# ====================
dev-db-up:
	docker compose --env-file .env.dev -f docker-compose-dev.yaml up -d
	@echo "✅ PostgreSQL запущена на $(DB_HOST):$(DB_PORT)"

# Остановка БД
dev-db-down:
	docker compose --env-file .env.dev -f docker-compose-dev.yaml down

# Статус БД
dev-db-status:
	docker compose --env-file .env.dev -f docker-compose-dev.yaml ps

# ====================
# Работа только с БД без compose (для эксперимента) подключаюсь к той же БД tgbot-postgres-db
# ====================
db-up:
	docker run -d \
      --name tgbot-postgres-db \
      -e POSTGRES_DB=$(DB_NAME) \
      -e POSTGRES_USER=$(DB_USER) \
      -e POSTGRES_PASSWORD=$(DB_PASSWORD) \
      -p $(DB_PORT):5432 \
      -v engbot_dev_data:/var/lib/postgresql/data \
      postgres:15-alpine

db-down:
	docker stop tgbot-postgres-db

db-del:
	docker rm tgbot-postgres-db

# ====================
# GOOSE migrations handling (DDL operations)
# ====================
migrate-up:
	goose -dir ./migrations postgres $(DATABASE_URL) up

migrate-down:
	goose -dir ./migrations postgres $(DATABASE_URL) down

migrate-status:
	goose -dir ./migrations postgres $(DATABASE_URL) status