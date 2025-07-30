package services

import (
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/mytheresa/go-hiring-challenge/repositories"
)

type ProductService struct {
	repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetCatalog(categoryCode string, priceLessThan *float64, offset, limit int) ([]models.Product, int64, error) {
	return s.repo.GetProductsWithFilter(categoryCode, priceLessThan, offset, limit)
}

func (s *ProductService) GetProductDetail(code string) (*models.Product, error) {
	return s.repo.FindByCode(code)
} 