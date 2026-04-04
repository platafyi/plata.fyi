package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/platafyi/plata.fyi/internal/database"
)

func TestCompaniesSearch(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		store    *MockStore
		wantCode int
	}{
		{
			name:     "200 empty q too short",
			query:    "a",
			store:    &MockStore{},
			wantCode: http.StatusOK,
		},
		{
			name:     "200 empty db query fails",
			query:    "ab",
			store:    &MockStore{CompaniesErr: errors.New("db error")},
			wantCode: http.StatusOK,
		},
		{
			name:  "200 results",
			query: "acme",
			store: &MockStore{
				Companies: []database.Company{{Name: "Acme DOO", RegNo: ""}},
			},
			wantCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := NewCompaniesHandler(tc.store)
			req := httptest.NewRequest(http.MethodGet, "/api/companies?q="+tc.query, nil)
			rec := httptest.NewRecorder()
			h.Search(rec, req)

			if rec.Code != tc.wantCode {
				t.Errorf("got %d, want %d", rec.Code, tc.wantCode)
			}
		})
	}
}
