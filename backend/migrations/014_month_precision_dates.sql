ALTER TABLE salary_submissions
  ALTER COLUMN created_at TYPE DATE USING DATE_TRUNC('month', created_at)::DATE,
  ALTER COLUMN updated_at TYPE DATE USING DATE_TRUNC('month', updated_at)::DATE;

ALTER TABLE salary_submissions
  ALTER COLUMN created_at SET DEFAULT DATE_TRUNC('month', NOW())::DATE,
  ALTER COLUMN updated_at SET DEFAULT DATE_TRUNC('month', NOW())::DATE;
