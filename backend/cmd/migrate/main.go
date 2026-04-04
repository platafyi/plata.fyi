package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://platafyi:platafyi@localhost:5432/platafyi?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	// Create migrations tracking table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		filename  VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`)
	if err != nil {
		log.Fatalf("create schema_migrations: %v", err)
	}

	migrationsDir := "migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		// Try relative to binary location
		exe, _ := os.Executable()
		migrationsDir = filepath.Join(filepath.Dir(exe), "../migrations")
	}

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("read migrations dir: %v", err)
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, f := range files {
		var count int
		err = db.QueryRow(`SELECT COUNT(*) FROM schema_migrations WHERE filename = $1`, f).Scan(&count)
		if err != nil {
			log.Fatalf("check migration %s: %v", f, err)
		}
		if count > 0 {
			fmt.Printf("skip (already applied): %s\n", f)
			continue
		}

		content, err := os.ReadFile(filepath.Join(migrationsDir, f))
		if err != nil {
			log.Fatalf("read %s: %v", f, err)
		}

		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("begin tx for %s: %v", f, err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			log.Fatalf("apply %s: %v", f, err)
		}

		if _, err := tx.Exec(`INSERT INTO schema_migrations (filename) VALUES ($1)`, f); err != nil {
			tx.Rollback()
			log.Fatalf("record %s: %v", f, err)
		}

		if err := tx.Commit(); err != nil {
			log.Fatalf("commit %s: %v", f, err)
		}

		fmt.Printf("applied: %s\n", f)
	}

	fmt.Println("migrations complete")
}
