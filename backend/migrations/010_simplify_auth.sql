-- 010_simplify_auth.sql
-- One token, 30 days. No sessions table. No used flag.
-- auth_tokens IS the session: created with owner_id upfront, valid for 30 days.

DROP TABLE IF EXISTS sessions;

DROP TABLE IF EXISTS auth_tokens;

CREATE TABLE auth_tokens (
    id         SERIAL PRIMARY KEY,
    token      VARCHAR(64) UNIQUE NOT NULL,
    owner_id   UUID NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL
);
