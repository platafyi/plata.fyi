package database_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/platafyi/plata.fyi/internal/database"
)

func testDSN() string {
	if v := os.Getenv("TEST_DB_URL"); v != "" {
		return v
	}
	return "postgres://platafyi:platafyi@localhost:5433/platafyi?sslmode=disable"
}

func testRawDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := database.New(testDSN())
	if err != nil {
		t.Skipf("DB not available: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Skipf("DB not reachable: %v", err)
	}
	return db
}

func testStore(t *testing.T) database.Store {
	t.Helper()
	return database.NewPostgresStore(testRawDB(t))
}

// TestSearchSalariesLatestFirst inserts 10 submissions then verifies that
// SearchSalaries returns them in a stable, consistent order (ORDER BY ctid DESC),
// and that the last inserted submission appears first.
func TestSearchSalariesLatestFirst(t *testing.T) {
	store := testStore(t)
	ctx := context.Background()

	ownerID := "00000000-0000-0000-0000-000000000001"

	industries, err := store.GetIndustries(ctx)
	if err != nil || len(industries) == 0 {
		t.Skip("no industries seeded — run migrations first")
	}
	cities, err := store.GetCities(ctx)
	if err != nil || len(cities) == 0 {
		t.Skip("no cities seeded — run migrations first")
	}
	industryID := industries[0].ID
	cityID := cities[0].ID

	var createdIDs []string
	t.Cleanup(func() {
		for _, id := range createdIDs {
			store.DeleteSubmission(ctx, id, ownerID) //nolint
		}
	})

	for i := 0; i < 10; i++ {
		sub, err := store.CreateSubmission(ctx, database.CreateSubmissionInput{
			OwnerID:         ownerID,
			CompanyName:     fmt.Sprintf("Company %d", i),
			JobTitle:        fmt.Sprintf("Engineer %d", i),
			IndustryID:      industryID,
			CityID:          cityID,
			Seniority:       "mid",
			WorkArrangement: "office",
			EmploymentType:  "full_time",
			BaseSalary:      50000 + i*1000,
			SalaryYear:      2024,
		})
		if err != nil {
			t.Fatalf("insert submission %d: %v", i, err)
		}
		createdIDs = append(createdIDs, sub.ID)
	}

	lastInsertedID := createdIDs[len(createdIDs)-1]

	firstOrder := salaryIDs(t, store, ctx)
	if len(firstOrder) < 10 {
		t.Fatalf("expected at least 10 results, got %d", len(firstOrder))
	}

	// Order must be stable across calls.
	for i := 0; i < 5; i++ {
		if !sameOrder(firstOrder, salaryIDs(t, store, ctx)) {
			t.Errorf("SearchSalaries returned a different order on call %d — expected stable ctid ordering", i+1)
		}
	}

	// The last inserted row must appear first (highest ctid).
	if firstOrder[0] != lastInsertedID {
		t.Errorf("expected last inserted submission %s to be first in results, got %s", lastInsertedID, firstOrder[0])
	}
}

func salaryIDs(t *testing.T, store database.Store, ctx context.Context) []string {
	t.Helper()
	subs, _, err := store.SearchSalaries(ctx, database.SearchFilters{PageSize: 100})
	if err != nil {
		t.Fatalf("SearchSalaries: %v", err)
	}
	ids := make([]string, len(subs))
	for i, s := range subs {
		ids[i] = s.ID
	}
	return ids
}

// TestSubmissionTimestampPrecision verifies that created_at and updated_at are
// stored as TIMESTAMPTZ, not DATE. Checks the schema directly to avoid non deterministic test.
func TestSubmissionTimestampPrecision(t *testing.T) {
	db := testRawDB(t)
	ctx := context.Background()

	for _, col := range []string{"created_at", "updated_at"} {
		var dataType string
		err := db.QueryRowContext(ctx,
			`SELECT data_type FROM information_schema.columns
			 WHERE table_name = 'salary_submissions' AND column_name = $1`,
			col,
		).Scan(&dataType)
		if err != nil {
			t.Fatalf("query column type for %s: %v", col, err)
		}
		if dataType != "timestamp with time zone" {
			t.Errorf("column %s type = %q, want 'timestamp with time zone' — run migration 018", col, dataType)
		}
	}
}

func sameOrder(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
