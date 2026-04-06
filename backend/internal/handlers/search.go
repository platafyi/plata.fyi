package handlers

import (
	"net/http"
	"strconv"

	"github.com/platafyi/plata.fyi/internal/database"
)

type SearchHandler struct {
	store database.Store
}

func NewSearchHandler(store database.Store) *SearchHandler {
	return &SearchHandler{store: store}
}

func (h *SearchHandler) Salaries(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	f := database.SearchFilters{
		IndustrySlug:    q.Get("industry"),
		CitySlug:        q.Get("city"),
		Seniority:       q.Get("seniority"),
		WorkArrangement: q.Get("arrangement"),
		MinSalary:       parseInt(q.Get("min_salary")),
		MaxSalary:       parseInt(q.Get("max_salary")),
		CompanyType:     q.Get("company_type"),
		Page:            parseInt(q.Get("page")),
		PageSize:        parseInt(q.Get("page_size")),
	}

	subs, total, err := h.store.SearchSalaries(r.Context(), f)
	if err != nil {
		jsonError(w, "Грешка при пребарување", http.StatusInternalServerError)
		return
	}
	if subs == nil {
		subs = []database.SalarySubmission{}
	}

	jsonOK(w, map[string]interface{}{
		"data":      subs,
		"total":     total,
		"page":      f.Page,
		"page_size": f.PageSize,
	})
}

func (h *SearchHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		jsonError(w, "Невалиден ID", http.StatusBadRequest)
		return
	}

	sub, err := h.store.GetSubmissionByID(r.Context(), id)
	if err != nil {
		jsonError(w, "Грешка при вчитување", http.StatusInternalServerError)
		return
	}
	if sub == nil {
		jsonError(w, "Не е пронајдено", http.StatusNotFound)
		return
	}

	jsonOK(w, sub)
}

func (h *SearchHandler) Stats(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	groupBy := q.Get("group_by") // "industry" or "city"

	f := database.SearchFilters{
		IndustrySlug: q.Get("industry"),
		CitySlug:     q.Get("city"),
		Seniority:    q.Get("seniority"),
		CompanyType:  q.Get("company_type"),
	}

	stats, err := h.store.GetSalaryStats(r.Context(), groupBy, f)
	if err != nil {
		jsonError(w, "Грешка при статистики", http.StatusInternalServerError)
		return
	}
	if stats == nil {
		stats = []database.SalaryStats{}
	}

	jsonOK(w, stats)
}

func parseInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
