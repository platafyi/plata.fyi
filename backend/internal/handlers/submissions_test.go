package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/platafyi/plata.fyi/internal/database"
	"github.com/platafyi/plata.fyi/internal/middleware"
)

const validSubmissionBody = `{"company_name":"Test Co","job_title":"Engineer","industry_id":1,"city_id":1,"seniority":"mid","work_arrangement":"office","employment_type":"full_time","base_salary":50000,"salary_year":2024}`

func makeSubmission(id string) database.SalarySubmission {
	now := time.Now()
	return database.SalarySubmission{
		ID: id, OwnerID: "owner-uuid", CompanyName: "Test Co", JobTitle: "Engineer",
		IndustryID: 1, IndustryName: "Tech", IndustrySlug: "tech",
		CityID: 1, CityName: "Skopje", CitySlug: "skopje",
		Seniority: "mid", YearsAtCompany: 2, YearsExperience: 5,
		WorkArrangement: "office", EmploymentType: "full_time",
		BaseSalary: 50000, SalaryYear: 2024, IsApproved: true,
		CreatedAt: now, UpdatedAt: now,
	}
}

func authStore(ownerID string) *MockStore {
	return &MockStore{OwnerByToken: &ownerID}
}

func TestValidateRequestStringLengths(t *testing.T) {
	h := NewSubmissionsHandler(&MockStore{}, "", "")

	base := func() *submissionRequest {
		return &submissionRequest{
			CompanyName: "Test Co", JobTitle: "Engineer",
			IndustryID: 1, CityID: 1, Seniority: "mid",
			WorkArrangement: "office", EmploymentType: "full_time",
			BaseSalary: 50000, SalaryYear: 2025,
		}
	}

	t.Run("company_name at limit (100 chars)", func(t *testing.T) {
		r := base()
		r.CompanyName = strings.Repeat("а", 100)
		if got := h.validateRequest(r); got != "" {
			t.Errorf("expected no error, got %q", got)
		}
	})
	t.Run("company_name too long (101 chars)", func(t *testing.T) {
		r := base()
		r.CompanyName = strings.Repeat("а", 101)
		if got := h.validateRequest(r); got == "" {
			t.Error("expected error for company_name > 100 chars")
		}
	})
	t.Run("job_title at limit (100 chars)", func(t *testing.T) {
		r := base()
		r.JobTitle = strings.Repeat("а", 100)
		if got := h.validateRequest(r); got != "" {
			t.Errorf("expected no error, got %q", got)
		}
	})
	t.Run("job_title too long (101 chars)", func(t *testing.T) {
		r := base()
		r.JobTitle = strings.Repeat("а", 101)
		if got := h.validateRequest(r); got == "" {
			t.Error("expected error for job_title > 100 chars")
		}
	})
	t.Run("company_reg_no at limit (20 chars)", func(t *testing.T) {
		r := base()
		r.CompanyRegNo = strings.Repeat("1", 20)
		if got := h.validateRequest(r); got != "" {
			t.Errorf("expected no error, got %q", got)
		}
	})
	t.Run("company_reg_no too long (21 chars)", func(t *testing.T) {
		r := base()
		r.CompanyRegNo = strings.Repeat("1", 21)
		if got := h.validateRequest(r); got == "" {
			t.Error("expected error for company_reg_no > 20 chars")
		}
	})
}

func TestValidateRequestYearsBounds(t *testing.T) {
	h := NewSubmissionsHandler(&MockStore{}, "", "")

	base := func() *submissionRequest {
		return &submissionRequest{
			CompanyName: "Test Co", JobTitle: "Engineer",
			IndustryID: 1, CityID: 1, Seniority: "mid",
			WorkArrangement: "office", EmploymentType: "full_time",
			BaseSalary: 50000, SalaryYear: 2025,
		}
	}

	cases := []struct {
		name    string
		modify  func(r *submissionRequest)
		wantErr bool
	}{
		{"years_at_company 0", func(r *submissionRequest) { r.YearsAtCompany = 0 }, false},
		{"years_at_company 60", func(r *submissionRequest) { r.YearsAtCompany = 60 }, false},
		{"years_at_company 61", func(r *submissionRequest) { r.YearsAtCompany = 61 }, true},
		{"years_at_company 9999", func(r *submissionRequest) { r.YearsAtCompany = 9999 }, true},
		{"years_at_company -1", func(r *submissionRequest) { r.YearsAtCompany = -1 }, true},
		{"years_experience 0", func(r *submissionRequest) { r.YearsExperience = 0 }, false},
		{"years_experience 60", func(r *submissionRequest) { r.YearsExperience = 60 }, false},
		{"years_experience 61", func(r *submissionRequest) { r.YearsExperience = 61 }, true},
		{"years_experience 9999", func(r *submissionRequest) { r.YearsExperience = 9999 }, true},
		{"years_experience -1", func(r *submissionRequest) { r.YearsExperience = -1 }, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := base()
			tc.modify(r)
			got := h.validateRequest(r)
			if tc.wantErr && got == "" {
				t.Error("expected error but got none")
			}
			if !tc.wantErr && got != "" {
				t.Errorf("expected no error but got %q", got)
			}
		})
	}
}

