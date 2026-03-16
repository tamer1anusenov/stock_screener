-- 003_create_stock_prices.sql
-- Latest price snapshot per stock (updated every 5-15 min by worker)
-- One row per stock — stock_id is both PK and FK

CREATE TABLE stock_prices (
    stock_id       UUID          PRIMARY KEY REFERENCES stocks (id) ON DELETE CASCADE,
    price          NUMERIC(12, 4) NOT NULL,
    change_percent NUMERIC(6, 4),
    volume         BIGINT,
    updated_at     TIMESTAMP     NOT NULL DEFAULT NOW()
);