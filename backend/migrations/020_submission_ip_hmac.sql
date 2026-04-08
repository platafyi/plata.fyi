ALTER TABLE salary_submissions
  ADD COLUMN submitter_ip_hmac TEXT;

CREATE INDEX idx_submissions_ip_hmac ON salary_submissions (submitter_ip_hmac);