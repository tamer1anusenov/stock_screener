-- 009_alter_stock_prices.sql
-- Add additional price fields for worker sync

ALTER TABLE stock_prices 
ADD COLUMN IF NOT EXISTS day_low NUMERIC(12, 4),
ADD COLUMN IF NOT EXISTS day_high NUMERIC(12, 4),
ADD COLUMN IF NOT EXISTS year_high NUMERIC(12, 4),
ADD COLUMN IF NOT EXISTS year_low NUMERIC(12, 4),
ADD COLUMN IF NOT EXISTS price_avg_50 NUMERIC(12, 4),
ADD COLUMN IF NOT EXISTS price_avg_200 NUMERIC(12, 4);

-- Add logo_url and description to stocks table
ALTER TABLE stocks 
ADD COLUMN IF NOT EXISTS logo_url TEXT,
ADD COLUMN IF NOT EXISTS description TEXT;

-- Create sync_status table for worker
CREATE TABLE IF NOT EXISTS sync_status (
    id SERIAL PRIMARY KEY,
    sync_type VARCHAR(50) UNIQUE NOT NULL,
    last_sync TIMESTAMP,
    status VARCHAR(20) DEFAULT 'idle',
    records_synced INTEGER DEFAULT 0,
    error_message TEXT
);

-- Initialize sync status rows
INSERT INTO sync_status (sync_type, status) 
VALUES ('quotes', 'idle'), ('history', 'idle'), ('fundamentals', 'idle')
ON CONFLICT (sync_type) DO NOTHING;
