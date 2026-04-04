-- 009_merge_tokens.sql
-- Promote auth_tokens to be long-lived sessions: add owner_id column.
-- On verify, owner_id is set and expiry extended to 30 days — same token is reused.

ALTER TABLE auth_tokens ADD COLUMN IF NOT EXISTS owner_id UUID;
