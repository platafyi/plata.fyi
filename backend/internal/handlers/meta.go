package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/platafyi/plata.fyi/internal/database"
)

type MetaHandler struct {
	store database.Store
}

func NewMetaHandler(store database.Store) *MetaHandler {
	return &MetaHandler{store: store}
}

func (h *MetaHandler) Industries(w http.ResponseWriter, r *http.Request) {
	industries, err := h.store.GetIndustries(r.Context())
	if err != nil {
		jsonError(w, "Грешка при вчитување на индустриите", http.StatusInternalServerError)
		return
	}
	jsonOK(w, industries)
}

func (h *MetaHandler) Cities(w http.ResponseWriter, r *http.Request) {
	cities, err := h.store.GetCities(r.Context())
	if err != nil {
		jsonError(w, "Грешка при вчитување на градовите", http.StatusInternalServerError)
		return
	}
	jsonOK(w, cities)
}

// ── Helpers ──────────────────────────────────────────────────────────────────

func jsonOK(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
