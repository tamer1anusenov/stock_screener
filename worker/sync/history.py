import yfinance as yf
import pandas as pd
import db
import config
from datetime import datetime, timedelta

def fetch_history(ticker, period="1y"):
    try:
        stock = yf.Ticker(ticker)
        hist = stock.history(period=period)
        
        if hist.empty:
            return []
        
        records = []
        for idx, row in hist.iterrows():
            records.append({
                'date': idx.strftime('%Y-%m-%d'),
                'open': float(row['Open']),
                'high': float(row['High']),
                'low': float(row['Low']),
                'close': float(row['Close']),
                'volume': int(row['Volume']),
            })
        
        return records
    except Exception as e:
        print(f"Error fetching history for {ticker}: {e}")
        return []

def insert_history(ticker, history):
    if not history:
        return 0
    
    with db.get_cursor() as cursor:
        cursor.execute("SELECT id FROM stocks WHERE ticker = %s", (ticker,))
        result = cursor.fetchone()
        if not result:
            return 0
        
        stock_id = result['id']
        
        for h in history:
            cursor.execute("""
                INSERT INTO stock_history (stock_id, date, open, high, low, close, volume)
                VALUES (%s, %s, %s, %s, %s, %s, %s)
                ON CONFLICT (stock_id, date) DO NOTHING
            """, (stock_id, h['date'], h['open'], h['high'], h['low'], h['close'], h['volume']))
    
    return len(history)

def update_sync_status(sync_type, status, records=0, error=None):
    with db.get_cursor() as cursor:
        cursor.execute("""
            UPDATE sync_status 
            SET status = %s, last_sync = NOW(), records_synced = %s, error_message = %s
            WHERE sync_type = %s
        """, (status, records, error, sync_type))

def sync_history():
    print("Starting history sync...")
    update_sync_status('history', 'running')
    
    try:
        with db.get_cursor() as cursor:
            cursor.execute("""
                SELECT s.ticker 
                FROM stocks s
                LEFT JOIN stock_history sh ON s.id = sh.stock_id
                GROUP BY s.id, s.ticker
                HAVING COUNT(sh.id) < 50
                LIMIT 50
            """)
            tickers = [row['ticker'] for row in cursor.fetchall()]
        
        if not tickers:
            print("All stocks have sufficient history")
            update_sync_status('history', 'completed', 0)
            return
        
        total = 0
        for ticker in tickers:
            history = fetch_history(ticker)
            count = insert_history(ticker, history)
            total += count
            print(f"Synced {count} history points for {ticker}")
        
        print(f"Total history records synced: {total}")
        update_sync_status('history', 'completed', total)
        
    except Exception as e:
        print(f"History sync error: {e}")
        update_sync_status('history', 'failed', error=str(e))

if __name__ == "__main__":
    db.init_db()
    sync_history()
