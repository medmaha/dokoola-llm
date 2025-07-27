.PHONY: build dev prod stop clean logs shell test

# Variables
IMAGE_NAME = llm-api
CONTAINER_NAME = llm-api
PROD_CONTAINER_NAME = llm-api-prod

# Build the Docker image
build:
	docker compose build

# Run in development mode
dev:
	docker compose up --build

# Run in development mode (detached)
dev-d:
	docker compose up --build -d

# Run in production mode
prod:
	docker compose -f docker-compose.prod.yml up --build -d

# Stop all containers
stop:
	docker compose down
	docker compose -f docker-compose.prod.yml down

# Stop and remove all containers, networks, and volumes
clean:
	docker compose down -v --remove-orphans
	docker compose -f docker-compose.prod.yml down -v --remove-orphans
	docker system prune -f

# View logs
logs:
	docker compose logs -f

# View production logs
logs-prod:
	docker compose -f docker-compose.prod.yml logs -f

# Get a shell in the running container
shell:
	docker exec -it $(CONTAINER_NAME) /bin/bash

# Get a shell in the production container
shell-prod:
	docker exec -it $(PROD_CONTAINER_NAME) /bin/bash

# Run tests in container
test:
	docker compose exec llm-api uv run pytest

# Check container health
health:
	docker inspect --format='{{.State.Health.Status}}' $(CONTAINER_NAME)

health-prod:
	docker inspect --format='{{.State.Health.Status}}' $(PROD_CONTAINER_NAME)

# Restart the service
restart:
	docker compose restart

# Restart production service
restart-prod:
	docker compose -f docker-compose.prod.yml restart
