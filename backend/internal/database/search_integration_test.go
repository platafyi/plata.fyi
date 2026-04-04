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

// TestSearchSalariesIsRandom inserts 10 submissions then calls SearchSalaries
// 20 times. It asserts that at least two calls return a different order
// proving the query uses ORDER BY random(), not ORDER BY created_at or any
// other deterministic column.
func TestSearchSalariesIsRandom(t *testing.T) {
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

	firstOrder := salaryIDs(t, store, ctx)
	if len(firstOrder) < 10 {
		t.Fatalf("expected at least 10 results, got %d", len(firstOrder))
	}

	const attempts = 20
	for i := 0; i < attempts; i++ {
		if !sameOrder(firstOrder, salaryIDs(t, store, ctx)) {
			return // different order observed, ORDER BY random() confirmed
		}
	}

	t.Errorf("SearchSalaries returned the same order all %d times, not using ORDER BY random()", attempts)
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
