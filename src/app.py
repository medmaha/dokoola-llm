import os

from dotenv import load_dotenv
from fastapi import FastAPI, Request, APIRouter
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
from pydantic import ValidationError

load_dotenv()

from src.logger import Logger
from src.config import settings
from src.routers import (
    healthcheck,
    jobs,
    actions,
    llm_text_completion,
)

from .middlewares import authorization_middleware, process_timer_middleware

# from src.constant import (
#     ALLOWED_ORIGNS,
#     DOKOOLA_X_LLM_SERVICE_KEY_NAME,
#     DOKOOLA_X_LLM_SERVICE_CLIENT_NAME,
#     DOKOOLA_X_LLM_SERVICE_SECRET_HASH_NAME
# )


logger = Logger(__name__)

app = FastAPI(
    title=settings.app_name,
    debug=settings.debug,
)


app.add_middleware(
    CORSMiddleware,
    allow_credentials=True,
    # allow_origins=ALLOWED_ORIGNS,
    # allow_methods=["get", "post"],
    # allow_headers=[
    #     DOKOOLA_X_LLM_SERVICE_KEY_NAME,
    #     DOKOOLA_X_LLM_SERVICE_CLIENT_NAME,
    #     DOKOOLA_X_LLM_SERVICE_SECRET_HASH_NAME
    # ],
    allow_origins="*",
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.middleware("http")
async def log_request(request: Request, call_next):
    return await process_timer_middleware(request, call_next)


@app.middleware("http")
async def auth(request: Request, call_next):
    return await authorization_middleware(request, call_next)


@app.exception_handler(ValidationError)
async def validation_error_handler(request: Request, exc: ValidationError):
    logger.exception("ERROR")
    return JSONResponse(status_code=400, content={"detail": exc.errors()})


print(settings.api_prefix)

router = APIRouter(prefix=settings.api_prefix, tags=["main"])


@router.get("/health")
async def health_check():
    """Health check endpoint."""
    return {
        "status": "healthy",
        "app": settings.app_name,
        "version": settings.app_version,
    }


router.include_router(jobs.router)
router.include_router(actions.router)
router.include_router(healthcheck.router)
router.include_router(llm_text_completion.router)


app.include_router(router)


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        "src.app:app", host=settings.host, port=settings.port, reload=settings.debug
    )
