import yfinance as yf
import db

def fetch_fundamentals(ticker):
    try:
        stock = yf.Ticker(ticker)
        info = stock.info
        
        return {
            'ticker': ticker,
            'market_cap': info.get('marketCap'),
            'pe_ratio': info.get('trailingPE'),
            'eps': info.get('trailingEps'),
            'revenue_growth': info.get('revenueGrowth'),
            'sector': info.get('sector'),
            'industry': info.get('industry'),
            'description': info.get('longBusinessSummary', '')[:1000],
            'logo_url': info.get('logo_url'),
            'ceo': info.get('ceo'),
            'website': info.get('website'),
        }
    except Exception as e:
        print(f"Error fetching fundamentals for {ticker}: {e}")
        return None

def update_stock_fundamentals(data):
    if not data:
        return False
    
    with db.get_cursor() as cursor:
        cursor.execute("""
            UPDATE stocks 
            SET sector = %s,
                industry = %s,
                market_cap = %s,
                pe_ratio = %s,
                eps = %s,
                revenue_growth = %s,
                description = %s,
                logo_url = %s,
                updated_at = NOW()
            WHERE ticker = %s
        """, (
            data.get('sector'),
            data.get('industry'),
            data.get('market_cap'),
            data.get('pe_ratio'),
            data.get('eps'),
            data.get('revenue_growth'),
            data.get('description'),
            data.get('logo_url'),
            data.get('ticker')
        ))
    
    return True

def sync_fundamentals():
    print("Starting fundamentals sync...")
    
    with db.get_cursor() as cursor:
        cursor.execute("SELECT ticker FROM stocks ORDER BY updated_at ASC NULLS FIRST LIMIT 100")
        tickers = [row['ticker'] for row in cursor.fetchall()]
    
    if not tickers:
        print("No stocks to sync fundamentals for")
        return
    
    success = 0
    for ticker in tickers:
        data = fetch_fundamentals(ticker)
        if data and update_stock_fundamentals(data):
            success += 1
            print(f"Synced fundamentals for {ticker}")
    
    print(f"Synced {success}/{len(tickers)} fundamentals")

if __name__ == "__main__":
    db.init_db()
    sync_fundamentals()
