DROP TRIGGER IF EXISTS salary_submissions_sort_key_trigger ON salary_submissions;
DROP FUNCTION IF EXISTS set_sort_key();
DROP SEQUENCE IF EXISTS salary_submissions_sort_key_seq;
ALTER TABLE salary_submissions DROP COLUMN IF EXISTS sort_key;
