import psycopg2
from psycopg2.extras import RealDictCursor
from contextlib import contextmanager
import config

def get_connection():
    return psycopg2.connect(
        host=config.DB_HOST,
        port=config.DB_PORT,
        user=config.DB_USER,
        password=config.DB_PASSWORD,
        dbname=config.DB_NAME,
        sslmode=config.DB_SSLMODE
    )

@contextmanager
def get_cursor():
    conn = get_connection()
    try:
        with conn.cursor(cursor_factory=RealDictCursor) as cursor:
            yield cursor
            conn.commit()
    except Exception as e:
        conn.rollback()
        raise e
    finally:
        conn.close()

def init_db():
    with get_cursor() as cursor:
        cursor.execute("""
            CREATE TABLE IF NOT EXISTS sync_status (
                id SERIAL PRIMARY KEY,
                sync_type VARCHAR(50) UNIQUE NOT NULL,
                last_sync TIMESTAMP,
                status VARCHAR(20) DEFAULT 'idle',
                records_synced INTEGER DEFAULT 0,
                error_message TEXT
            )
        """)
        
        cursor.execute("""
            INSERT INTO sync_status (sync_type, status) 
            VALUES ('quotes', 'idle'), ('history', 'idle'), ('fundamentals', 'idle')
            ON CONFLICT (sync_type) DO NOTHING
        """)
