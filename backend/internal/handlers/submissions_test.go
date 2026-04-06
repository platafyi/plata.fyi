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

func TestValidateRequestSalaryBounds(t *testing.T) {
	h := NewSubmissionsHandler(&MockStore{}, "")

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
		h := NewSubmissionsHandler(store, "")
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
		h := NewSubmissionsHandler(store, "")
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
		h := NewSubmissionsHandler(store, "")
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
		h := NewSubmissionsHandler(store, "")
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
		h := NewSubmissionsHandler(store, "")
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
		h := NewSubmissionsHandler(store, "")
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
		h := NewSubmissionsHandler(store, "")
		req := httptest.NewRequest(http.MethodPost, "/api/submissions", strings.NewReader("not-json"))
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Create)).ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got %d, want 400", rec.Code)
		}
	})

	t.Run("400 validate fails", func(t *testing.T) {
		store := authStore(ownerID)
		h := NewSubmissionsHandler(store, "")
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
		h := NewSubmissionsHandler(store, "")
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
		h := NewSubmissionsHandler(store, "")
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
		h := NewSubmissionsHandler(store, "")
		req := httptest.NewRequest(http.MethodPut, "/api/submissions/sub-123", strings.NewReader(validSubmissionBody))
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Update)).ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("got %d, want 401", rec.Code)
		}
	})

	t.Run("400 empty id", func(t *testing.T) {
		store := authStore(ownerID)
		h := NewSubmissionsHandler(store, "")
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
		h := NewSubmissionsHandler(store, "")
		req := httptest.NewRequest(http.MethodPut, "/api/submissions/sub-123", strings.NewReader("not-json"))
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Update)).ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got %d, want 400", rec.Code)
		}
	})

	t.Run("400 validate fails", func(t *testing.T) {
		store := authStore(ownerID)
		h := NewSubmissionsHandler(store, "")
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
		h := NewSubmissionsHandler(store, "")
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
		h := NewSubmissionsHandler(store, "")
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
		h := NewSubmissionsHandler(store, "")
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
		h := NewSubmissionsHandler(store, "")
		req := httptest.NewRequest(http.MethodDelete, "/api/submissions/sub-123", nil)
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Delete)).ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("got %d, want 401", rec.Code)
		}
	})

	t.Run("400 empty id", func(t *testing.T) {
		store := authStore(ownerID)
		h := NewSubmissionsHandler(store, "")
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
		h := NewSubmissionsHandler(store, "")
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
		h := NewSubmissionsHandler(store, "")
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
		h := NewSubmissionsHandler(store, "")
		req := httptest.NewRequest(http.MethodDelete, "/api/submissions/sub-123", nil)
		req.Header.Set("Authorization", "Bearer tok")
		rec := httptest.NewRecorder()
		middleware.Auth(store)(http.HandlerFunc(h.Delete)).ServeHTTP(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Errorf("got %d, want 204", rec.Code)
		}
	})
}
