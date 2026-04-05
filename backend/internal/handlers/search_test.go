package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/platafyi/plata.fyi/internal/database"
)

func TestSalaries(t *testing.T) {
	tests := []struct {
		name     string
		store    *MockStore
		wantCode int
	}{
		{
			name:     "200 ok",
			store:    &MockStore{},
			wantCode: http.StatusOK,
		},
		{
			name:     "500 search fails",
			store:    &MockStore{SearchErr: errors.New("db error")},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := NewSearchHandler(tc.store)
			req := httptest.NewRequest(http.MethodGet, "/api/salaries", nil)
			rec := httptest.NewRecorder()
			h.Salaries(rec, req)

			if rec.Code != tc.wantCode {
				t.Errorf("got %d, want %d", rec.Code, tc.wantCode)
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	now := time.Now()
	sub := database.SalarySubmission{
		ID: "sub-1", OwnerID: "owner-1", CompanyName: "Acme", JobTitle: "Engineer",
		IndustryID: 1, IndustryName: "Tech", IndustrySlug: "tech",
		CityID: 1, CityName: "Skopje", CitySlug: "skopje",
		Seniority: "mid", YearsAtCompany: 2, YearsExperience: 5,
		WorkArrangement: "office", EmploymentType: "full_time",
		BaseSalary: 50000, SalaryYear: 2024, IsApproved: true,
		CreatedAt: now, UpdatedAt: now,
	}

	tests := []struct {
		name     string
		id       string
		store    *MockStore
		wantCode int
	}{
		{
			name:     "400 missing id",
			id:       "",
			store:    &MockStore{},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "404 not found",
			id:       "sub-1",
			store:    &MockStore{SubmissionByID: nil},
			wantCode: http.StatusNotFound,
		},
		{
			name:     "500 db error",
			id:       "sub-1",
			store:    &MockStore{SubmissionByIDErr: errors.New("db error")},
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "200 ok",
			id:       "sub-1",
			store:    &MockStore{SubmissionByID: &sub},
			wantCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := NewSearchHandler(tc.store)
			req := httptest.NewRequest(http.MethodGet, "/api/salaries/"+tc.id, nil)
			if tc.id != "" {
				req.SetPathValue("id", tc.id)
			}
			rec := httptest.NewRecorder()
			h.GetByID(rec, req)

			if rec.Code != tc.wantCode {
				t.Errorf("got %d, want %d", rec.Code, tc.wantCode)
			}
		})
	}
}

// TestSalariesOrderNotByCreatedAt verifies that the handler returns results in
// exactly the order the store provides, not re-sorted by created_at.
// Ordering is delegated to the store (ORDER BY ctid DESC in PostgresStore),
// so the handler must never impose its own date-based sort.
func TestSalariesOrderNotByCreatedAt(t *testing.T) {
	older := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	newer := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	// Store returns newer-first — opposite of what created_at ASC would give.
	// If the handler re-sorted by created_at it would flip these.
	storeOrder := []database.SalarySubmission{
		{ID: "newer", CreatedAt: newer, UpdatedAt: newer},
		{ID: "older", CreatedAt: older, UpdatedAt: older},
	}

	store := &MockStore{SearchResults: storeOrder, SearchTotal: 2}
	h := NewSearchHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/api/salaries", nil)
	rec := httptest.NewRecorder()
	h.Salaries(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got %d, want 200", rec.Code)
	}

	body := rec.Body.String()
	newerPos := strings.Index(body, "newer")
	olderPos := strings.Index(body, "older")
	if newerPos == -1 || olderPos == -1 {
		t.Fatal("expected both IDs in response body")
	}
	if newerPos > olderPos {
		t.Error("handler re-sorted results by created_at: older appeared before newer")
	}
}

func TestStats(t *testing.T) {
	tests := []struct {
		name     string
		store    *MockStore
		wantCode int
	}{
		{
			name:     "200 ok",
			store:    &MockStore{},
			wantCode: http.StatusOK,
		},
		{
			name:     "500 query fails",
			store:    &MockStore{StatsErr: errors.New("db error")},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := NewSearchHandler(tc.store)
			req := httptest.NewRequest(http.MethodGet, "/api/salaries/stats", nil)
			rec := httptest.NewRecorder()
			h.Stats(rec, req)

			if rec.Code != tc.wantCode {
				t.Errorf("got %d, want %d", rec.Code, tc.wantCode)
			}
		})
	}
}
