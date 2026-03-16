-- 006_create_swipes.sql
-- Records every left/right swipe per user per stock
-- UNIQUE constraint prevents a user from swiping the same stock twice

CREATE TYPE swipe_direction AS ENUM ('left', 'right');

CREATE TABLE swipes (
    id         UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID            NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    stock_id   UUID            NOT NULL REFERENCES stocks (id) ON DELETE CASCADE,
    direction  swipe_direction NOT NULL,
    created_at TIMESTAMP       NOT NULL DEFAULT NOW(),

    UNIQUE (user_id, stock_id)
);

-- Discover query: exclude all stocks this user already swiped
CREATE INDEX idx_swipes_user         ON swipes (user_id);
CREATE INDEX idx_swipes_user_created ON swipes (user_id, created_at DESC);
