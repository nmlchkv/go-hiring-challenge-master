package repositories

import (
	"github.com/mytheresa/go-hiring-challenge/models"
)

type ProductRepository interface {
	GetProductsWithFilter(categoryCode string, priceLessThan *float64, offset, limit int) ([]models.Product, int64, error)
	FindByCode(code string) (*models.Product, error)
}

type GormProductRepository struct {
	db *gorm.DB
}

func NewGormProductRepository(db *gorm.DB) *GormProductRepository {
	return &GormProductRepository{db: db}
}

func (r *GormProductRepository) GetProductsWithFilter(categoryCode string, priceLessThan *float64, offset, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64
	db := r.db.Model(&models.Product{}).Preload("Variants").Preload("Category")
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

func (r *GormProductRepository) FindByCode(code string) (*models.Product, error) {
	var product models.Product
	if err := r.db.Preload("Variants").Preload("Category").Where("code = ?", code).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
} 