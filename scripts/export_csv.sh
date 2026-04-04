#!/usr/bin/env bash
# Export anonymous salary data from the production DB as CSV.
# Used for CI/CD
set -euo pipefail

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
OUTPUT="${1:-plata_fyi_${TIMESTAMP}.csv}"

QUERY="
SELECT
  s.company_name,
  s.job_title,
  i.name   AS industry,
  c.name   AS city,
  s.seniority,
  s.work_arrangement,
  s.employment_type,
  s.years_experience,
  s.years_at_company,
  s.base_salary,
  s.salary_year,
  to_char(s.created_at, 'YYYY-MM') AS submitted_month
FROM salary_submissions s
JOIN industries i ON i.id = s.industry_id
JOIN cities     c ON c.id = s.city_id
WHERE s.is_approved = true
ORDER BY s.created_at DESC
"

echo "Exporting to ${OUTPUT}..."

kubectl exec -n platafyi deployment/postgres -- \
  psql -U platafyi platafyi -c "\COPY (${QUERY}) TO STDOUT WITH CSV HEADER" \
  > "${OUTPUT}"

echo "Done: ${OUTPUT} ($(wc -l < "${OUTPUT}") rows)"