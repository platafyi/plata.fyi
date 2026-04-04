CREATE SEQUENCE salary_submissions_sort_key_seq START 1000;

ALTER TABLE salary_submissions
  ADD COLUMN sort_key BIGINT NULL;

-- Backfill existing rows with random sort_keys in random order
UPDATE salary_submissions
SET sort_key = sub.rk
FROM (
  SELECT id, (row_number() OVER (ORDER BY random())) * (floor(random() * 1000 + 1)::BIGINT) AS rk
  FROM salary_submissions
) sub
WHERE salary_submissions.id = sub.id;

-- Advance sequence past the max backfilled value
SELECT setval('salary_submissions_sort_key_seq', COALESCE((SELECT MAX(sort_key) FROM salary_submissions), 1000) + 1000);

CREATE OR REPLACE FUNCTION set_sort_key()
RETURNS TRIGGER AS $$
BEGIN
  PERFORM nextval('salary_submissions_sort_key_seq');
  PERFORM setval(
    'salary_submissions_sort_key_seq',
    lastval() + floor(random() * 1000 + 1)::BIGINT
  );
  NEW.sort_key = lastval();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER salary_submissions_sort_key_trigger
BEFORE INSERT ON salary_submissions
FOR EACH ROW EXECUTE FUNCTION set_sort_key();
