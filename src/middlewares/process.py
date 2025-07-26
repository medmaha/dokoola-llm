import time

from fastapi import Request, Response

from src.logger import Logger

logger = Logger(__name__)

async def process_timer_middleware(request: Request, call_next):
    start_time = time.perf_counter()

    logger.debug('Incoming request:', extra={'method': request.method, 'url': request.url})
    logger.debug('Headers:', extra={'headers': request.headers})
    logger.debug('Query parameters:', extra={'query_params': request.query_params})

    response: Response = await call_next(request)
    process_time = time.perf_counter() - start_time

    url = request.url
    method = request.method
    duration = str(process_time)

    logger.info(f"[{method}] - {duration} | {response.status_code} {url}")
    response.headers["X-Process-Time"] = str(process_time)
    return response

