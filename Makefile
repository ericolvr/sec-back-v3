# ANSI color codes
COLOR_RESET=\033[0m
COLOR_BOLD=\033[1m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m
COLOR_BLUE=\033[34m
COLOR_RED=\033[31m

# Variáveis
MAIN_PATH=cmd/main.go
ENV ?= local
ENV_FILE = .env.$(ENV)
DB_VOLUME_NAME=postgres_data_nr1

.PHONY: help run install db-start db-stop db-clean db-migrate db-reset build test dev


help:
	@echo ""
	@echo "  $(COLOR_YELLOW)Available targets:$(COLOR_RESET)"
	@echo "  $(COLOR_BLUE)Local Development:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)install$(COLOR_RESET)		- Install dependencies"
	@echo "  $(COLOR_GREEN)run$(COLOR_RESET)			- Run development server"
	@echo "  $(COLOR_GREEN)dev$(COLOR_RESET)			- Setup and run complete development environment"
	@echo ""
	@echo "  $(COLOR_BLUE)Database Management:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)db-start$(COLOR_RESET)		- Start Postgres container"
	@echo "  $(COLOR_GREEN)db-stop$(COLOR_RESET)		- Stop and remove database container"
	@echo "  $(COLOR_GREEN)db-clean$(COLOR_RESET)		- Clean database data (remove container + volume)"
	@echo "  $(COLOR_GREEN)db-migrate$(COLOR_RESET)		- Run database migrations"
	@echo "  $(COLOR_GREEN)db-reset$(COLOR_RESET)		- Reset database (clean + start + migrate)"
	@echo ""
	@echo "  $(COLOR_BLUE)Build & Test:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)build$(COLOR_RESET)		- Build the application"

	@echo ""

install:
	@echo "$(COLOR_YELLOW)Installing Go dependencies...$(COLOR_RESET)"
	go mod download
	go mod tidy
	@echo "$(COLOR_GREEN)✅ Dependencies installed successfully!$(COLOR_RESET)"

run:
	@echo "$(COLOR_YELLOW)Starting development server with $(ENV) environment...$(COLOR_RESET)"
	@if [ "$(ENV)" = "local" ]; then \
		if [ ! -f .env ]; then \
			if [ -f .env-sample ]; then \
				echo "$(COLOR_YELLOW)⚠️  .env not found, creating from .env-sample...$(COLOR_RESET)"; \
				cp .env-sample .env; \
			else \
				echo "$(COLOR_RED)❌ Neither .env nor .env-sample found!$(COLOR_RESET)"; \
				exit 1; \
			fi; \
		fi; \
		echo "$(COLOR_BLUE)Using .env file$(COLOR_RESET)"; \
	else \
		if [ ! -f $(ENV_FILE) ]; then \
			echo "$(COLOR_YELLOW)⚠️  Environment file $(ENV_FILE) not found, creating from .env-sample...$(COLOR_RESET)"; \
			cp .env-sample $(ENV_FILE); \
		fi; \
		echo "$(COLOR_BLUE)Loading environment from: $(ENV_FILE)$(COLOR_RESET)"; \
		cp $(ENV_FILE) .env; \
	fi
	go run $(MAIN_PATH)

db-start:
	@echo "$(COLOR_YELLOW)Starting Postgres container for $(ENV) environment...$(COLOR_RESET)"
	@if [ "$(ENV)" = "local" ]; then \
		if [ ! -f .env ]; then \
			if [ -f .env-sample ]; then \
				echo "$(COLOR_YELLOW)⚠️  .env not found, creating from .env-sample...$(COLOR_RESET)"; \
				cp .env-sample .env; \
			else \
				echo "$(COLOR_RED)❌ Neither .env nor .env-sample found!$(COLOR_RESET)"; \
				exit 1; \
			fi; \
		fi; \
		echo "$(COLOR_BLUE)Using .env file$(COLOR_RESET)"; \
	else \
		if [ ! -f $(ENV_FILE) ]; then \
			echo "$(COLOR_YELLOW)⚠️  Environment file $(ENV_FILE) not found, creating from .env-sample...$(COLOR_RESET)"; \
			cp .env-sample $(ENV_FILE); \
		fi; \
		echo "$(COLOR_BLUE)Loading environment from: $(ENV_FILE)$(COLOR_RESET)"; \
		cp $(ENV_FILE) .env; \
	fi
	docker compose --env-file .env up postgres -d
	@echo "$(COLOR_GREEN)✅ Database container started!$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)Database: $$(grep DB_NAME .env | cut -d'=' -f2) on localhost:$$(grep DB_PORT .env | cut -d'=' -f2)$(COLOR_RESET)"

db-stop:
	@echo "$(COLOR_YELLOW)Stopping and removing database container...$(COLOR_RESET)"
	docker compose down postgres
	@echo "$(COLOR_GREEN)✅ Database container removed!$(COLOR_RESET)"

db-clean:
	@echo "$(COLOR_YELLOW)Cleaning database data...$(COLOR_RESET)"
	docker compose down postgres
	docker volume rm $(DB_VOLUME_NAME) 2>/dev/null || true
	@echo "$(COLOR_GREEN)✅ Database data cleaned!$(COLOR_RESET)"

db-migrate:
	@echo "$(COLOR_YELLOW)Running database migrations...$(COLOR_RESET)"
	@if [ ! -f .env ]; then \
		if [ -f .env-sample ]; then \
			echo "$(COLOR_YELLOW)⚠️  .env not found, creating from .env-sample...$(COLOR_RESET)"; \
			cp .env-sample .env; \
		else \
			echo "$(COLOR_RED)❌ Neither .env nor .env-sample found!$(COLOR_RESET)"; \
			exit 1; \
		fi; \
	fi
	@echo "$(COLOR_BLUE)Waiting for database to be ready...$(COLOR_RESET)"
	@until docker exec nr1-postgres_version2 pg_isready -U $$(grep DB_USER .env | cut -d'=' -f2) -d $$(grep DB_NAME .env | cut -d'=' -f2) > /dev/null 2>&1; do \
		echo "$(COLOR_YELLOW)⏳ Waiting for database...$(COLOR_RESET)"; \
		sleep 2; \
	done
	@echo "$(COLOR_BLUE)Running SQL migrations...$(COLOR_RESET)"
	docker exec -i nr1-postgres_version2 psql -U $$(grep DB_USER .env | cut -d'=' -f2) -d $$(grep DB_NAME .env | cut -d'=' -f2) < scripts/init.sql
	@echo "$(COLOR_GREEN)✅ Database migrations completed!$(COLOR_RESET)"

db-reset: db-clean db-start db-migrate
	@echo "$(COLOR_GREEN)✅ Database reset completed!$(COLOR_RESET)"

build:
	@echo "$(COLOR_YELLOW)Building application...$(COLOR_RESET)"
	go build -o bin/nr1-backend $(MAIN_PATH)
	@echo "$(COLOR_GREEN)✅ Build completed! Binary: bin/nr1-backend$(COLOR_RESET)"


dev: db-start db-migrate run
	@echo "$(COLOR_GREEN)✅ Development environment ready!$(COLOR_RESET)"