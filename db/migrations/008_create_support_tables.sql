-- 008_create_user_preferences.sql
-- Future personalization layer (schema-ready, not used in MVP)

CREATE TABLE user_preferences (
    user_id         UUID   PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    preferred_sector TEXT,
    growth_bias      DOUBLE PRECISION NOT NULL DEFAULT 0,
    market_cap_bias  DOUBLE PRECISION NOT NULL DEFAULT 0,
    updated_at       TIMESTAMP        NOT NULL DEFAULT NOW()
);

-- 009_create_refresh_metadata.sql
-- Tracks last successful data refresh per job type
-- Keys: 'price_refresh', 'fundamentals_refresh', 'history_refresh'

CREATE TABLE refresh_metadata (
    key      TEXT      PRIMARY KEY,
    last_run TIMESTAMP NOT NULL
);
