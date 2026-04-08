package handlers

import (
	"context"
	"time"

	"github.com/platafyi/plata.fyi/internal/database"
)

type MockStore struct {
	PingErr error

	Industries    []database.Industry
	IndustriesErr error

	Cities    []database.City
	CitiesErr error

	InsertTokenErr   error
	OwnerByToken     *string
	OwnerByTokenErr  error
	DeleteTokenErr   error
	ExchangedOwnerID *string
	ExchangeTokenErr error

	Submissions       []database.SalarySubmission
	SubmissionsErr    error
	CreatedSubmission *database.SalarySubmission
	CreateErr         error
	SubmissionByID    *database.SalarySubmission
	SubmissionByIDErr error
	UpdateErr         error
	DeleteErr         error
	IPHMACCount       int
	IPHMACErr         error

	SearchResults []database.SalarySubmission
	SearchTotal   int
	SearchErr     error
	StatsResults  []database.SalaryStats
	StatsErr      error

	Companies    []database.Company
	CompaniesErr error
	JobTitles    []string
	JobTitlesErr error
}

func (m *MockStore) Ping(_ context.Context) error { return m.PingErr }

func (m *MockStore) GetIndustries(_ context.Context) ([]database.Industry, error) {
	return m.Industries, m.IndustriesErr
}

func (m *MockStore) GetCities(_ context.Context) ([]database.City, error) {
	return m.Cities, m.CitiesErr
}

func (m *MockStore) InsertToken(_ context.Context, _ string) error { return m.InsertTokenErr }

func (m *MockStore) GetOwnerByToken(_ context.Context, _ string) (*string, error) {
	return m.OwnerByToken, m.OwnerByTokenErr
}

func (m *MockStore) DeleteToken(_ context.Context, _ string) error { return m.DeleteTokenErr }

func (m *MockStore) ExchangeToken(_ context.Context, _, _ string) (*string, error) {
	return m.ExchangedOwnerID, m.ExchangeTokenErr
}

func (m *MockStore) CreateSubmission(_ context.Context, _ database.CreateSubmissionInput) (*database.SalarySubmission, error) {
	return m.CreatedSubmission, m.CreateErr
}

func (m *MockStore) GetSubmissionsByOwner(_ context.Context, _ string) ([]database.SalarySubmission, error) {
	return m.Submissions, m.SubmissionsErr
}

func (m *MockStore) GetSubmissionByID(_ context.Context, _ string) (*database.SalarySubmission, error) {
	return m.SubmissionByID, m.SubmissionByIDErr
}

func (m *MockStore) UpdateSubmission(_ context.Context, _, _ string, _ database.CreateSubmissionInput) error {
	return m.UpdateErr
}

func (m *MockStore) DeleteSubmission(_ context.Context, _, _ string) error { return m.DeleteErr }

func (m *MockStore) CountRecentSubmissionsByIPHMAC(_ context.Context, _ string, _ time.Time) (int, error) {
	return m.IPHMACCount, m.IPHMACErr
}

func (m *MockStore) NullifyOldIPHMACs(_ context.Context) error { return nil }

func (m *MockStore) SearchSalaries(_ context.Context, _ database.SearchFilters) ([]database.SalarySubmission, int, error) {
	return m.SearchResults, m.SearchTotal, m.SearchErr
}

func (m *MockStore) GetSalaryStats(_ context.Context, _ string, _ database.SearchFilters) ([]database.SalaryStats, error) {
	return m.StatsResults, m.StatsErr
}

func (m *MockStore) SearchCompanies(_ context.Context, _ string) ([]database.Company, error) {
	return m.Companies, m.CompaniesErr
}

func (m *MockStore) SearchJobTitles(_ context.Context, _ string) ([]string, error) {
	return m.JobTitles, m.JobTitlesErr
}
