
from dotenv import load_dotenv
load_dotenv()

from src.app import app
from src.logger import Logger

logger= Logger(__name__)
logger.info("Application Up an running")