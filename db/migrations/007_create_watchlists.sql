-- 007_create_watchlists.sql
-- Stocks a user saved (right-swipe result)
-- Separate from swipes for clean querying

CREATE TABLE watchlists (
    id         UUID      PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    stock_id   UUID      NOT NULL REFERENCES stocks (id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE (user_id, stock_id)
);

-- Fetch a user's full watchlist
CREATE INDEX idx_watchlist_user ON watchlists (user_id);
