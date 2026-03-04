package variants

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type mockDetailService struct {
	GetProductDetailFunc func(code string) (*models.Product, error)
}

func (m *mockDetailService) GetProductDetail(code string) (*models.Product, error) {
	return m.GetProductDetailFunc(code)
}

func TestHandleGet_ProductFound(t *testing.T) {
	h := NewHandler(&mockDetailService{
		GetProductDetailFunc: func(code string) (*models.Product, error) {
			return &models.Product{
				Code:  "PROD001",
				Price: decimal.NewFromFloat(10.99),
				Category: models.Category{Code: "CLOTHING", Name: "Clothing"},
				Variants: []models.Variant{
					{ID: 1, Name: "Variant A", SKU: "SKU001A", Price: decimal.NewFromFloat(10.99)},
				},
			}, nil
		},
	})

	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog/{id}", h.HandleGet)

	r := httptest.NewRequest("GET", "/catalog/PROD001", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var payload struct {
		Data ProductDetail `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if payload.Data.Code != "PROD001" {
		t.Fatalf("expected code PROD001, got %s", payload.Data.Code)
	}
	if payload.Data.Category.Code != "CLOTHING" {
		t.Fatalf("expected category code CLOTHING, got %s", payload.Data.Category.Code)
	}
	if len(payload.Data.Variants) != 1 {
		t.Fatalf("expected 1 variant, got %d", len(payload.Data.Variants))
	}
}

func TestHandleGet_ProductNotFound(t *testing.T) {
	h := NewHandler(&mockDetailService{
		GetProductDetailFunc: func(code string) (*models.Product, error) {
			return nil, gorm.ErrRecordNotFound
		},
	})
	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog/{id}", h.HandleGet)

	r := httptest.NewRequest("GET", "/catalog/UNKNOWN", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
} 

func TestHandleGet_MissingID(t *testing.T) {
	h := NewHandler(&mockDetailService{
		GetProductDetailFunc: func(code string) (*models.Product, error) {
			return nil, nil
		},
	})
	r := httptest.NewRequest("GET", "/catalog", nil)
	w := httptest.NewRecorder()
	h.HandleGet(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	var payload map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}
	if payload["error"] != "missing id" {
		t.Fatalf("expected error %q, got %q", "missing id", payload["error"])
	}
}