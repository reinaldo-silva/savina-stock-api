package repository

import (
	domain "github.com/reinaldo-silva/savina-stock/internal/domain/product"
	"gorm.io/gorm"
)

type GormProductRepository struct {
	db *gorm.DB
}

func NewGormProductRepository(db *gorm.DB) *GormProductRepository {
	return &GormProductRepository{db: db}
}

func (r *GormProductRepository) GetAll() ([]domain.Product, error) {
	var products []domain.Product
	result := r.db.Find(&products)
	return products, result.Error
}

func (r *GormProductRepository) Create(p domain.Product) error {
	return r.db.Create(&p).Error
}

func (r *GormProductRepository) FindBySlug(slug string) (*domain.Product, error) {
	var product domain.Product
	result := r.db.Where("slug = ?", slug).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}
