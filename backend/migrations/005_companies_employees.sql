-- 005_companies_employees.sql
ALTER TABLE companies ADD COLUMN IF NOT EXISTS employees INTEGER;
