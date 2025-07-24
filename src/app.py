import os

from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
from pydantic import ValidationError

 

# from .errors import (
#     BadRequestException,
#     NotFoundException,
# )

from src.logger import Logger
from src.routers import (
    healthcheck,
    llm_text_completion,
)


API_BASE_PATH = os.environ.get('API_BASE_PATH', '/api')
app = FastAPI(root_path=API_BASE_PATH, openapi_version='3.0.1')

app.add_middleware(
    CORSMiddleware,
    allow_credentials=True,
    allow_origins=['*'],
    allow_methods=['*'],
    allow_headers=['*']
)

logger = Logger(__name__)

@app.middleware('http')
async def log_request(request: Request, call_next):
    # Log information about the incoming request

    logger.debug('Incoming request:', extra={'method': request.method, 'url': request.url})
    logger.debug('Headers:', extra={'headers': request.headers})
    logger.debug('Query parameters:', extra={'query_params': request.query_params})
    logger.debug('Request body:', extra={'body': await request.body()})

    # Call the next middleware or route handler
    response = await call_next(request)
    return response

 

@app.exception_handler(ValidationError)
async def validation_error_handler(request: Request, exc: ValidationError):
    logger.exception('ERROR')
    return JSONResponse(
        status_code=400,
        content={'detail': exc.errors()}
    )

 

# @app.exception_handler(NotFoundException)
# async def not_found_exception_handler(request: Request, exc: NotFoundException):
#     logger.exception('ERROR')
#     return JSONResponse(
#         status_code=404,
#         content={'detail': exc.detail}
#     )

 

# @app.exception_handler(BadRequestException)
# async def bad_request_exception_handler(request: Request, exc: BadRequestException):
#     logger.exception('ERROR')
#     return JSONResponse(
#         status_code=400,
#         content={'detail': exc.detail}
#     )


app.include_router(healthcheck.router)
app.include_router(llm_text_completion.router)
