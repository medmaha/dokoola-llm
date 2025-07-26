import os
from fastapi import Request, Response
from fastapi.responses import JSONResponse

from src.constant import (
    ALLOWED_SERVICES,
    DOKOOLA_X_LLM_SERVICE_KEY_NAME,
    DOKOOLA_X_LLM_SERVICE_CLIENT_NAME,
    DOKOOLA_X_LLM_SERVICE_SECRET_HASH_NAME
)
from src.logger import Logger
logger = Logger(__name__)

async def authorization_middleware(request: Request, call_next):
    headers = request.headers
    logger.info(f"Processing authorization request from {headers.get('host', 'unknown host')}")

    llm_service_key = headers.get(DOKOOLA_X_LLM_SERVICE_KEY_NAME)
    llm_service_client_name = headers.get(DOKOOLA_X_LLM_SERVICE_CLIENT_NAME)
    llm_service_secret_hash = headers.get(DOKOOLA_X_LLM_SERVICE_SECRET_HASH_NAME)

    logger.debug(f"Received service key: {llm_service_key}")
    service = ALLOWED_SERVICES.get(llm_service_key)

    if not service:
        logger.warning(f"Invalid service key attempted: {llm_service_key}")
        return JSONResponse(({
            "message":"403: Forbidden request!",
            "reason":"Invalid llm service key provided"
        }), status_code=403)

    # Check if provided has is the same
    if llm_service_secret_hash != service.get("secret_hash"):
        logger.warning(f"Invalid secret hash attempted for service: {llm_service_key}")
        return JSONResponse(({
            "message":"403: Forbidden request!",
            "reason":"Invalid secret hash provided"
        }), status_code=403)
    
    if llm_service_client_name != service.get("client_name"):
        logger.warning(f"Invalid client name attempted: {llm_service_client_name}")
        return JSONResponse(({
            "message":"403: Forbidden request!",
            "reason":"Invalid service-client provided"
        }), status_code=403)

    service_host = service.get("host")
    request_host = headers.get("host", "").split(":")[0]  # Extract host without port

    if service_host != request_host:
        logger.warning(f"Invalid host origin attempted: {request_host}, expected: {service_host}")
        return JSONResponse({
            "message": "403: Forbidden request!",
            "reason": "Invalid host origin"
        }, status_code=403)

    logger.info(f"Authorization successful for service: {llm_service_key}")
    response: Response = await call_next(request)
    return response