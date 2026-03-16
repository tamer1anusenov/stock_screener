-- 004_create_stock_history.sql
-- OHLCV time-series for charts
-- Composite PK prevents duplicate candles

CREATE TABLE stock_history (
    stock_id UUID          NOT NULL REFERENCES stocks (id) ON DELETE CASCADE,
    date     DATE          NOT NULL,
    open     NUMERIC(12,4) NOT NULL,
    high     NUMERIC(12,4) NOT NULL,
    low      NUMERIC(12,4) NOT NULL,
    close    NUMERIC(12,4) NOT NULL,
    volume   BIGINT,

    PRIMARY KEY (stock_id, date)
);

-- Fast range scans for chart queries: WHERE stock_id = $1 ORDER BY date DESC
CREATE INDEX idx_history_stock_date ON stock_history (stock_id, date DESC);