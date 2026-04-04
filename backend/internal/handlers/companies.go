package handlers

import (
	"net/http"

	"github.com/platafyi/plata.fyi/internal/database"
)

type CompaniesHandler struct {
	store database.Store
}

func NewCompaniesHandler(store database.Store) *CompaniesHandler {
	return &CompaniesHandler{store: store}
}

type CompaniesResponse struct {
	Results []database.Company `json:"results"`
}

func (h *CompaniesHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if len(q) < 2 {
		jsonOK(w, CompaniesResponse{Results: []database.Company{}})
		return
	}

	results, err := h.store.SearchCompanies(r.Context(), q)
	if err != nil {
		jsonOK(w, CompaniesResponse{Results: []database.Company{}})
		return
	}
	if results == nil {
		results = []database.Company{}
	}
	jsonOK(w, CompaniesResponse{Results: results})
}
