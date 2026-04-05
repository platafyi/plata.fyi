package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func nullStringToPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

// GetIndustries returns all industries ordered by name.
func (s *PostgresStore) GetIndustries(ctx context.Context) ([]Industry, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, slug, name FROM industries ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("get industries: %w", err)
	}
	defer rows.Close()

	var out []Industry
	for rows.Next() {
		var i Industry
		if err := rows.Scan(&i.ID, &i.Slug, &i.Name); err != nil {
			return nil, err
		}
		out = append(out, i)
	}
	return out, rows.Err()
}

// GetCities returns all cities except "remote", with Skopje sorted first.
func (s *PostgresStore) GetCities(ctx context.Context) ([]City, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, slug, name FROM cities WHERE slug != 'remote' ORDER BY CASE slug WHEN 'skopje' THEN 0 ELSE 1 END, name`)
	if err != nil {
		return nil, fmt.Errorf("get cities: %w", err)
	}
	defer rows.Close()

	var out []City
	for rows.Next() {
		var c City
		if err := rows.Scan(&c.ID, &c.Slug, &c.Name); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// InsertToken stores a new magic-link token. The DB assigns owner_id and expires_at.
func (s *PostgresStore) InsertToken(ctx context.Context, token string) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO tokens (token) VALUES ($1)`, token)
	return err
}

