package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/platafyi/plata.fyi/internal/database"
	"github.com/platafyi/plata.fyi/internal/middleware"
	"golang.org/x/time/rate"
)

type SubmissionsHandler struct {
	store           database.Store
	turnstileSecret string
	ipHMACSecret    string
	createRL        *middleware.KeyRateLimiter
}

func NewSubmissionsHandler(store database.Store, turnstileSecret, ipHMACSecret string) *SubmissionsHandler {
	return &SubmissionsHandler{
		store:           store,
		turnstileSecret: turnstileSecret,
		ipHMACSecret:    ipHMACSecret,
		createRL:        middleware.NewKeyRateLimiter(rate.Limit(3.0/3600), 3), // 3/hour per IP
	}
}

// hmacIP returns HMAC-SHA256(secret, ip) as a hex string.
// The raw IP is never stored. Without the secret key the hash cannot be reversed.
// Keep IP_HMAC_SECRET strictly separate from DB backups, never log or export it.
func hmacIP(secret, ip string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(ip))
	return hex.EncodeToString(mac.Sum(nil))
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

var validCompanyTypes = map[string]bool{
	"domestic": true, "foreign": true,
}

const minSalaryMKD = 26_046 // minimum wage MK 2026
const maxSalaryMKD = 2_000_000

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
	CompanyType     string         `json:"company_type"`
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
	if utf8.RuneCountInString(strings.TrimSpace(req.CompanyName)) > 100 {
		return "Името на компанијата е предолго (максимум 100 знаци)"
	}
	if strings.TrimSpace(req.JobTitle) == "" {
		return "Работната позиција е задолжителна"
	}
	if utf8.RuneCountInString(strings.TrimSpace(req.JobTitle)) > 100 {
		return "Работната позиција е предолга (максимум 100 знаци)"
	}
	if len(req.CompanyRegNo) > 20 {
		return "Регистарскиот број е предолг (максимум 20 знаци)"
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
	if req.YearsAtCompany < 0 || req.YearsAtCompany > 60 {
		return "Годините во компанијата мора да бидат помеѓу 0 и 60"
	}
	if req.YearsExperience < 0 || req.YearsExperience > 60 {
		return "Годините искуство мора да бидат помеѓу 0 и 60"
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
	effectiveMin := minSalaryMKD
	if req.EmploymentType == "part_time" {
		effectiveMin = 1
	}
	if req.BaseSalary < effectiveMin || req.BaseSalary > maxSalaryMKD {
		return "Платата мора да биде помеѓу 26.046 и 2.000.000 МКД"
	}
	if req.CompanyType == "" {
		req.CompanyType = "domestic"
	}
	if !validCompanyTypes[req.CompanyType] {
		return "Невалиден тип на компанија"
	}
	currentYear := time.Now().Year()
	if req.SalaryYear == 0 {
		req.SalaryYear = currentYear
	}
	if req.SalaryYear < 2000 || req.SalaryYear > currentYear {
		return "Невалидна година на плата"
	}
	if len(req.Bonuses) > 10 {
		return "Максималниот број на бонуси е 10"
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
	r.Body = http.MaxBytesReader(w, r.Body, 64*1024)
	var req submissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			jsonError(w, "Барањето е премногу големо", http.StatusRequestEntityTooLarge)
			return
		}
		jsonError(w, "Невалидно тело на барањето", http.StatusBadRequest)
		return
	}

	remoteIP := middleware.RealIP(r)

	// Velocity check: block if this IP has submitted 3+ times in the last 24h.
	// The IP is never stored, only its HMAC-SHA256 hash is used for comparison.
	if h.ipHMACSecret != "" {
		ipHash := hmacIP(h.ipHMACSecret, remoteIP)
		count, err := h.store.CountRecentSubmissionsByIPHMAC(r.Context(), ipHash, time.Now().Add(-12*time.Hour))
		if err != nil {
			jsonError(w, "Грешка при проверка", http.StatusInternalServerError)
			return
		}
		if count >= 3 {
			w.Header().Set("Retry-After", "86400")
			jsonError(w, "Премногу барања. Обидете се подоцна.", http.StatusTooManyRequests)
			return
		}
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
			w.Header().Set("Retry-After", "3600")
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

	var ipHash string
	if h.ipHMACSecret != "" {
		ipHash = hmacIP(h.ipHMACSecret, remoteIP)
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
		CompanyType:     req.CompanyType,
		SubmitterIPHMAC: ipHash,
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

	r.Body = http.MaxBytesReader(w, r.Body, 64*1024)
	var req submissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			jsonError(w, "Барањето е премногу големо", http.StatusRequestEntityTooLarge)
			return
		}
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
		CompanyType:     req.CompanyType,
		Bonuses:         bonuses,
	})
	if errors.Is(err, sql.ErrNoRows) {
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
	if errors.Is(err, sql.ErrNoRows) {
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
