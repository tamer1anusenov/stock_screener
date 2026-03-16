-- 002_create_stocks.sql
-- Core stock entity: semi-static fundamentals

CREATE TABLE stocks (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticker         VARCHAR(10)       NOT NULL UNIQUE,
    company_name   TEXT              NOT NULL,
    sector         TEXT              NOT NULL,
    industry       TEXT,

    market_cap     BIGINT,
    pe_ratio       NUMERIC(10, 4),
    eps            NUMERIC(10, 4),
    revenue_growth NUMERIC(6, 4),

    ranking_score  DOUBLE PRECISION  NOT NULL DEFAULT 0,

    created_at     TIMESTAMP         NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMP         NOT NULL DEFAULT NOW()
);

-- Used in discover query ordering
CREATE INDEX idx_stocks_ranking   ON stocks (ranking_score DESC);
CREATE INDEX idx_stocks_sector    ON stocks (sector);
CREATE INDEX idx_stocks_marketcap ON stocks (market_cap);