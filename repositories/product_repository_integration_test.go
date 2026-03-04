//go:build integration

package repositories

import (
	"os"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/database"
)

// newTestRepo assumes that the database is up and seeded
// (e.g. via `make docker-up` + `make seed`) before running tests.
func newTestRepo(t *testing.T) *GormProductRepository {
	t.Helper()

	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	if user == "" || pass == "" || dbname == "" || port == "" {
		t.Skip("database env vars not set, skipping repository integration tests")
	}

	db, closeFn := database.New(user, pass, dbname, port)
	t.Cleanup(func() { _ = closeFn() })

	return NewGormProductRepository(db)
}

func TestGetProductsWithFilter_TotalIndependentOfPagination(t *testing.T) {
	repo := newTestRepo(t)

	_, total1, err := repo.GetProductsWithFilter("", nil, 0, 2)
	if err != nil {
		t.Fatalf("unexpected error on first page: %v", err)
	}
	if total1 == 0 {
		t.Fatalf("expected total > 0")
	}

	_, total2, err := repo.GetProductsWithFilter("", nil, 2, 2)
	if err != nil {
		t.Fatalf("unexpected error on second page: %v", err)
	}
	if total1 != total2 {
		t.Fatalf("expected total to be independent of offset/limit, got %d and %d", total1, total2)
	}
}

func TestGetProductsWithFilter_CategoryFilterAffectsTotal(t *testing.T) {
	repo := newTestRepo(t)

	_, totalAll, err := repo.GetProductsWithFilter("", nil, 0, 100)
	if err != nil {
		t.Fatalf("unexpected error for all products: %v", err)
	}
	if totalAll == 0 {
		t.Fatalf("expected totalAll > 0")
	}

	_, totalShoes, err := repo.GetProductsWithFilter("shoes", nil, 0, 100)
	if err != nil {
		t.Fatalf("unexpected error for shoes category: %v", err)
	}
	if totalShoes == 0 {
		t.Fatalf("expected some products in shoes category")
	}
	if totalShoes >= totalAll {
		t.Fatalf("expected totalShoes < totalAll, got %d vs %d", totalShoes, totalAll)
	}
}

func TestGetProductsWithFilter_PriceLessThanFilterAndCategoryPreload(t *testing.T) {
	repo := newTestRepo(t)

	priceLimit := 10.0
	products, total, err := repo.GetProductsWithFilter("", &priceLimit, 0, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total == 0 {
		t.Fatalf("expected some products with price < %v", priceLimit)
	}
	if len(products) == 0 {
		t.Fatalf("expected at least one product in items slice")
	}

	for _, p := range products {
		if p.Price.InexactFloat64() >= priceLimit {
			t.Fatalf("product %s has price %.2f >= %.2f", p.Code, p.Price.InexactFloat64(), priceLimit)
		}
		if p.Category.Code == "" || p.Category.Name == "" {
			t.Fatalf("expected category to be preloaded for product %s, got %#v", p.Code, p.Category)
		}
	}
}

