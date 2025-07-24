import os
import uuid
import logging

ERROR = logging.ERROR
WARNING = logging.WARNING
INFO = logging.INFO
DEBUG = logging.DEBUG

LOG_LEVEL = os.environ.get('API_LOGGING_LEVEL', logging.INFO)

class CorrelationIdFilter(logging.Filter):
    def filter(self, record):
        record.correlation_id = str(uuid.uuid4()).split("-")[1]
        record.api_request_id = str(uuid.uuid4()).split("-")[1]
        return True

class Logger():
    def __init__(self, logger_name: str):
        self.logger = logging.getLogger(logger_name)
        self.logger.addFilter(CorrelationIdFilter())

    def log(
        self,
        level: int,
        message: str,
        *args,
        **kwargs
    ):
        self.logger.log(level, message, *args, **kwargs)

    def exception(self, message: str):
        self.logger.exception(message)

    def error(self, message: str, *args, **kwargs):
        self.logger.error(message, *args, **kwargs)

    def warning(self, message: str, *args, **kwargs):
        self.logger.warning(message, *args, **kwargs)

    def info(self, message: str, *args, **kwargs):
        self.logger.info(message, *args, **kwargs)

    def debug(self, message: str, *args, **kwargs):
        self.logger.debug(message, *args, **kwargs)
 

logger = Logger(__name__)
logger.info(f'dokoola-api-llm logging level is {LOG_LEVEL}')