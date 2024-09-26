package repository

import (
	domain "github.com/reinaldo-silva/savina-stock/internal/domain/product"
	"gorm.io/gorm"
)

type GormProductRepository struct {
	db *gorm.DB
}

func NewGormProductRepository(db *gorm.DB) domain.ProductRepository {
	return &GormProductRepository{db}
}

func (r *GormProductRepository) GetAll() ([]domain.Product, error) {
	var products []domain.Product
	result := r.db.Find(&products)
	return products, result.Error
}

func (r *GormProductRepository) Create(product domain.Product) error {
	result := r.db.Create(&product)
	return result.Error
}
