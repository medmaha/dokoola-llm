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
	$(COMPOSE) build

# Run in production mode
start:
	$(COMPOSE) up --build -d

# Rebuild the service
rebuild:
	@echo "🔄 Rebuilding and redeploying production containers..."
	@echo "Step 1: Building new image..."
	$(COMPOSE) build --no-cache
	@echo "Step 2: Gracefully stopping old containers..."
	$(COMPOSE) down --remove-orphans
	@echo "Step 3: Starting new containers..."
	$(COMPOSE) up -d
	@echo "Step 4: Cleaning up unused Docker resources..."
	sudo docker system prune -f
	@echo "✅ Rebuild and redeploy completed."
	
# Stop all containers
stop:
	$(COMPOSE) down
	$(COMPOSE) down

# Stop and remove all containers, networks, and volumes
clean:
	$(COMPOSE) down -v --remove-orphans
	$(COMPOSE) down -v --remove-orphans
	docker system prune -f

# View logs
logs:
	$(COMPOSE) logs -f

# Get a shell in the container
shell:
	docker exec -it $(PROD_CONTAINER_NAME) /bin/bash

# Service health check
health:
	docker inspect --format='{{.State.Health.Status}}' $(PROD_CONTAINER_NAME)

# Restart service
restart:
	$(COMPOSE) restart
