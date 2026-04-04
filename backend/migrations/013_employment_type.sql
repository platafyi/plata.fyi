ALTER TABLE salary_submissions
  ADD COLUMN employment_type VARCHAR(10) NOT NULL DEFAULT 'full_time',
  ADD COLUMN hours_per_week  SMALLINT NULL;

ALTER TABLE salary_submissions
  ADD CONSTRAINT salary_submissions_employment_type_check
    CHECK (employment_type IN ('full_time', 'part_time'));

ALTER TABLE salary_submissions
  ADD CONSTRAINT salary_submissions_hours_per_week_check
    CHECK (hours_per_week IS NULL OR (hours_per_week >= 1 AND hours_per_week <= 168));
