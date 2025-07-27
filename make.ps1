param(
    [Parameter(Mandatory=$true)]
    [ValidateSet("build", "build-prod", "dev", "dev-d", "prod", "rebuild", "rebuild-prod", "stop", "clean", "logs", "logs-prod", "shell", "shell-prod", "test", "health", "health-prod", "restart", "restart-prod")]
    [string]$Command
)

# Variables
$IMAGE_NAME = "llm-api"
$CONTAINER_NAME = "llm-api"
$PROD_CONTAINER_NAME = "llm-api-prod"

switch ($Command) {
    "build" {
        Write-Host "Building Docker image..." -ForegroundColor Green
        docker-compose build
    }
    "build-prod" {
        Write-Host "Building Docker image..." -ForegroundColor Green
        docker-compose -f docker-compose.prod.yml build
    }
    "dev" {
        Write-Host "Running in development mode..." -ForegroundColor Green
        docker-compose up --build
    }
    "dev-d" {
        Write-Host "Running in development mode (detached)..." -ForegroundColor Green
        docker-compose up --build -d
    }
    "prod" {
        Write-Host "Running in production mode..." -ForegroundColor Green
        docker-compose -f docker-compose.prod.yml up --build -d
    }
    "rebuild" {
        Write-Host "Rebuilding and redeploying development containers..." -ForegroundColor Magenta
        Write-Host "Step 1: Building new image..." -ForegroundColor Yellow
        docker-compose build --no-cache
        
        Write-Host "Step 2: Gracefully stopping old containers..." -ForegroundColor Yellow
        docker-compose down --remove-orphans
        
        Write-Host "Step 3: Starting new containers..." -ForegroundColor Yellow
        docker-compose up -d
        
        Write-Host "Step 4: Waiting for health check..." -ForegroundColor Yellow
        $maxAttempts = 30
        $attempt = 0
        do {
            Start-Sleep -Seconds 2
            $attempt++
            $healthStatus = docker inspect --format='{{.State.Health.Status}}' $CONTAINER_NAME 2>$null
            Write-Host "Health check attempt $attempt/30: $healthStatus" -ForegroundColor Cyan
        } while ($healthStatus -ne "healthy" -and $attempt -lt $maxAttempts)
        
        if ($healthStatus -eq "healthy") {
            Write-Host "✅ Rebuild completed successfully! Container is healthy." -ForegroundColor Green
        } else {
            Write-Host "⚠️ Rebuild completed but health check failed. Check logs with: .\make.ps1 logs" -ForegroundColor Red
        }
    }
    "rebuild-prod" {
        Write-Host "Rebuilding and redeploying production containers..." -ForegroundColor Magenta
        Write-Host "Step 1: Building new image..." -ForegroundColor Yellow
        docker-compose -f docker-compose.prod.yml build --no-cache
        
        Write-Host "Step 2: Gracefully stopping old containers..." -ForegroundColor Yellow
        docker-compose -f docker-compose.prod.yml down --remove-orphans
        
        Write-Host "Step 3: Starting new containers..." -ForegroundColor Yellow
        docker-compose -f docker-compose.prod.yml up -d
        
        Write-Host "Step 4: Waiting for health check..." -ForegroundColor Yellow
        $maxAttempts = 30
        $attempt = 0
        do {
            Start-Sleep -Seconds 2
            $attempt++
            $healthStatus = docker inspect --format='{{.State.Health.Status}}' $PROD_CONTAINER_NAME 2>$null
            Write-Host "Health check attempt $attempt/30: $healthStatus" -ForegroundColor Cyan
        } while ($healthStatus -ne "healthy" -and $attempt -lt $maxAttempts)
        
        if ($healthStatus -eq "healthy") {
            Write-Host "✅ Production rebuild completed successfully! Container is healthy." -ForegroundColor Green
        } else {
            Write-Host "⚠️ Production rebuild completed but health check failed. Check logs with: .\make.ps1 logs-prod" -ForegroundColor Red
        }
    }
    "stop" {
        Write-Host "Stopping all containers..." -ForegroundColor Yellow
        docker-compose down
        docker-compose -f docker-compose.prod.yml down
    }
    "clean" {
        Write-Host "Cleaning up containers, networks, and volumes..." -ForegroundColor Red
        docker-compose down -v --remove-orphans
        docker-compose -f docker-compose.prod.yml down -v --remove-orphans
        docker system prune -f
    }
    "logs" {
        Write-Host "Viewing logs..." -ForegroundColor Cyan
        docker-compose logs -f
    }
    "logs-prod" {
        Write-Host "Viewing production logs..." -ForegroundColor Cyan
        docker-compose -f docker-compose.prod.yml logs -f
    }
    "shell" {
        Write-Host "Opening shell in container..." -ForegroundColor Cyan
        docker exec -it $CONTAINER_NAME /bin/bash
    }
    "shell-prod" {
        Write-Host "Opening shell in production container..." -ForegroundColor Cyan
        docker exec -it $PROD_CONTAINER_NAME /bin/bash
    }
    "test" {
        Write-Host "Running tests..." -ForegroundColor Green
        docker-compose exec llm-api uv run pytest
    }
    "health" {
        Write-Host "Checking container health..." -ForegroundColor Blue
        docker inspect --format='{{.State.Health.Status}}' $CONTAINER_NAME
    }
    "health-prod" {
        Write-Host "Checking container health..." -ForegroundColor Blue
        docker inspect --format='{{.State.Health.Status}}' $PROD_CONTAINER_NAME
    }
    "restart" {
        Write-Host "Restarting service..." -ForegroundColor Yellow
        docker-compose restart
    }
    "restart-prod" {
        Write-Host "Restarting production service..." -ForegroundColor Yellow
        docker-compose -f docker-compose.prod.yml restart
    }
}