// GetOwnerByToken returns the owner_id for a valid, non-expired token, or nil if not found.
func (s *PostgresStore) GetOwnerByToken(ctx context.Context, token string) (*string, error) {
	var ownerID string
	err := s.db.QueryRowContext(ctx,
		`SELECT owner_id::text FROM tokens WHERE token = $1 AND expires_at > NOW()`,
		token,
	).Scan(&ownerID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &ownerID, nil
}

// DeleteToken removes a token, used on logout.
func (s *PostgresStore) DeleteToken(ctx context.Context, token string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM tokens WHERE token = $1`, token)
	return err
}

// ExchangeToken atomically validates magicToken, issues a new session token with the
// same owner_id, and deletes the magic token. Returns the owner_id, or nil if the
// magic token is invalid or expired.
func (s *PostgresStore) ExchangeToken(ctx context.Context, magicToken, sessionToken string) (*string, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var ownerID string
	err = tx.QueryRowContext(ctx,
		`SELECT owner_id::text FROM tokens WHERE token = $1 AND expires_at > NOW()`,
		magicToken,
	).Scan(&ownerID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO tokens (token, owner_id) VALUES ($1, $2::uuid)`,
		sessionToken, ownerID,
	)
	if err != nil {
		return nil, fmt.Errorf("insert session token: %w", err)
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM tokens WHERE token = $1`, magicToken)
	if err != nil {
		return nil, fmt.Errorf("delete magic token: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &sessionToken, nil
}

type CreateSubmissionInput struct {
	OwnerID         string
	CompanyName     string
	CompanyRegNo    string
	JobTitle        string
	IndustryID      int
	CityID          int
	Seniority       string
	YearsAtCompany  int
	YearsExperience int
	WorkArrangement string
	EmploymentType  string
	HoursPerWeek    *int
	BaseSalary      int
	SalaryYear      int
	Bonuses         []BonusInput
}

type BonusInput struct {
	BonusType string
	Amount    int
	Frequency string
}

// CreateSubmission inserts a new salary submission and its bonuses in a transaction,
// then best-effort upserts the company name for autocomplete.
func (s *PostgresStore) CreateSubmission(ctx context.Context, inp CreateSubmissionInput) (*SalarySubmission, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var regNo interface{}
	if inp.CompanyRegNo != "" {
		regNo = inp.CompanyRegNo
	}

	var sub SalarySubmission
	employmentType := inp.EmploymentType
	if employmentType == "" {
		employmentType = "full_time"
	}
	err = tx.QueryRowContext(ctx,
		`INSERT INTO salary_submissions
		 (owner_id, company_name, company_reg_no, job_title, industry_id, city_id,
		  seniority, years_at_company, years_experience, work_arrangement,
		  employment_type, hours_per_week, base_salary, salary_year)
		 VALUES ($1::uuid,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		 RETURNING id, created_at, updated_at`,
		inp.OwnerID, inp.CompanyName, regNo, inp.JobTitle,
		inp.IndustryID, inp.CityID, inp.Seniority,
		inp.YearsAtCompany, inp.YearsExperience, inp.WorkArrangement,
		employmentType, inp.HoursPerWeek, inp.BaseSalary, inp.SalaryYear,
	).Scan(&sub.ID, &sub.CreatedAt, &sub.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert submission: %w", err)
	}

	for _, b := range inp.Bonuses {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO bonuses (submission_id, bonus_type, amount, frequency)
			 VALUES ($1,$2,$3,$4)`,
			sub.ID, b.BonusType, b.Amount, b.Frequency,
		)
		if err != nil {
			return nil, fmt.Errorf("insert bonus: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Upsert company for autocomplete (best-effort, ignore errors)
	s.upsertCompany(ctx, inp.CompanyName, inp.CompanyRegNo)

	sub.OwnerID = inp.OwnerID
	sub.CompanyName = inp.CompanyName
	sub.JobTitle = inp.JobTitle
	sub.IndustryID = inp.IndustryID
	sub.CityID = inp.CityID
	sub.Seniority = inp.Seniority
	sub.YearsAtCompany = inp.YearsAtCompany
	sub.YearsExperience = inp.YearsExperience
	sub.WorkArrangement = inp.WorkArrangement
	sub.EmploymentType = employmentType
	sub.HoursPerWeek = inp.HoursPerWeek
	sub.BaseSalary = inp.BaseSalary
	sub.SalaryYear = inp.SalaryYear
	sub.IsApproved = true

	return &sub, nil
}

// GetSubmissionsByOwner returns all submissions for the given owner, with bonuses loaded.
func (s *PostgresStore) GetSubmissionsByOwner(ctx context.Context, ownerID string) ([]SalarySubmission, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT s.id, s.owner_id, s.company_name, s.company_reg_no, s.job_title,
		        s.industry_id, i.name, i.slug, s.city_id, c.name, c.slug,
		        s.seniority, s.years_at_company, s.years_experience,
		        s.work_arrangement, s.employment_type, s.hours_per_week,
		        s.base_salary, s.salary_year, s.is_approved, s.created_at, s.updated_at
		 FROM salary_submissions s
		 JOIN industries i ON i.id = s.industry_id
		 JOIN cities c ON c.id = s.city_id
		 WHERE s.owner_id = $1::uuid
		 ORDER BY s.created_at DESC, s.ctid DESC`,
		ownerID,
	)
	if err != nil {
		return nil, fmt.Errorf("get submissions by owner: %w", err)
	}
	defer rows.Close()

	var out []SalarySubmission
	for rows.Next() {
		var sub SalarySubmission
		var regNo sql.NullString
		err := rows.Scan(
			&sub.ID, &sub.OwnerID, &sub.CompanyName, &regNo, &sub.JobTitle,
			&sub.IndustryID, &sub.IndustryName, &sub.IndustrySlug,
			&sub.CityID, &sub.CityName, &sub.CitySlug,
			&sub.Seniority, &sub.YearsAtCompany, &sub.YearsExperience,
			&sub.WorkArrangement, &sub.EmploymentType, &sub.HoursPerWeek,
			&sub.BaseSalary, &sub.SalaryYear, &sub.IsApproved, &sub.CreatedAt, &sub.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		sub.CompanyRegNo = nullStringToPtr(regNo)
		out = append(out, sub)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range out {
		bonuses, err := s.getBonuses(ctx, out[i].ID)
		if err != nil {
			return nil, err
		}
		out[i].Bonuses = bonuses
	}

	return out, nil
}

// GetSubmissionByID returns a single submission with bonuses, or nil if not found.
func (s *PostgresStore) GetSubmissionByID(ctx context.Context, id string) (*SalarySubmission, error) {
	var sub SalarySubmission
	var regNo sql.NullString
	err := s.db.QueryRowContext(ctx,
		`SELECT s.id, s.owner_id, s.company_name, s.company_reg_no, s.job_title,
		        s.industry_id, i.name, i.slug, s.city_id, c.name, c.slug,
		        s.seniority, s.years_at_company, s.years_experience,
		        s.work_arrangement, s.employment_type, s.hours_per_week,
		        s.base_salary, s.salary_year, s.is_approved, s.created_at, s.updated_at
		 FROM salary_submissions s
		 JOIN industries i ON i.id = s.industry_id
		 JOIN cities c ON c.id = s.city_id
		 WHERE s.id = $1`,
		id,
	).Scan(
		&sub.ID, &sub.OwnerID, &sub.CompanyName, &regNo, &sub.JobTitle,
		&sub.IndustryID, &sub.IndustryName, &sub.IndustrySlug,
		&sub.CityID, &sub.CityName, &sub.CitySlug,
		&sub.Seniority, &sub.YearsAtCompany, &sub.YearsExperience,
		&sub.WorkArrangement, &sub.EmploymentType, &sub.HoursPerWeek,
		&sub.BaseSalary, &sub.SalaryYear, &sub.IsApproved, &sub.CreatedAt, &sub.UpdatedAt,
	)
	sub.CompanyRegNo = nullStringToPtr(regNo)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	bonuses, err := s.getBonuses(ctx, sub.ID)
	if err != nil {
		return nil, err
	}
	sub.Bonuses = bonuses
	return &sub, nil
}

// UpdateSubmission replaces all fields and bonuses for a submission the caller owns.
// Returns sql.ErrNoRows if the id/ownerID pair does not match.
func (s *PostgresStore) UpdateSubmission(ctx context.Context, id, ownerID string, inp CreateSubmissionInput) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var regNo interface{}
	if inp.CompanyRegNo != "" {
		regNo = inp.CompanyRegNo
	}

	employmentType := inp.EmploymentType
	if employmentType == "" {
		employmentType = "full_time"
	}
	res, err := tx.ExecContext(ctx,
		`UPDATE salary_submissions SET
		  company_name=$1, company_reg_no=$2, job_title=$3, industry_id=$4, city_id=$5,
		  seniority=$6, years_at_company=$7, years_experience=$8, work_arrangement=$9,
		  employment_type=$10, hours_per_week=$11, base_salary=$12, salary_year=$13, updated_at=DATE_TRUNC('month', NOW())::DATE
		 WHERE id=$14 AND owner_id=$15::uuid`,
		inp.CompanyName, regNo, inp.JobTitle, inp.IndustryID, inp.CityID,
		inp.Seniority, inp.YearsAtCompany, inp.YearsExperience, inp.WorkArrangement,
		employmentType, inp.HoursPerWeek, inp.BaseSalary, inp.SalaryYear, id, ownerID,
	)
	if err != nil {
		return fmt.Errorf("update submission: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM bonuses WHERE submission_id = $1`, id); err != nil {
		return err
	}
	for _, b := range inp.Bonuses {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO bonuses (submission_id, bonus_type, amount, frequency)
			 VALUES ($1,$2,$3,$4)`,
			id, b.BonusType, b.Amount, b.Frequency,
		)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	s.upsertCompany(ctx, inp.CompanyName, inp.CompanyRegNo)
	return nil
}

// DeleteSubmission removes a submission the caller owns.
// Returns sql.ErrNoRows if the id/ownerID pair does not match.
func (s *PostgresStore) DeleteSubmission(ctx context.Context, id, ownerID string) error {
	res, err := s.db.ExecContext(ctx,
		`DELETE FROM salary_submissions WHERE id = $1 AND owner_id = $2::uuid`,
		id, ownerID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *PostgresStore) upsertCompany(ctx context.Context, name, regNo string) {
	var reg interface{}
	if regNo != "" {
		reg = regNo
	}
	s.db.ExecContext(ctx, `INSERT INTO companies (name, reg_no) VALUES ($1, $2) ON CONFLICT (name) DO NOTHING`, name, reg)
}

func (s *PostgresStore) getBonuses(ctx context.Context, submissionID string) ([]Bonus, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, submission_id, bonus_type, amount, frequency, created_at
		 FROM bonuses WHERE submission_id = $1 ORDER BY created_at`,
		submissionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Bonus
	for rows.Next() {
		var b Bonus
		if err := rows.Scan(&b.ID, &b.SubmissionID, &b.BonusType, &b.Amount, &b.Frequency, &b.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// SearchSalaries returns a paginated, filtered list of approved submissions and the total count.
func (s *PostgresStore) SearchSalaries(ctx context.Context, f SearchFilters) ([]SalarySubmission, int, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 || f.PageSize > 100 {
		f.PageSize = 20
	}

	args := []interface{}{}
	where := []string{"s.is_approved = TRUE"}
	argN := 1

	if f.IndustrySlug != "" {
		args = append(args, f.IndustrySlug)
		where = append(where, fmt.Sprintf("i.slug = $%d", argN))
		argN++
	}
	if f.CitySlug != "" {
		args = append(args, f.CitySlug)
		where = append(where, fmt.Sprintf("c.slug = $%d", argN))
		argN++
	}
	if f.Seniority != "" {
		args = append(args, f.Seniority)
		where = append(where, fmt.Sprintf("s.seniority = $%d", argN))
		argN++
	}
	if f.WorkArrangement != "" {
		args = append(args, f.WorkArrangement)
		where = append(where, fmt.Sprintf("s.work_arrangement = $%d", argN))
		argN++
	}
	if f.MinSalary > 0 {
		args = append(args, f.MinSalary)
		where = append(where, fmt.Sprintf("s.base_salary >= $%d", argN))
		argN++
	}
	if f.MaxSalary > 0 {
		args = append(args, f.MaxSalary)
		where = append(where, fmt.Sprintf("s.base_salary <= $%d", argN))
		argN++
	}

	whereSQL := strings.Join(where, " AND ")

	var total int
	err := s.db.QueryRowContext(ctx,
		fmt.Sprintf(`SELECT COUNT(*) FROM salary_submissions s
		 JOIN industries i ON i.id = s.industry_id
		 JOIN cities c ON c.id = s.city_id
		 WHERE %s`, whereSQL),
		args...,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count salaries: %w", err)
	}

	offset := (f.Page - 1) * f.PageSize
	args = append(args, f.PageSize, offset)

	rows, err := s.db.QueryContext(ctx,
		fmt.Sprintf(`SELECT s.id, s.company_name, s.job_title,
		        s.industry_id, i.name, i.slug, s.city_id, c.name, c.slug,
		        s.seniority, s.years_at_company, s.years_experience,
		        s.work_arrangement, s.base_salary, s.created_at
		 FROM salary_submissions s
		 JOIN industries i ON i.id = s.industry_id
		 JOIN cities c ON c.id = s.city_id
		 WHERE %s
		 ORDER BY s.created_at DESC, s.ctid DESC
		 LIMIT $%d OFFSET $%d`, whereSQL, argN, argN+1),
		args...,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("search salaries: %w", err)
	}
	defer rows.Close()

	var out []SalarySubmission
	for rows.Next() {
		var sub SalarySubmission
		err := rows.Scan(
			&sub.ID, &sub.CompanyName, &sub.JobTitle,
			&sub.IndustryID, &sub.IndustryName, &sub.IndustrySlug,
			&sub.CityID, &sub.CityName, &sub.CitySlug,
			&sub.Seniority, &sub.YearsAtCompany, &sub.YearsExperience,
			&sub.WorkArrangement, &sub.BaseSalary, &sub.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, sub)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

// GetSalaryStats returns aggregated salary statistics grouped by industry or city.
func (s *PostgresStore) GetSalaryStats(ctx context.Context, groupBy string, f SearchFilters) ([]SalaryStats, error) {
	where := []string{"s.is_approved = TRUE"}
	args := []interface{}{}
	argN := 1

	if f.IndustrySlug != "" {
		args = append(args, f.IndustrySlug)
		where = append(where, fmt.Sprintf("i.slug = $%d", argN))
		argN++
	}
	if f.CitySlug != "" {
		args = append(args, f.CitySlug)
		where = append(where, fmt.Sprintf("c.slug = $%d", argN))
		argN++
	}
	if f.Seniority != "" {
		args = append(args, f.Seniority)
		where = append(where, fmt.Sprintf("s.seniority = $%d", argN))
		argN++
	}

	whereSQL := strings.Join(where, " AND ")

	var groupCol, groupLabel string
	switch groupBy {
	case "city":
		groupCol = "c.slug"
		groupLabel = "c.name"
	default:
		groupCol = "i.slug"
		groupLabel = "i.name"
	}

	query := fmt.Sprintf(`
		SELECT %s, %s,
		       COUNT(*),
		       AVG(s.base_salary),
		       PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY s.base_salary),
		       MIN(s.base_salary),
		       MAX(s.base_salary)
		FROM salary_submissions s
		JOIN industries i ON i.id = s.industry_id
		JOIN cities c ON c.id = s.city_id
		WHERE %s
		GROUP BY %s, %s
		ORDER BY COUNT(*) DESC
	`, groupCol, groupLabel, whereSQL, groupCol, groupLabel)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get stats: %w", err)
	}
	defer rows.Close()

	var out []SalaryStats
	for rows.Next() {
		var st SalaryStats
		if err := rows.Scan(&st.GroupKey, &st.GroupVal, &st.Count, &st.Average, &st.Median, &st.Min, &st.Max); err != nil {
			return nil, err
		}
		out = append(out, st)
	}
	return out, rows.Err()
}

// SearchCompanies returns up to 10 companies whose names contain q (case-insensitive).
func (s *PostgresStore) SearchCompanies(ctx context.Context, q string) ([]Company, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT name, COALESCE(reg_no, '') FROM companies
		 WHERE name ILIKE $1
		 ORDER BY name
		 LIMIT 10`,
		"%"+q+"%",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Company
	for rows.Next() {
		var c Company
		if err := rows.Scan(&c.Name, &c.RegNo); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// SearchJobTitles returns up to 10 distinct job title strings matching q, drawn from
// both approved submissions and the curated job_titles table.
func (s *PostgresStore) SearchJobTitles(ctx context.Context, q string) ([]string, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT DISTINCT name FROM (
		    SELECT job_title AS name FROM salary_submissions WHERE job_title ILIKE $1 AND is_approved = TRUE
		    UNION
		    SELECT name FROM job_titles WHERE name ILIKE $1
		 ) t
		 ORDER BY name
		 LIMIT 10`,
		"%"+q+"%",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}
