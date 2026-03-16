import yfinance as yf
import pandas as pd
import db
import config
from datetime import datetime

def fetch_quotes(tickers):
    if not tickers:
        return []
    
    data = yf.download(tickers, group_by='ticker', threads=True, progress=False)
    
    results = []
    for ticker in tickers:
        try:
            if len(tickers) == 1:
                ticker_data = data
            else:
                ticker_data = data[ticker]
            
            if ticker_data.empty:
                continue
                
            latest = ticker_data.iloc[-1]
            
            results.append({
                'ticker': ticker,
                'price': float(latest.get('Close', 0)) if pd.notna(latest.get('Close')) else 0,
                'change_percent': float(latest.get('Close', 0) - latest.get('Open', 0)) / latest.get('Open', 1) * 100 if pd.notna(latest.get('Close')) and pd.notna(latest.get('Open')) else 0,
                'volume': int(latest.get('Volume', 0)) if pd.notna(latest.get('Volume')) else 0,
                'day_low': float(latest.get('Low', 0)) if pd.notna(latest.get('Low')) else 0,
                'day_high': float(latest.get('High', 0)) if pd.notna(latest.get('High')) else 0,
                'year_high': float(ticker_data.get('High').max()) if 'High' in ticker_data.columns and not ticker_data.empty else 0,
                'year_low': float(ticker_data.get('Low').min()) if 'Low' in ticker_data.columns and not ticker_data.empty else 0,
            })
        except Exception as e:
            print(f"Error fetching quote for {ticker}: {e}")
            continue
    
    return results

def update_sync_status(sync_type, status, records=0, error=None):
    with db.get_cursor() as cursor:
        cursor.execute("""
            UPDATE sync_status 
            SET status = %s, last_sync = NOW(), records_synced = %s, error_message = %s
            WHERE sync_type = %s
        """, (status, records, error, sync_type))

def bulk_upsert_prices(quotes):
    if not quotes:
        return 0
        
    with db.get_cursor() as cursor:
        for q in quotes:
            cursor.execute("""
                INSERT INTO stock_prices 
                (stock_id, price, change_percent, volume, day_low, day_high, year_high, year_low, updated_at)
                SELECT 
                    id, %s, %s, %s, %s, %s, %s, %s, NOW()
                FROM stocks 
                WHERE ticker = %s
                ON CONFLICT (stock_id) DO UPDATE SET
                    price = EXCLUDED.price,
                    change_percent = EXCLUDED.change_percent,
                    volume = EXCLUDED.volume,
                    day_low = EXCLUDED.day_low,
                    day_high = EXCLUDED.day_high,
                    year_high = EXCLUDED.year_high,
                    year_low = EXCLUDED.year_low,
                    updated_at = NOW()
            """, (q['price'], q['change_percent'], q['volume'], q['day_low'], 
                  q['day_high'], q['year_high'], q['year_low'], q['ticker']))
    
    return len(quotes)

def sync_quotes():
    print("Starting quote sync...")
    update_sync_status('quotes', 'running')
    
    try:
        with db.get_cursor() as cursor:
            cursor.execute("SELECT ticker FROM stocks")
            tickers = [row['ticker'] for row in cursor.fetchall()]
        
        if not tickers:
            print("No stocks in database, skipping quote sync")
            update_sync_status('quotes', 'idle')
            return
        
        all_quotes = []
        for i in range(0, len(tickers), config.BATCH_SIZE):
            batch = tickers[i:i + config.BATCH_SIZE]
            quotes = fetch_quotes(batch)
            all_quotes.extend(quotes)
            print(f"Fetched {len(quotes)} quotes from batch {i//config.BATCH_SIZE + 1}")
        
        count = bulk_upsert_prices(all_quotes)
        print(f"Synced {count} quotes")
        update_sync_status('quotes', 'completed', count)
        
    except Exception as e:
        print(f"Quote sync error: {e}")
        update_sync_status('quotes', 'failed', error=str(e))

if __name__ == "__main__":
    db.init_db()
    sync_quotes()
