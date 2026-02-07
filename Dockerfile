# Multi-stage build for minimal image size
# Stage 1: Build stage
FROM python:3.12-slim AS builder

WORKDIR /app

# Set environment variables for pip
ENV PIP_NO_CACHE_DIR=1 \
    PIP_DISABLE_PIP_VERSION_CHECK=1 \
    PYTHONDONTWRITEBYTECODE=1

# Install build dependencies if needed
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    && rm -rf /var/lib/apt/lists/*

# Copy dependency files
COPY pyproject.toml ./

# Install dependencies to a specific location
RUN pip install --target=/app/dependencies -e .

# Stage 2: Runtime stage
FROM python:3.12-slim

WORKDIR /app

# Set environment variables
ENV PYTHONPATH=/app:/app/dependencies \
    PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1

# Install only runtime dependencies (curl for healthcheck)
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get clean

# Copy installed dependencies from builder
COPY --from=builder /app/dependencies /app/dependencies

# Copy application code
COPY . .

# Create non-root user
RUN useradd --create-home --shell /bin/bash --no-log-init appuser && \
    chown -R appuser:appuser /app
USER appuser

# Expose port
EXPOSE 8000

# Run the application
CMD ["python", "-m", "fastapi", "run", "main.py", "--host", "0.0.0.0", "--port", "8000"]
