-- 005_create_users.sql
-- User accounts (MVP version — password/OAuth added later)

CREATE TABLE users (
    id         UUID      PRIMARY KEY DEFAULT gen_random_uuid(),
    email      TEXT      NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);