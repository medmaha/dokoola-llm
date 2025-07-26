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

    llm_service_key = headers.get(DOKOOLA_X_LLM_SERVICE_KEY_NAME)
    llm_service_client_name = headers.get(DOKOOLA_X_LLM_SERVICE_CLIENT_NAME)
    llm_service_secret_hash = headers.get(DOKOOLA_X_LLM_SERVICE_SECRET_HASH_NAME)

    service = ALLOWED_SERVICES.get(llm_service_key)

    if not service:
        return JSONResponse(({
            "message":"403: Forbidden request!",
            "reason":"Invalid llm service key provided"
        }), status_code=403)

    # Check if provided has is the same
    if llm_service_secret_hash != service.get("secret_hash"):
        return JSONResponse(({
            "message":"403: Forbidden request!",
            "reason":"Invalid secret hash provided"
        }), status_code=403)
    
    if llm_service_client_name != service.get("client_name"):
        return JSONResponse(({
            "message":"403: Forbidden request!",
            "reason":"Invalid service-client provided"
        }), status_code=403)

    service_host = service.get("host")
    request_host = headers.get("host", "").split(":")[0]  # Extract host without port

    if service_host != request_host:
        return JSONResponse({
            "message": "403: Forbidden request!",
            "reason": "Invalid host origin"
        }, status_code=403)

    response: Response = await call_next(request)
    return response