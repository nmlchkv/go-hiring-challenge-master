package catalog

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
)

type mockCatalogService struct {
	errToReturn error
}

func (m *mockCatalogService) GetCatalog(categoryCode string, priceLessThan *float64, offset, limit int) ([]models.Product, int64, error) {
	if m.errToReturn != nil {
		return nil, 0, m.errToReturn
	}
	products := []models.Product{
		{Code: "PROD001", Price: decimal.NewFromFloat(10.99), Category: models.Category{Code: "CLOTHING", Name: "Clothing"}},
		{Code: "PROD002", Price: decimal.NewFromFloat(12.49), Category: models.Category{Code: "SHOES", Name: "Shoes"}},
	}
	return products, 2, nil
}

func TestHandleGet_OK(t *testing.T) {
	h := NewCatalogHandler(&mockCatalogService{})
	r := httptest.NewRequest("GET", "/catalog", nil)
	w := httptest.NewRecorder()
	h.HandleGet(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var payload struct {
		Data Response `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if payload.Data.Total != 2 {
		t.Fatalf("expected total 2, got %d", payload.Data.Total)
	}
	if len(payload.Data.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(payload.Data.Items))
	}
}

func TestHandleGet_ValidationErrors(t *testing.T) {
	h := NewCatalogHandler(&mockCatalogService{})

	tests := []struct {
		name       string
		url        string
		wantCode   int
		wantErrMsg string
	}{
		{"limit zero", "/catalog?limit=0", http.StatusBadRequest, "limit must be in [1; 100]"},
		{"limit too big", "/catalog?limit=101", http.StatusBadRequest, "limit must be in [1; 100]"},
		{"offset negative", "/catalog?offset=-1", http.StatusBadRequest, "invalid offset"},
		{"price zero", "/catalog?priceLessThan=0", http.StatusBadRequest, "priceLessThan must be > 0"},
		{"price invalid", "/catalog?priceLessThan=aaa", http.StatusBadRequest, "invalid priceLessThan"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()
			h.HandleGet(w, r)
			if w.Code != tt.wantCode {
				t.Fatalf("expected %d, got %d", tt.wantCode, w.Code)
			}
			var payload map[string]string
			if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
				t.Fatalf("failed to unmarshal error response: %v", err)
			}
			if payload["error"] != tt.wantErrMsg {
				t.Fatalf("expected error %q, got %q", tt.wantErrMsg, payload["error"])
			}
		})
	}
}

func TestHandleGet_InternalError(t *testing.T) {
	h := NewCatalogHandler(&mockCatalogService{errToReturn: errors.New("db error")})
	r := httptest.NewRequest("GET", "/catalog", nil)
	w := httptest.NewRecorder()
	h.HandleGet(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}

	var payload map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}
	if payload["error"] != "internal server error" {
		t.Fatalf("expected error %q, got %q", "internal server error", payload["error"])
	}
}