from src.app import app
from src.logger import Logger

logger = Logger(__name__)
logger.info("Application Up an running")


__all__ = ("app",)
