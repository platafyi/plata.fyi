package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/platafyi/plata.fyi/internal/database"
)

func TestIndustries(t *testing.T) {
	tests := []struct {
		name     string
		store    *MockStore
		wantCode int
	}{
		{
			name: "200 ok",
			store: &MockStore{
				Industries: []database.Industry{{ID: 1, Slug: "tech", Name: "Technology"}},
			},
			wantCode: http.StatusOK,
		},
		{
			name:     "500 query fails",
			store:    &MockStore{IndustriesErr: errors.New("db error")},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := NewMetaHandler(tc.store)
			req := httptest.NewRequest(http.MethodGet, "/api/industries", nil)
			rec := httptest.NewRecorder()
			h.Industries(rec, req)

			if rec.Code != tc.wantCode {
				t.Errorf("got %d, want %d", rec.Code, tc.wantCode)
			}
		})
	}
}

func TestCities(t *testing.T) {
	tests := []struct {
		name     string
		store    *MockStore
		wantCode int
	}{
		{
			name: "200 ok",
			store: &MockStore{
				Cities: []database.City{{ID: 1, Slug: "skopje", Name: "Скопје"}},
			},
			wantCode: http.StatusOK,
		},
		{
			name:     "500 query fails",
			store:    &MockStore{CitiesErr: errors.New("db error")},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := NewMetaHandler(tc.store)
			req := httptest.NewRequest(http.MethodGet, "/api/cities", nil)
			rec := httptest.NewRecorder()
			h.Cities(rec, req)

			if rec.Code != tc.wantCode {
				t.Errorf("got %d, want %d", rec.Code, tc.wantCode)
			}
		})
	}
}
