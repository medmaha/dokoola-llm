.PHONY: build build-prod rebuild rebuild-prod dev prod stop clean logs shell test

# Variables
IMAGE_NAME = llm-api
CONTAINER_NAME = llm-api
PROD_CONTAINER_NAME = llm-api-prod

# Build the Docker image
api:
	uv run fastapi dev --host 0.0.0.0 --port 8080

build:
	docker compose build

build-prod:
	docker compose -f docker-compose.prod.yml build

# Run in development mode
dev:
	docker compose up --build

# Run in development mode (detached)
dev-d:
	docker compose up --build -d

# Run in production mode
prod:
	docker compose -f docker-compose.prod.yml up --build -d

# Rebuild and redeploy development containers
rebuild:
	@echo "🔄 Rebuilding and redeploying development containers..."
	@echo "Step 1: Building new image..."
	docker compose build --no-cache
	@echo "Step 2: Gracefully stopping old containers..."
	docker compose down --remove-orphans
	@echo "Step 3: Starting new containers..."
	docker compose up -d
	@echo "Step 4: Waiting for health check..."
	@for i in $$(seq 1 30); do \
		sleep 2; \
		STATUS=$$(docker inspect --format='{{.State.Health.Status}}' $(CONTAINER_NAME) 2>/dev/null || echo "starting"); \
		echo "Health check attempt $$i/30: $$STATUS"; \
		if [ "$$STATUS" = "healthy" ]; then \
			echo "✅ Rebuild completed successfully! Container is healthy."; \
			break; \
		fi; \
		if [ $$i -eq 30 ]; then \
			echo "⚠️ Rebuild completed but health check failed. Check logs with: make logs"; \
		fi; \
	done

# Rebuild and redeploy production containers
rebuild-prod:
	@echo "🔄 Rebuilding and redeploying production containers..."
	@echo "Step 1: Building new image..."
	docker compose -f docker-compose.prod.yml build --no-cache
	@echo "Step 2: Gracefully stopping old containers..."
	docker compose -f docker-compose.prod.yml down --remove-orphans
	@echo "Step 3: Starting new containers..."
	docker compose -f docker-compose.prod.yml up -d
	@echo "Step 4: Waiting for health check..."
	@for i in $$(seq 1 30); do \
		sleep 2; \
		STATUS=$$(docker inspect --format='{{.State.Health.Status}}' $(PROD_CONTAINER_NAME) 2>/dev/null || echo "starting"); \
		echo "Health check attempt $$i/30: $$STATUS"; \
		if [ "$$STATUS" = "healthy" ]; then \
			echo "✅ Production rebuild completed successfully! Container is healthy."; \
			break; \
		fi; \
		if [ $$i -eq 30 ]; then \
			echo "⚠️ Production rebuild completed but health check failed. Check logs with: make logs-prod"; \
		fi; \
	done

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