func TestValidateRequestBonusCap(t *testing.T) {
	h := NewSubmissionsHandler(&MockStore{}, "", "")

	makeReq := func(n int) *submissionRequest {
		bonuses := make([]bonusRequest, n)
		for i := range bonuses {
			bonuses[i] = bonusRequest{BonusType: "annual", Amount: 1000, Frequency: "annual"}
		}
		return &submissionRequest{
			CompanyName: "Test Co", JobTitle: "Engineer",
			IndustryID: 1, CityID: 1, Seniority: "mid",
			WorkArrangement: "office", EmploymentType: "full_time",
			BaseSalary: 50000, SalaryYear: 2025,
			Bonuses: bonuses,
		}
	}

	cases := []struct {
		name    string
		count   int
		wantErr bool
	}{
		{"0 bonuses", 0, false},
		{"1 bonus", 1, false},
		{"10 bonuses", 10, false},
		{"11 bonuses", 11, true},
		{"100 bonuses", 100, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := h.validateRequest(makeReq(tc.count))
			if tc.wantErr && got == "" {
				t.Error("expected error but got none")
			}
			if !tc.wantErr && got != "" {
				t.Errorf("expected no error but got %q", got)
			}
		})
	}
}

func TestValidateRequestSalaryBounds(t *testing.T) {
	h := NewSubmissionsHandler(&MockStore{}, "", "")

	baseReq := func(employmentType string, salary int) *submissionRequest {
		return &submissionRequest{
			CompanyName:     "Test Co",
			JobTitle:        "Engineer",
			IndustryID:      1,
			CityID:          1,
			Seniority:       "mid",
			WorkArrangement: "office",
			EmploymentType:  employmentType,
			BaseSalary:      salary,
			SalaryYear:      2025,
		}
	}

	cases := []struct {
		name    string
		req     *submissionRequest
		wantErr bool
	}{
		// full_time: minimum wage applies
		{"full_time zero", baseReq("full_time", 0), true},
		{"full_time below min", baseReq("full_time", minSalaryMKD-1), true},
		{"full_time at min", baseReq("full_time", minSalaryMKD), false},
		{"full_time typical", baseReq("full_time", 55_000), false},
		{"full_time at max", baseReq("full_time", maxSalaryMKD), false},
		{"full_time above max", baseReq("full_time", maxSalaryMKD+1), true},
		// part_time: no minimum wage floor, only max wage
		{"part_time zero", baseReq("part_time", 0), true},
		{"part_time below full_time min", baseReq("part_time", minSalaryMKD-1), false},
		{"part_time typical", baseReq("part_time", 15_000), false},
		{"part_time at max", baseReq("part_time", maxSalaryMKD), false},
		{"part_time above max", baseReq("part_time", maxSalaryMKD+1), true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := h.validateRequest(tc.req)
			if tc.wantErr && got == "" {
				t.Errorf("expected validation error but got none")
			}
			if !tc.wantErr && got != "" {
				t.Errorf("expected no error but got %q", got)
			}
		})
	}
}

func TestList(t *testing.T) {
	ownerID := "owner-uuid"

	t.Run("401 no auth", func(t *testing.T) {
		store := &MockStore{}
		h := NewSubmissionsHandler(store, "", "")
		req := httptest.NewRequest(http.MethodGet, "/api/submissions", nil)
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.List)).ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("got %d, want 401", rec.Code)
		}
	})

	t.Run("500 db query fails", func(t *testing.T) {
		store := authStore(ownerID)
		store.SubmissionsErr = errors.New("db error")
		h := NewSubmissionsHandler(store, "", "")
		req := httptest.NewRequest(http.MethodGet, "/api/submissions", nil)
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.List)).ServeHTTP(rec, req)
		if rec.Code != http.StatusInternalServerError {
			t.Errorf("got %d, want 500", rec.Code)
		}
	})

	t.Run("200 empty list", func(t *testing.T) {
		store := authStore(ownerID)
		h := NewSubmissionsHandler(store, "", "")
		req := httptest.NewRequest(http.MethodGet, "/api/submissions", nil)
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.List)).ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("got %d, want 200", rec.Code)
		}
	})
}

