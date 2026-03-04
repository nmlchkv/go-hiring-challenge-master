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
	product, err := s.repo.FindByCode(code)
	if err != nil {
		return nil, err
	}

	for i := range product.Variants {
		if product.Variants[i].Price.IsZero() {
			product.Variants[i].Price = product.Price
		}
	}

	return product, nil
} 