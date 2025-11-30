.PHONY: api build rebuild dev start stop clean logs shell test health restart

# Variables
IMAGE_NAME = llm-api
CONTAINER_NAME = llm-api
PROD_CONTAINER_NAME = llm-api-prod

COMPOSE = sudo docker-compose -f docker-compose.yaml

# Build the Docker image
api:
	uv run fastapi dev --host 0.0.0.0 --port 8080

# Build docker container
build:
	$(COMPOSE) -f docker-compose.prod.yml build

# Run in production mode
start:
	$(COMPOSE) -f docker-compose.prod.yml up --build -d

# Rebuild the service
rebuild:
	@echo "🔄 Rebuilding and redeploying production containers..."
	@echo "Step 1: Building new image..."
	$(COMPOSE) -f docker-compose.prod.yml build --no-cache
	@echo "Step 2: Gracefully stopping old containers..."
	$(COMPOSE) -f docker-compose.prod.yml down --remove-orphans
	@echo "Step 3: Starting new containers..."
	$(COMPOSE) -f docker-compose.prod.yml up -d
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
	$(COMPOSE) down
	$(COMPOSE) -f docker-compose.prod.yml down

# Stop and remove all containers, networks, and volumes
clean:
	$(COMPOSE) down -v --remove-orphans
	$(COMPOSE) -f docker-compose.prod.yml down -v --remove-orphans
	docker system prune -f

# View logs
logs:
	$(COMPOSE) -f docker-compose.prod.yml logs -f

# Get a shell in the container
shell:
	docker exec -it $(PROD_CONTAINER_NAME) /bin/bash

# Service health check
health:
	docker inspect --format='{{.State.Health.Status}}' $(PROD_CONTAINER_NAME)

# Restart service
restart:
	$(COMPOSE) -f docker-compose.prod.yml restart