func TestCreate(t *testing.T) {
	ownerID := "owner-uuid"

	t.Run("401 no auth", func(t *testing.T) {
		store := &MockStore{}
		h := NewSubmissionsHandler(store, "", "")
		req := httptest.NewRequest(http.MethodPost, "/api/submissions", strings.NewReader(validSubmissionBody))
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Create)).ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("got %d, want 401", rec.Code)
		}
	})

	t.Run("500 limit check fails", func(t *testing.T) {
		store := authStore(ownerID)
		store.SubmissionsErr = errors.New("db error")
		h := NewSubmissionsHandler(store, "", "")
		req := httptest.NewRequest(http.MethodPost, "/api/submissions", strings.NewReader(validSubmissionBody))
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Create)).ServeHTTP(rec, req)
		if rec.Code != http.StatusInternalServerError {
			t.Errorf("got %d, want 500", rec.Code)
		}
	})

	t.Run("403 limit reached", func(t *testing.T) {
		store := authStore(ownerID)
		store.Submissions = []database.SalarySubmission{
			makeSubmission("sub-1"),
			makeSubmission("sub-2"),
			makeSubmission("sub-3"),
		}
		h := NewSubmissionsHandler(store, "", "")
		req := httptest.NewRequest(http.MethodPost, "/api/submissions", strings.NewReader(validSubmissionBody))
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Create)).ServeHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Errorf("got %d, want 403", rec.Code)
		}
	})

	t.Run("400 invalid json", func(t *testing.T) {
		store := authStore(ownerID)
		h := NewSubmissionsHandler(store, "", "")
		req := httptest.NewRequest(http.MethodPost, "/api/submissions", strings.NewReader("not-json"))
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Create)).ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got %d, want 400", rec.Code)
		}
	})

	t.Run("413 body too large", func(t *testing.T) {
		store := authStore(ownerID)
		h := NewSubmissionsHandler(store, "", "")
		// build a JSON body slightly over 64KB (64*1024 = 65536 bytes)
		bigBody := `{"company_name":"` + strings.Repeat("a", 66000) + `"}`
		req := httptest.NewRequest(http.MethodPost, "/api/submissions", strings.NewReader(bigBody))
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Create)).ServeHTTP(rec, req)
		if rec.Code != http.StatusRequestEntityTooLarge {
			t.Errorf("got %d, want 413", rec.Code)
		}
	})

	t.Run("429 rate limited has Retry-After header", func(t *testing.T) {
		sub := makeSubmission("new-sub")
		ownerStr := "owner-uuid"
		store := &MockStore{
			OwnerByToken:      &ownerStr,
			CreatedSubmission: &sub,
		}
		h := NewSubmissionsHandler(store, "", "") // empty secret = turnstile skipped
		body := `{"company_name":"Test Co","job_title":"Engineer","industry_id":1,"city_id":1,"seniority":"mid","work_arrangement":"office","employment_type":"full_time","base_salary":50000,"salary_year":2025,"turnstile_token":"tok"}`

		// Exhaust the burst bucket (3 allowed per hour)
		for i := 0; i < 3; i++ {
			req := httptest.NewRequest(http.MethodPost, "/api/submissions", strings.NewReader(body))
			rec := httptest.NewRecorder()
			h.Create(rec, req)
		}

		// 4th request must be rate limited
		req := httptest.NewRequest(http.MethodPost, "/api/submissions", strings.NewReader(body))
		rec := httptest.NewRecorder()
		h.Create(rec, req)

		if rec.Code != http.StatusTooManyRequests {
			t.Errorf("got %d, want 429", rec.Code)
		}
		if rec.Header().Get("Retry-After") == "" {
			t.Error("missing Retry-After header on 429 response")
		}
	})

	t.Run("400 validate fails", func(t *testing.T) {
		store := authStore(ownerID)
		h := NewSubmissionsHandler(store, "", "")
		body := `{"company_name":"","job_title":"","industry_id":0,"city_id":0,"seniority":"","work_arrangement":"","base_salary":0}`
		req := httptest.NewRequest(http.MethodPost, "/api/submissions", strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Create)).ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got %d, want 400", rec.Code)
		}
	})

	t.Run("500 create fails", func(t *testing.T) {
		store := authStore(ownerID)
		store.CreateErr = errors.New("db error")
		h := NewSubmissionsHandler(store, "", "")
		req := httptest.NewRequest(http.MethodPost, "/api/submissions", strings.NewReader(validSubmissionBody))
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Create)).ServeHTTP(rec, req)
		if rec.Code != http.StatusInternalServerError {
			t.Errorf("got %d, want 500", rec.Code)
		}
	})

	t.Run("201 happy path", func(t *testing.T) {
		sub := makeSubmission("new-sub-id")
		store := authStore(ownerID)
		store.CreatedSubmission = &sub
		h := NewSubmissionsHandler(store, "", "")
		req := httptest.NewRequest(http.MethodPost, "/api/submissions", strings.NewReader(validSubmissionBody))
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Create)).ServeHTTP(rec, req)
		if rec.Code != http.StatusCreated {
			t.Errorf("got %d, want 201", rec.Code)
		}
	})
}

