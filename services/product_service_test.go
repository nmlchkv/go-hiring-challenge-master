package services

import (
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
)

type mockRepo struct {
	FindByCodeFunc func(code string) (*models.Product, error)
	GetProductsWithFilterFunc func(categoryCode string, priceLessThan *float64, offset, limit int) ([]models.Product, int64, error)
}

func (m *mockRepo) FindByCode(code string) (*models.Product, error) {
	return m.FindByCodeFunc(code)
}
func (m *mockRepo) GetProductsWithFilter(categoryCode string, priceLessThan *float64, offset, limit int) ([]models.Product, int64, error) {
	return m.GetProductsWithFilterFunc(categoryCode, priceLessThan, offset, limit)
}

func TestGetProductDetail(t *testing.T) {
	repo := &mockRepo{
		FindByCodeFunc: func(code string) (*models.Product, error) {
			return &models.Product{Code: code, Price: decimal.NewFromFloat(10.99)}, nil
		},
	}
	s := NewProductService(repo)
	p, err := s.GetProductDetail("PROD001")
	if err != nil || p.Code != "PROD001" {
		t.Fatalf("unexpected error or wrong product")
	}
}

func TestGetCatalog(t *testing.T) {
	repo := &mockRepo{
		GetProductsWithFilterFunc: func(categoryCode string, priceLessThan *float64, offset, limit int) ([]models.Product, int64, error) {
			return []models.Product{{Code: "PROD001"}}, 1, nil
		},
	}
	s := NewProductService(repo)
	products, total, err := s.GetCatalog("", nil, 0, 10)
	if err != nil || total != 1 || len(products) != 1 {
		t.Fatalf("unexpected error or wrong result")
	}
} 