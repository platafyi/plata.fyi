package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	tests := []struct {
		name     string
		store    *MockStore
		wantCode int
	}{
		{
			name:     "200 db ok",
			store:    &MockStore{},
			wantCode: http.StatusOK,
		},
		{
			name:     "503 db error",
			store:    &MockStore{PingErr: errors.New("db error")},
			wantCode: http.StatusServiceUnavailable,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := NewHealthHandler(tc.store)
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rec := httptest.NewRecorder()
			h.Health(rec, req)

			if rec.Code != tc.wantCode {
				t.Errorf("got %d, want %d", rec.Code, tc.wantCode)
			}
		})
	}
}
