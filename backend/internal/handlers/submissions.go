package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/platafyi/plata.fyi/internal/database"
	"github.com/platafyi/plata.fyi/internal/middleware"
	"golang.org/x/time/rate"
)

type SubmissionsHandler struct {
	store           database.Store
	turnstileSecret string
	createRL        *middleware.KeyRateLimiter
}

func NewSubmissionsHandler(store database.Store, turnstileSecret string) *SubmissionsHandler {
	return &SubmissionsHandler{
		store:           store,
		turnstileSecret: turnstileSecret,
		createRL:        middleware.NewKeyRateLimiter(rate.Limit(3.0/3600), 3), // 3/hour per IP
	}
}

var validSeniorities = map[string]bool{
	"intern": true, "junior": true, "mid": true, "senior": true,
	"lead": true, "principal": true, "manager": true, "director": true, "executive": true,
}

var validArrangements = map[string]bool{
	"office": true, "hybrid": true, "remote": true,
}

var validEmploymentTypes = map[string]bool{
	"full_time": true, "part_time": true,
}

var validBonusTypes = map[string]bool{
	"annual": true, "performance": true, "signing": true, "project": true, "other": true,
}

var validFrequencies = map[string]bool{
	"monthly": true, "quarterly": true, "annual": true, "one_time": true,
}

type submissionRequest struct {
	CompanyName     string         `json:"company_name"`
	CompanyRegNo    string         `json:"company_reg_no"`
	JobTitle        string         `json:"job_title"`
	IndustryID      int            `json:"industry_id"`
	CityID          int            `json:"city_id"`
	Seniority       string         `json:"seniority"`
	YearsAtCompany  int            `json:"years_at_company"`
	YearsExperience int            `json:"years_experience"`
	WorkArrangement string         `json:"work_arrangement"`
	EmploymentType  string         `json:"employment_type"`
	HoursPerWeek    *int           `json:"hours_per_week"`
	BaseSalary      int            `json:"base_salary"`
	SalaryYear      int            `json:"salary_year"`
	Bonuses         []bonusRequest `json:"bonuses"`
	TurnstileToken  string         `json:"turnstile_token"`
}

type bonusRequest struct {
	BonusType string `json:"bonus_type"`
	Amount    int    `json:"amount"`
	Frequency string `json:"frequency"`
}

func (h *SubmissionsHandler) validateRequest(req *submissionRequest) string {
	if strings.TrimSpace(req.CompanyName) == "" {
		return "Името на компанијата е задолжително"
	}
	if strings.TrimSpace(req.JobTitle) == "" {
		return "Работната позиција е задолжителна"
	}
	if req.IndustryID <= 0 {
		return "Изберете индустрија"
	}
	if req.CityID <= 0 {
		return "Изберете град"
	}
	if !validSeniorities[req.Seniority] {
		return "Невалидно ниво на искуство"
	}
	if req.YearsAtCompany < 0 {
		return "Годините во компанијата не можат да бидат негативни"
	}
	if req.YearsExperience < 0 {
		return "Годините искуство не можат да бидат негативни"
	}
	if !validArrangements[req.WorkArrangement] {
		return "Невалиден начин на работа"
	}
	if req.EmploymentType == "" {
		req.EmploymentType = "full_time"
	}
	if !validEmploymentTypes[req.EmploymentType] {
		return "Невалиден тип на вработување"
	}
	if req.EmploymentType == "part_time" && (req.HoursPerWeek == nil || *req.HoursPerWeek < 1 || *req.HoursPerWeek > 168) {
		h := 30
		req.HoursPerWeek = &h
	}
	if req.EmploymentType == "full_time" {
		req.HoursPerWeek = nil
	}
	if req.BaseSalary <= 0 {
		return "Основната плата мора да биде позитивна"
	}
	currentYear := time.Now().Year()
	if req.SalaryYear == 0 {
		req.SalaryYear = currentYear
	}
	if req.SalaryYear < 2000 || req.SalaryYear > currentYear {
		return "Невалидна година на плата"
	}
	for _, b := range req.Bonuses {
		if !validBonusTypes[b.BonusType] {
			return "Невалиден тип на бонус"
		}
		if b.Amount <= 0 {
			return "Износот на бонусот мора да биде позитивен"
		}
		if !validFrequencies[b.Frequency] {
			return "Невалидна фреквенција на бонус"
		}
	}
	return ""
}

func (h *SubmissionsHandler) List(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := middleware.GetOwnerID(r.Context())
	if !ok {
		jsonError(w, "Потребна е автентикација", http.StatusUnauthorized)
		return
	}

	subs, err := h.store.GetSubmissionsByOwner(r.Context(), ownerID)
	if err != nil {
		jsonError(w, "Грешка при вчитување", http.StatusInternalServerError)
		return
	}
	if subs == nil {
		subs = []database.SalarySubmission{}
	}
	jsonOK(w, subs)
}

