import os
from dotenv import load_dotenv

load_dotenv()

DB_HOST = os.getenv("DB_HOST", "localhost")
DB_PORT = int(os.getenv("DB_PORT", "5432"))
DB_USER = os.getenv("DB_USER", "postgres")
DB_PASSWORD = os.getenv("DB_PASSWORD", "1234")
DB_NAME = os.getenv("DB_NAME", "stockswipe")
DB_SSLMODE = os.getenv("DB_SSLMODE", "disable")

QUOTE_SYNC_INTERVAL = 900
HISTORY_SYNC_INTERVAL = 86400
FUNDAMENTALS_SYNC_INTERVAL = 43200

BATCH_SIZE = 100
