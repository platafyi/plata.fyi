package handlers

import (
	"net/http"

	"github.com/platafyi/plata.fyi/internal/database"
)

type JobTitlesHandler struct {
	store database.Store
}

func NewJobTitlesHandler(store database.Store) *JobTitlesHandler {
	return &JobTitlesHandler{store: store}
}

func (h *JobTitlesHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if len(q) < 2 {
		jsonOK(w, []string{})
		return
	}

	results, err := h.store.SearchJobTitles(r.Context(), q)
	if err != nil {
		jsonOK(w, []string{})
		return
	}
	if results == nil {
		results = []string{}
	}
	jsonOK(w, results)
}