func (h *SubmissionsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req submissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Невалидно тело на барањето", http.StatusBadRequest)
		return
	}

	var ownerID string
	var newSessionToken string

	token := bearerToken(r)
	if token != "" {
		// Authenticated path: existing session
		ownerIDPtr, err := h.store.GetOwnerByToken(r.Context(), token)
		if err != nil {
			jsonError(w, "Грешка при автентикација", http.StatusInternalServerError)
			return
		}
		if ownerIDPtr == nil {
			jsonError(w, "Потребна е автентикација", http.StatusUnauthorized)
			return
		}
		ownerID = *ownerIDPtr
	} else {
		// Unauthenticated path: verify Turnstile, issue new session
		if req.TurnstileToken == "" {
			jsonError(w, "Потребна е верификација", http.StatusBadRequest)
			return
		}
		remoteIP := middleware.RealIP(r)
		if !h.createRL.Allow(remoteIP) {
			jsonError(w, "Премногу барања. Обидете се подоцна.", http.StatusTooManyRequests)
			return
		}
		if err := verifyTurnstile(r.Context(), h.turnstileSecret, req.TurnstileToken, remoteIP); err != nil {
			jsonError(w, "Верификацијата не успеа", http.StatusForbidden)
			return
		}

		newSessionToken = randHex(32)
		if err := h.store.InsertToken(r.Context(), newSessionToken); err != nil {
			jsonError(w, "Грешка при создавање сесија", http.StatusInternalServerError)
			return
		}
		ownerIDPtr, err := h.store.GetOwnerByToken(r.Context(), newSessionToken)
		if err != nil || ownerIDPtr == nil {
			jsonError(w, "Грешка при создавање сесија", http.StatusInternalServerError)
			return
		}
		ownerID = *ownerIDPtr
	}

	// Enforce 3-submission limit per owner
	existing, err := h.store.GetSubmissionsByOwner(r.Context(), ownerID)
	if err != nil {
		jsonError(w, "Грешка при проверка", http.StatusInternalServerError)
		return
	}
	if len(existing) >= 3 {
		jsonError(w, "Максимален број на записи е 3 по токен", http.StatusForbidden)
		return
	}

	if errMsg := h.validateRequest(&req); errMsg != "" {
		jsonError(w, errMsg, http.StatusBadRequest)
		return
	}

	bonuses := make([]database.BonusInput, len(req.Bonuses))
	for i, b := range req.Bonuses {
		bonuses[i] = database.BonusInput{
			BonusType: b.BonusType,
			Amount:    b.Amount,
			Frequency: b.Frequency,
		}
	}

	sub, err := h.store.CreateSubmission(r.Context(), database.CreateSubmissionInput{
		OwnerID:         ownerID,
		CompanyName:     strings.TrimSpace(req.CompanyName),
		CompanyRegNo:    strings.TrimSpace(req.CompanyRegNo),
		JobTitle:        strings.TrimSpace(req.JobTitle),
		IndustryID:      req.IndustryID,
		CityID:          req.CityID,
		Seniority:       req.Seniority,
		YearsAtCompany:  req.YearsAtCompany,
		YearsExperience: req.YearsExperience,
		WorkArrangement: req.WorkArrangement,
		EmploymentType:  req.EmploymentType,
		HoursPerWeek:    req.HoursPerWeek,
		BaseSalary:      req.BaseSalary,
		SalaryYear:      req.SalaryYear,
		Bonuses:         bonuses,
	})
	if err != nil {
		jsonError(w, "Грешка при зачувување", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if newSessionToken != "" {
		type createResponse struct {
			*database.SalarySubmission
			SessionToken string `json:"session_token"`
		}
		json.NewEncoder(w).Encode(createResponse{SalarySubmission: sub, SessionToken: newSessionToken})
	} else {
		json.NewEncoder(w).Encode(sub)
	}
}

func (h *SubmissionsHandler) Update(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := middleware.GetOwnerID(r.Context())
	if !ok {
		jsonError(w, "Потребна е автентикација", http.StatusUnauthorized)
		return
	}

	id := submissionIDFromPath(r)
	if id == "" {
		jsonError(w, "Недостасува ID", http.StatusBadRequest)
		return
	}

	var req submissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Невалидно тело на барањето", http.StatusBadRequest)
		return
	}

	if errMsg := h.validateRequest(&req); errMsg != "" {
		jsonError(w, errMsg, http.StatusBadRequest)
		return
	}

	bonuses := make([]database.BonusInput, len(req.Bonuses))
	for i, b := range req.Bonuses {
		bonuses[i] = database.BonusInput{BonusType: b.BonusType, Amount: b.Amount, Frequency: b.Frequency}
	}

	err := h.store.UpdateSubmission(r.Context(), id, ownerID, database.CreateSubmissionInput{
		OwnerID:         ownerID,
		CompanyName:     strings.TrimSpace(req.CompanyName),
		CompanyRegNo:    strings.TrimSpace(req.CompanyRegNo),
		JobTitle:        strings.TrimSpace(req.JobTitle),
		IndustryID:      req.IndustryID,
		CityID:          req.CityID,
		Seniority:       req.Seniority,
		YearsAtCompany:  req.YearsAtCompany,
		YearsExperience: req.YearsExperience,
		WorkArrangement: req.WorkArrangement,
		EmploymentType:  req.EmploymentType,
		HoursPerWeek:    req.HoursPerWeek,
		BaseSalary:      req.BaseSalary,
		SalaryYear:      req.SalaryYear,
		Bonuses:         bonuses,
	})
	if err == sql.ErrNoRows {
		jsonError(w, "Не е пронајдено или немате право да го уредите", http.StatusNotFound)
		return
	}
	if err != nil {
		jsonError(w, "Грешка при ажурирање", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SubmissionsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := middleware.GetOwnerID(r.Context())
	if !ok {
		jsonError(w, "Потребна е автентикација", http.StatusUnauthorized)
		return
	}

	id := submissionIDFromPath(r)
	if id == "" {
		jsonError(w, "Недостасува ID", http.StatusBadRequest)
		return
	}

	err := h.store.DeleteSubmission(r.Context(), id, ownerID)
	if err == sql.ErrNoRows {
		jsonError(w, "Не е пронајдено или немате право да го избришете", http.StatusNotFound)
		return
	}
	if err != nil {
		jsonError(w, "Грешка при бришење", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// submissionIDFromPath extracts the UUID from paths like /api/submissions/UUID
func submissionIDFromPath(r *http.Request) string {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	// Expected: api / submissions / <id>
	if len(parts) >= 3 {
		return parts[len(parts)-1]
	}
	return ""
}
