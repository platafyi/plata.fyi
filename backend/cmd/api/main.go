package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/platafyi/plata.fyi/internal/config"
	"github.com/platafyi/plata.fyi/internal/database"
	"github.com/platafyi/plata.fyi/internal/handlers"
	"github.com/platafyi/plata.fyi/internal/middleware"

	"golang.org/x/time/rate"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := database.New(cfg.DBURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	store := database.NewPostgresStore(db)

	// Handlers
	healthH := handlers.NewHealthHandler(store)
	metaH := handlers.NewMetaHandler(store)
	authH := handlers.NewAuthHandler(store)
	submissionsH := handlers.NewSubmissionsHandler(store, cfg.TurnstileSecret, cfg.IPHMACSecret)
	searchH := handlers.NewSearchHandler(store)
	companiesH := handlers.NewCompaniesHandler(store)
	jobTitlesH := handlers.NewJobTitlesHandler(store)

	// Rate limiters
	globalRL := middleware.NewIPRateLimiter(rate.Limit(10), 30)    // 10 req/s burst 30
	authRL := middleware.NewIPRateLimiter(rate.Limit(5.0/3600), 5) // 5 req/hour per IP

	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/health", healthH.Health)
	mux.HandleFunc("/api/industries", metaH.Industries)
	mux.HandleFunc("/api/cities", metaH.Cities)
	mux.HandleFunc("GET /api/salaries", searchH.Salaries)
	mux.HandleFunc("GET /api/salaries/stats", searchH.Stats)
	mux.HandleFunc("GET /api/salaries/{id}", searchH.GetByID)
	mux.HandleFunc("/api/companies", companiesH.Search)
	mux.HandleFunc("/api/job-titles", jobTitlesH.Search)

	// Auth routes (with stricter rate limit)
	authMux := http.NewServeMux()
	authMux.HandleFunc("/api/auth/session", authH.DeleteSession)
	mux.Handle("/api/auth/", authRL.Middleware(authMux))

	// Submission routes
	authMiddleware := middleware.Auth(store)

	// GET /api/submissions — requires auth
	mux.Handle("GET /api/submissions", authMiddleware(http.HandlerFunc(submissionsH.List)))

	// POST /api/submissions — handles auth internally (existing token or Turnstile for new session)
	mux.HandleFunc("POST /api/submissions", submissionsH.Create)

	// PUT/DELETE /api/submissions/{id} — requires auth
	mux.Handle("/api/submissions/", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			submissionsH.Update(w, r)
		case http.MethodDelete:
			submissionsH.Delete(w, r)
		default:
			http.Error(w, `{"error":"Метод не е дозволен"}`, http.StatusMethodNotAllowed)
		}
	})))

	// Apply global middleware: CORS → rate limit → router
	handler := middleware.CORS(globalRL.Middleware(mux))

	addr := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{Addr: addr, Handler: handler}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// Periodically nullify IP hashes older than 12h to minimise data retention.
	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			if err := store.NullifyOldIPHMACs(context.Background()); err != nil {
				log.Printf("nullify ip hmacs: %v", err)
			}
		}
	}()

	go func() {
		log.Printf("starting server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	<-quit
	log.Printf("shutting down.")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown: %v", err)
	}
	log.Printf("shutdown complete")
}
