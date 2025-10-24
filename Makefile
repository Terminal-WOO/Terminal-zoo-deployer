.PHONY: build run dev stop clean test lint

# Docker image name and tag
IMAGE_NAME = ai-co
IMAGE_TAG = latest

# Docker compose commands
DC = docker compose

# Development commands
dev:
	npm install
	npm run dev

# Docker commands
build:
	$(DC) build

run:
	$(DC) up -d

stop:
	$(DC) down

# Clean up commands
clean:
	$(DC) down -v
	rm -rf node_modules
	rm -rf .nuxt
	rm -rf .output

# Application commands
install:
	npm ci

lint:
	npm run lint

test:
	npm run test

# Production commands
build-prod:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

run-prod:
	docker run -d -p 3000:3000 $(IMAGE_NAME):$(IMAGE_TAG)

# Helper commands
logs:
	$(DC) logs -f

ps:
	$(DC) ps

restart:
	$(DC) restart

# Help command
help:
	@echo "Available commands:"
	@echo "  make dev          - Start development server"
	@echo "  make build        - Build using docker-compose"
	@echo "  make run          - Run using docker-compose"
	@echo "  make stop         - Stop all containers"
	@echo "  make clean        - Clean up containers and build files"
	@echo "  make install      - Install dependencies"
	@echo "  make lint         - Run linter"
	@echo "  make test         - Run tests"
	@echo "  make build-prod   - Build production Docker image"
	@echo "  make run-prod     - Run production Docker container"
	@echo "  make logs         - View container logs"
	@echo "  make ps           - List running containers"
	@echo "  make restart      - Restart containers"