package database

import "time"


type Industry struct {
	ID   int    `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type City struct {
	ID   int    `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type SalarySubmission struct {
	ID              string  `json:"id"`
	OwnerID         string  `json:"-"`
	CompanyName     string  `json:"company_name"`
	CompanyRegNo    *string `json:"company_reg_no,omitempty"`
	JobTitle        string         `json:"job_title"`
	IndustryID      int            `json:"industry_id"`
	IndustryName    string         `json:"industry_name,omitempty"`
	IndustrySlug    string         `json:"industry_slug,omitempty"`
	CityID          int            `json:"city_id"`
	CityName        string         `json:"city_name,omitempty"`
	CitySlug        string         `json:"city_slug,omitempty"`
	Seniority       string         `json:"seniority"`
	YearsAtCompany  int            `json:"years_at_company"`
	YearsExperience int            `json:"years_experience"`
	WorkArrangement string         `json:"work_arrangement"`
	EmploymentType  string         `json:"employment_type"`
	HoursPerWeek    *int           `json:"hours_per_week,omitempty"`
	BaseSalary      int            `json:"base_salary"`
	SalaryYear      int            `json:"salary_year"`
	IsApproved      bool           `json:"is_approved"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	Bonuses         []Bonus        `json:"bonuses,omitempty"`
}

type Bonus struct {
	ID           string    `json:"id"`
	SubmissionID string    `json:"submission_id"`
	BonusType    string    `json:"bonus_type"`
	Amount       int       `json:"amount"`
	Frequency    string    `json:"frequency"`
	CreatedAt    time.Time `json:"created_at"`
}

// SalaryStats aggregated statistics
type SalaryStats struct {
	Count    int     `json:"count"`
	Average  float64 `json:"average"`
	Median   float64 `json:"median"`
	Min      int     `json:"min"`
	Max      int     `json:"max"`
	GroupKey string  `json:"group_key,omitempty"`
	GroupVal string  `json:"group_val,omitempty"`
}

type Company struct {
	Name  string `json:"name"`
	RegNo string `json:"reg_no,omitempty"`
}

// SearchFilters for the /api/salaries endpoint
type SearchFilters struct {
	IndustrySlug    string
	CitySlug        string
	Seniority       string
	WorkArrangement string
	MinSalary       int
	MaxSalary       int
	Page            int
	PageSize        int
}