func TestUpdate(t *testing.T) {
	ownerID := "owner-uuid"

	t.Run("401 no auth", func(t *testing.T) {
		store := &MockStore{}
		h := NewSubmissionsHandler(store, "", "")
		req := httptest.NewRequest(http.MethodPut, "/api/submissions/sub-123", strings.NewReader(validSubmissionBody))
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Update)).ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("got %d, want 401", rec.Code)
		}
	})

	t.Run("400 empty id", func(t *testing.T) {
		store := authStore(ownerID)
		h := NewSubmissionsHandler(store, "", "")
		// path /api/submissions only has 2 parts → empty ID
		req := httptest.NewRequest(http.MethodPut, "/api/submissions", strings.NewReader(validSubmissionBody))
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Update)).ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got %d, want 400", rec.Code)
		}
	})

	t.Run("400 invalid json", func(t *testing.T) {
		store := authStore(ownerID)
		h := NewSubmissionsHandler(store, "", "")
		req := httptest.NewRequest(http.MethodPut, "/api/submissions/sub-123", strings.NewReader("not-json"))
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Update)).ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got %d, want 400", rec.Code)
		}
	})

	t.Run("413 body too large", func(t *testing.T) {
		store := authStore(ownerID)
		h := NewSubmissionsHandler(store, "", "")
		bigBody := `{"company_name":"` + strings.Repeat("a", 66000) + `"}` // 66019 bytes > 64KB limit
		req := httptest.NewRequest(http.MethodPut, "/api/submissions/sub-123", strings.NewReader(bigBody))
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Update)).ServeHTTP(rec, req)
		if rec.Code != http.StatusRequestEntityTooLarge {
			t.Errorf("got %d, want 413", rec.Code)
		}
	})

	t.Run("400 validate fails", func(t *testing.T) {
		store := authStore(ownerID)
		h := NewSubmissionsHandler(store, "", "")
		body := `{"company_name":"","job_title":"","industry_id":0,"city_id":0,"seniority":"","work_arrangement":"","base_salary":0}`
		req := httptest.NewRequest(http.MethodPut, "/api/submissions/sub-123", strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Update)).ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got %d, want 400", rec.Code)
		}
	})

	t.Run("404 not found", func(t *testing.T) {
		store := authStore(ownerID)
		store.UpdateErr = sql.ErrNoRows
		h := NewSubmissionsHandler(store, "", "")
		req := httptest.NewRequest(http.MethodPut, "/api/submissions/sub-123", strings.NewReader(validSubmissionBody))
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Update)).ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Errorf("got %d, want 404", rec.Code)
		}
	})

	t.Run("500 update fails", func(t *testing.T) {
		store := authStore(ownerID)
		store.UpdateErr = errors.New("db error")
		h := NewSubmissionsHandler(store, "", "")
		req := httptest.NewRequest(http.MethodPut, "/api/submissions/sub-123", strings.NewReader(validSubmissionBody))
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Update)).ServeHTTP(rec, req)
		if rec.Code != http.StatusInternalServerError {
			t.Errorf("got %d, want 500", rec.Code)
		}
	})

	t.Run("204 happy path", func(t *testing.T) {
		store := authStore(ownerID)
		h := NewSubmissionsHandler(store, "", "")
		req := httptest.NewRequest(http.MethodPut, "/api/submissions/sub-123", strings.NewReader(validSubmissionBody))
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Update)).ServeHTTP(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Errorf("got %d, want 204", rec.Code)
		}
	})
}

