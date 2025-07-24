


from fastapi import APIRouter
from src.logger import Logger
from src.routers.models.healthcheck import HealthCheckResponse


logger = Logger(__name__)
router = APIRouter(
    tags=["Health Check"]
)

@router.get(
    "/health/", 
    response_model=HealthCheckResponse,
    description="Checks the health of the API.",
    responses={
        200: {"detail": "Health check passed"},
        500: {"detail": "Internal server error"},
    },
)
def healthcheck():
    try:
        logger.info("healthcheck request received")
        return HealthCheckResponse(status="OK", message="API is running")
    except Exception as e:
        logger.error(f"healthcheck request failed: {e}")
        return HealthCheckResponse(status="ERROR", message="API is not running")
    finally:
        logger.info("healthcheck request completed")
    
