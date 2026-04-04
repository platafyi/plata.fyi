-- 002_indexes.sql

-- Auth lookups
CREATE INDEX IF NOT EXISTS idx_auth_tokens_token ON auth_tokens(token);
CREATE INDEX IF NOT EXISTS idx_auth_tokens_email_hash ON auth_tokens(email_hash);
CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions(token);
CREATE INDEX IF NOT EXISTS idx_sessions_email_hash ON sessions(email_hash);

-- Salary search filters
CREATE INDEX IF NOT EXISTS idx_submissions_industry ON salary_submissions(industry_id) WHERE is_approved = TRUE;
CREATE INDEX IF NOT EXISTS idx_submissions_city ON salary_submissions(city_id) WHERE is_approved = TRUE;
CREATE INDEX IF NOT EXISTS idx_submissions_seniority ON salary_submissions(seniority) WHERE is_approved = TRUE;
CREATE INDEX IF NOT EXISTS idx_submissions_arrangement ON salary_submissions(work_arrangement) WHERE is_approved = TRUE;
CREATE INDEX IF NOT EXISTS idx_submissions_salary ON salary_submissions(base_salary) WHERE is_approved = TRUE;
CREATE INDEX IF NOT EXISTS idx_submissions_created ON salary_submissions(created_at DESC) WHERE is_approved = TRUE;

-- Ownership lookups
CREATE INDEX IF NOT EXISTS idx_submissions_session_hash ON salary_submissions(session_hash);

-- Bonuses
CREATE INDEX IF NOT EXISTS idx_bonuses_submission ON bonuses(submission_id);
