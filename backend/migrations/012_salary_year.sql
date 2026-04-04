-- 012_salary_year.sql
-- When was this salary earned? Defaults to current year. Can go back to 2000.

ALTER TABLE salary_submissions
    ADD COLUMN salary_year SMALLINT NOT NULL DEFAULT EXTRACT(YEAR FROM NOW());

ALTER TABLE salary_submissions
    ADD CONSTRAINT salary_submissions_salary_year_check
    CHECK (salary_year >= 2000 AND salary_year <= 2100);