func TestDelete(t *testing.T) {
	ownerID := "owner-uuid"

	t.Run("401 no auth", func(t *testing.T) {
		store := &MockStore{}
		h := NewSubmissionsHandler(store, "", "")
		req := httptest.NewRequest(http.MethodDelete, "/api/submissions/sub-123", nil)
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Delete)).ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("got %d, want 401", rec.Code)
		}
	})

	t.Run("400 empty id", func(t *testing.T) {
		store := authStore(ownerID)
		h := NewSubmissionsHandler(store, "", "")
		req := httptest.NewRequest(http.MethodDelete, "/api/submissions", nil)
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Delete)).ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got %d, want 400", rec.Code)
		}
	})

	t.Run("404 not found", func(t *testing.T) {
		store := authStore(ownerID)
		store.DeleteErr = sql.ErrNoRows
		h := NewSubmissionsHandler(store, "", "")
		req := httptest.NewRequest(http.MethodDelete, "/api/submissions/sub-123", nil)
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Delete)).ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Errorf("got %d, want 404", rec.Code)
		}
	})

	t.Run("500 delete fails", func(t *testing.T) {
		store := authStore(ownerID)
		store.DeleteErr = errors.New("db error")
		h := NewSubmissionsHandler(store, "", "")
		req := httptest.NewRequest(http.MethodDelete, "/api/submissions/sub-123", nil)
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Delete)).ServeHTTP(rec, req)
		if rec.Code != http.StatusInternalServerError {
			t.Errorf("got %d, want 500", rec.Code)
		}
	})

	t.Run("204 happy path", func(t *testing.T) {
		store := authStore(ownerID)
		h := NewSubmissionsHandler(store, "", "")
		req := httptest.NewRequest(http.MethodDelete, "/api/submissions/sub-123", nil)
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Delete)).ServeHTTP(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Errorf("got %d, want 204", rec.Code)
		}
	})
}

func TestIPVelocityCheck(t *testing.T) {
	ownerID := "owner-uuid"
	const secret = "test-secret"

	t.Run("429 when IP has 3+ submissions in last 24h", func(t *testing.T) {
		store := authStore(ownerID)
		store.IPHMACCount = 3 // mock: already 3 submissions from this IP
		h := NewSubmissionsHandler(store, "", secret)
		req := httptest.NewRequest(http.MethodPost, "/api/submissions", strings.NewReader(validSubmissionBody))
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Create)).ServeHTTP(rec, req)
		if rec.Code != http.StatusTooManyRequests {
			t.Errorf("got %d, want 429", rec.Code)
		}
		if rec.Header().Get("Retry-After") == "" {
			t.Error("missing Retry-After header on velocity-limited response")
		}
	})

	t.Run("201 when IP has 4 submissions (over limit)", func(t *testing.T) {
		sub := makeSubmission("new-sub")
		store := authStore(ownerID)
		store.IPHMACCount = 4
		store.CreatedSubmission = &sub
		h := NewSubmissionsHandler(store, "", secret)
		req := httptest.NewRequest(http.MethodPost, "/api/submissions", strings.NewReader(validSubmissionBody))
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Create)).ServeHTTP(rec, req)
		if rec.Code != http.StatusTooManyRequests {
			t.Errorf("got %d, want 201", rec.Code)
		}
	})

	t.Run("velocity check skipped when secret is empty", func(t *testing.T) {
		sub := makeSubmission("new-sub")
		store := authStore(ownerID)
		store.IPHMACCount = 99 // would trigger if secret were set
		store.CreatedSubmission = &sub
		h := NewSubmissionsHandler(store, "", "") // no secret
		req := httptest.NewRequest(http.MethodPost, "/api/submissions", strings.NewReader(validSubmissionBody))
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Create)).ServeHTTP(rec, req)
		if rec.Code != http.StatusCreated {
			t.Errorf("got %d, want 201 (velocity check must be skipped in dev mode)", rec.Code)
		}
	})
}

func TestHMACIP(t *testing.T) {
	t.Run("same IP and secret produce same hash", func(t *testing.T) {
		h1 := hmacIP("secret", "1.2.3.4")
		h2 := hmacIP("secret", "1.2.3.4")
		if h1 != h2 {
			t.Error("expected deterministic output")
		}
	})

	t.Run("different IPs produce different hashes", func(t *testing.T) {
		h1 := hmacIP("secret", "1.2.3.4")
		h2 := hmacIP("secret", "1.2.3.5")
		if h1 == h2 {
			t.Error("expected different IPs to produce different hashes")
		}
	})

	t.Run("same IP with different secrets produce different hashes", func(t *testing.T) {
		h1 := hmacIP("secret-a", "1.2.3.4")
		h2 := hmacIP("secret-b", "1.2.3.4")
		if h1 == h2 {
			t.Error("expected different secrets to produce different hashes")
		}
	})

	t.Run("hash is 64 hex chars (SHA-256 output)", func(t *testing.T) {
		h := hmacIP("secret", "1.2.3.4")
		if len(h) != 64 {
			t.Errorf("expected 64 chars, got %d", len(h))
		}
	})
}
