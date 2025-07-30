package models

import (
	"gorm.io/gorm"
)

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

// GetProductsWithFilter returns products with pagination and filters applied.
func (r *ProductsRepository) GetProductsWithFilter(categoryCode string, priceLessThan *float64, offset, limit int) ([]Product, int64, error) {
	var products []Product
	var total int64
	db := r.db.Model(&Product{}).Preload("Variants").Preload("Category")
	if categoryCode != "" {
		db = db.Joins("JOIN categories ON categories.id = products.category_id").Where("categories.code = ?", categoryCode)
	}
	if priceLessThan != nil {
		db = db.Where("products.price < ?", *priceLessThan)
	}
	db.Count(&total)
	if err := db.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (r *ProductsRepository) DB() *gorm.DB {
	return r.db
}

// ProductsRepositoryMock implements the interface for testing purposes.
// (FindByCode is needed for the variants handler)
type ProductsRepositoryMock struct {
	FindByCodeFunc func(code string) (*Product, error)
}

func (m *ProductsRepositoryMock) DB() *gorm.DB {
	return nil
}

func (m *ProductsRepositoryMock) GetProductsWithFilter(categoryCode string, priceLessThan *float64, offset, limit int) ([]Product, int64, error) {
	return nil, 0, nil
}

func (m *ProductsRepositoryMock) FindByCode(code string) (*Product, error) {
	return m.FindByCodeFunc(code)
}
