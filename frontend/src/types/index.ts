export interface Industry {
  id: number;
  slug: string;
  name: string;
}

export interface City {
  id: number;
  slug: string;
  name: string;
}

export interface Bonus {
  id: string;
  submission_id: string;
  bonus_type: string;
  amount: number;
  frequency: string;
  created_at: string;
}

export interface SalarySubmission {
  id: string;
  company_name: string;
  company_reg_no?: string;
  job_title: string;
  industry_id: number;
  industry_name: string;
  industry_slug: string;
  city_id: number;
  city_name: string;
  city_slug: string;
  seniority: string;
  years_at_company: number;
  years_experience: number;
  work_arrangement: string;
  employment_type: string;
  company_type?: string;
  hours_per_week?: number;
  base_salary: number;
  salary_year: number;
  is_approved: boolean;
  created_at: string;
  updated_at: string;
  bonuses: Bonus[];
}

export interface SalaryStats {
  count: number;
  average: number;
  median: number;
  min: number;
  max: number;
  group_key: string;
  group_val: string;
}

export interface SearchFilters {
  industry?: string;
  city?: string;
  seniority?: string;
  arrangement?: string;
  company_type?: string;
  min_salary?: string;
  max_salary?: string;
  page?: string;
  page_size?: string;
}

export interface SalaryListResponse {
  data: SalarySubmission[];
  total: number;
  page: number;
  page_size: number;
}

export interface BonusInput {
  bonus_type: string;
  amount: number;
  frequency: string;
}

export interface SubmissionInput {
  company_name: string;
  company_reg_no?: string;
  job_title: string;
  industry_id: number;
  city_id: number;
  seniority: string;
  years_at_company: number;
  years_experience: number;
  work_arrangement: string;
  employment_type: string;
  company_type: string;
  hours_per_week?: number;
  base_salary: number;
  salary_year: number;
  bonuses: BonusInput[];
}

export interface Company {
  name: string;
  reg_no?: string;
}

export interface CompaniesResponse {
  results: Company[];
  suggestion?: string;
}
