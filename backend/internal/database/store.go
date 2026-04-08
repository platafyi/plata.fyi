package database

import (
	"context"
	"database/sql"
	"time"
)

type Store interface {
	Ping(ctx context.Context) error

	GetIndustries(ctx context.Context) ([]Industry, error)
	GetCities(ctx context.Context) ([]City, error)

	InsertToken(ctx context.Context, token string) error
	GetOwnerByToken(ctx context.Context, token string) (*string, error)
	DeleteToken(ctx context.Context, token string) error
	ExchangeToken(ctx context.Context, magicToken, sessionToken string) (*string, error)

	CreateSubmission(ctx context.Context, inp CreateSubmissionInput) (*SalarySubmission, error)
	GetSubmissionsByOwner(ctx context.Context, ownerID string) ([]SalarySubmission, error)
	GetSubmissionByID(ctx context.Context, id string) (*SalarySubmission, error)
	UpdateSubmission(ctx context.Context, id, ownerID string, inp CreateSubmissionInput) error
	DeleteSubmission(ctx context.Context, id, ownerID string) error
	CountRecentSubmissionsByIPHMAC(ctx context.Context, ipHMAC string, since time.Time) (int, error)
	NullifyOldIPHMACs(ctx context.Context) error

	SearchSalaries(ctx context.Context, f SearchFilters) ([]SalarySubmission, int, error)
	GetSalaryStats(ctx context.Context, groupBy string, f SearchFilters) ([]SalaryStats, error)

	SearchCompanies(ctx context.Context, q string) ([]Company, error)
	SearchJobTitles(ctx context.Context, q string) ([]string, error)
}

type PostgresStore struct{ db *sql.DB }

func NewPostgresStore(db *sql.DB) *PostgresStore { return &PostgresStore{db: db} }

func (s *PostgresStore) Ping(ctx context.Context) error { return s.db.PingContext(ctx) }

// Compiler check — fails at build time if PostgresStore is missing any method.
var _ Store = (*PostgresStore)(nil)
