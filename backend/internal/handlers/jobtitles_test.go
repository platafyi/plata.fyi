package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJobTitlesSearch(t *testing.T) {
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
			query:    "en",
			store:    &MockStore{JobTitlesErr: errors.New("db error")},
			wantCode: http.StatusOK,
		},
		{
			name:     "200 results",
			query:    "engineer",
			store:    &MockStore{JobTitles: []string{"Software Engineer"}},
			wantCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := NewJobTitlesHandler(tc.store)
			req := httptest.NewRequest(http.MethodGet, "/api/job-titles?q="+tc.query, nil)
			rec := httptest.NewRecorder()
			h.Search(rec, req)

			if rec.Code != tc.wantCode {
				t.Errorf("got %d, want %d", rec.Code, tc.wantCode)
			}
		})
	}
}
