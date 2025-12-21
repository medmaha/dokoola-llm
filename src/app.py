import os

from dotenv import load_dotenv
from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
from pydantic import ValidationError

load_dotenv()

from src.logger import Logger
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

API_BASE_PATH = os.environ.get("API_BASE_PATH", "/api")
app = FastAPI(root_path=API_BASE_PATH, openapi_version="3.0.1")

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


app.include_router(jobs.router)
app.include_router(actions.router)
app.include_router(healthcheck.router)
app.include_router(llm_text_completion.router)
