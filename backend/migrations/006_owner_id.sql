-- 006_owner_id.sql
-- Replace email_hash/session_hash with random owner_id UUID.
-- Email is never stored. Each login session gets a new owner_id.

-- auth_tokens: drop email_hash (no longer needed)
ALTER TABLE auth_tokens DROP COLUMN IF EXISTS email_hash;

-- sessions: drop email_hash, add owner_id UUID generated on verify
ALTER TABLE sessions DROP COLUMN IF EXISTS email_hash;
ALTER TABLE sessions ADD COLUMN IF NOT EXISTS owner_id UUID NOT NULL DEFAULT gen_random_uuid();

-- salary_submissions: replace session_hash with owner_id UUID
ALTER TABLE salary_submissions DROP COLUMN IF EXISTS session_hash;
ALTER TABLE salary_submissions ADD COLUMN IF NOT EXISTS owner_id UUID;

-- index
DROP INDEX IF EXISTS idx_submissions_session_hash;
CREATE INDEX IF NOT EXISTS idx_submissions_owner_id ON salary_submissions(owner_id);
CREATE INDEX IF NOT EXISTS idx_sessions_owner_id ON sessions(owner_id);
