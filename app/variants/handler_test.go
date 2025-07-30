package variants

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/mytheresa/go-hiring-challenge/services"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type mockService struct {
	FindByCodeFunc func(code string) (*models.Product, error)
}

func (m *mockService) GetProductDetail(code string) (*models.Product, error) {
	return m.FindByCodeFunc(code)
}

func TestHandleGet_ProductFound(t *testing.T) {
	h := NewHandler(&mockService{
		FindByCodeFunc: func(code string) (*models.Product, error) {
			return &models.Product{
				Code:  "PROD001",
				Price: decimal.NewFromFloat(10.99),
				Category: models.Category{Name: "Clothing"},
				Variants: []models.Variant{{Name: "Variant A", SKU: "SKU001A", Price: decimal.NewFromFloat(0)}},
			}, nil
		},
	})
	r := httptest.NewRequest("GET", "/catalog/PROD001", nil)
	w := httptest.NewRecorder()
	h.HandleGet(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestHandleGet_ProductNotFound(t *testing.T) {
	h := NewHandler(&mockService{
		FindByCodeFunc: func(code string) (*models.Product, error) {
			return nil, gorm.ErrRecordNotFound
		},
	})
	r := httptest.NewRequest("GET", "/catalog/UNKNOWN", nil)
	w := httptest.NewRecorder()
	h.HandleGet(w, r)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
} 