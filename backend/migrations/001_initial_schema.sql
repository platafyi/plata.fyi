-- 001_initial_schema.sql

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Lookup tables
CREATE TABLE IF NOT EXISTS industries (
    id   SERIAL PRIMARY KEY,
    slug VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS cities (
    id   SERIAL PRIMARY KEY,
    slug VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL
);

-- Magic link tokens (15-minute TTL)
CREATE TABLE IF NOT EXISTS auth_tokens (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email_hash VARCHAR(64) NOT NULL,
    token      VARCHAR(64) UNIQUE NOT NULL,
    used       BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Sessions (30-day TTL)
CREATE TABLE IF NOT EXISTS sessions (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email_hash VARCHAR(64) NOT NULL,
    token      VARCHAR(64) UNIQUE NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Core salary data — no email field ever
CREATE TABLE IF NOT EXISTS salary_submissions (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_hash      VARCHAR(64) NOT NULL,
    company_name      VARCHAR(200) NOT NULL,
    company_reg_no    VARCHAR(50),
    job_title         VARCHAR(200) NOT NULL,
    industry_id       INTEGER NOT NULL REFERENCES industries(id),
    city_id           INTEGER NOT NULL REFERENCES cities(id),
    seniority         VARCHAR(50) NOT NULL,
    years_at_company  SMALLINT NOT NULL CHECK (years_at_company >= 0),
    years_experience  SMALLINT NOT NULL CHECK (years_experience >= 0),
    work_arrangement  VARCHAR(20) NOT NULL,
    base_salary       INTEGER NOT NULL CHECK (base_salary > 0),
    is_approved       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Bonuses (one-to-many per submission)
CREATE TABLE IF NOT EXISTS bonuses (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id UUID NOT NULL REFERENCES salary_submissions(id) ON DELETE CASCADE,
    bonus_type    VARCHAR(50) NOT NULL,
    amount        INTEGER NOT NULL CHECK (amount > 0),
    frequency     VARCHAR(20) NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
