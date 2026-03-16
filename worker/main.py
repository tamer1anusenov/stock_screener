import sys
import os
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

import time
import config
import db
from sync.sp500 import SP500_TICKERS
from sync.quotes import sync_quotes
from sync.history import sync_history
from sync.fundamentals import sync_fundamentals

def init_stocks():
    with db.get_cursor() as cursor:
        for ticker in SP500_TICKERS:
            cursor.execute("""
                INSERT INTO stocks (ticker, company_name, sector)
                VALUES (%s, %s, 'Unknown')
                ON CONFLICT (ticker) DO NOTHING
            """, (ticker, ticker))
        
        cursor.execute("SELECT COUNT(*) as cnt FROM stocks")
        result = cursor.fetchone()
        count = result['cnt'] if result else 0
        print(f"Initialized {count} stocks in database")

def run():
    print("Starting Stock Screener Worker...")
    db.init_db()
    
    init_stocks()
    
    history_counter = 0
    
    while True:
        print(f"\n=== Loop {time.strftime('%Y-%m-%d %H:%M:%S')} ===")
        
        sync_quotes()
        
        history_counter += 1
        if history_counter >= 4:
            sync_history()
            sync_fundamentals()
            history_counter = 0
        
        print(f"Sleeping for {config.QUOTE_SYNC_INTERVAL} seconds...")
        time.sleep(config.QUOTE_SYNC_INTERVAL)

if __name__ == "__main__":
    run()
