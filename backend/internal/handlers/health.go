package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/platafyi/plata.fyi/internal/database"
)

type HealthHandler struct {
	store database.Store
}

func NewHealthHandler(store database.Store) *HealthHandler {
	return &HealthHandler{store: store}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	status := "ok"
	code := http.StatusOK

	if err := h.store.Ping(r.Context()); err != nil {
		status = "db_error"
		code = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}
