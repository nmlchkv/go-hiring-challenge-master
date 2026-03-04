package repositories

import (
	"github.com/mytheresa/go-hiring-challenge/models"
	"gorm.io/gorm"
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
	base := r.db.Model(&models.Product{})
	if categoryCode != "" {
		base = base.Joins("JOIN categories ON categories.id = products.category_id").Where("categories.code = ?", categoryCode)
	}
	if priceLessThan != nil {
		base = base.Where("products.price < ?", *priceLessThan)
	}

	// count total products after filters, before pagination; protect against future JOIN duplicates
	if err := base.Distinct("products.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	itemsQuery := r.db.Model(&models.Product{}).Preload("Category")
	if categoryCode != "" {
		itemsQuery = itemsQuery.Joins("JOIN categories ON categories.id = products.category_id").Where("categories.code = ?", categoryCode)
	}
	if priceLessThan != nil {
		itemsQuery = itemsQuery.Where("products.price < ?", *priceLessThan)
	}

	if err := itemsQuery.Offset(offset).Limit(limit).Order("code").Find(&products).Error; err != nil {
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