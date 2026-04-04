-- 004_companies.sql
-- Companies are populated organically from user submissions.

CREATE TABLE IF NOT EXISTS companies (
    id     SERIAL PRIMARY KEY,
    name   VARCHAR(200) NOT NULL UNIQUE,
    reg_no VARCHAR(50)
);

CREATE INDEX IF NOT EXISTS idx_companies_name_trgm ON companies(name);
