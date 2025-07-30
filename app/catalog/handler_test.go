package catalog

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/mytheresa/go-hiring-challenge/services"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
)

type mockService struct{}

func (m *mockService) GetCatalog(categoryCode string, priceLessThan *float64, offset, limit int) ([]models.Product, int64, error) {
	products := []models.Product{
		{Code: "PROD001", Price: decimal.NewFromFloat(10.99), Category: models.Category{Name: "Clothing"}},
		{Code: "PROD002", Price: decimal.NewFromFloat(12.49), Category: models.Category{Name: "Shoes"}},
	}
	return products, 2, nil
}

func TestHandleGet_OK(t *testing.T) {
	h := NewCatalogHandler(&mockService{})
	r := httptest.NewRequest("GET", "/catalog", nil)
	w := httptest.NewRecorder()
	h.HandleGet(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
} 