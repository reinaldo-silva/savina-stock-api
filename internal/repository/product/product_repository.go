package product_repository

import (
	"errors"
	"time"

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

func (r *GormProductRepository) DeleteBySlug(slug string) error {
	result := r.db.Where("slug = ?", slug).Delete(&domain.Product{})
	if result.RowsAffected == 0 {
		return errors.New("product not found")
	}
	return result.Error
}

func (r *GormProductRepository) UpdateBySlug(slug string, updatedProduct domain.Product) (domain.Product, error) {
	var existingProduct domain.Product

	if err := r.db.Where("slug = ?", slug).First(&existingProduct).Error; err != nil {
		return existingProduct, err
	}

	existingProduct.Name = updatedProduct.Name
	existingProduct.Description = updatedProduct.Description
	existingProduct.Price = updatedProduct.Price
	existingProduct.Cost = updatedProduct.Cost
	existingProduct.Stock = updatedProduct.Stock
	existingProduct.UpdatedAt = time.Now()

	if err := r.db.Save(&existingProduct).Error; err != nil {
		return existingProduct, err
	}

	return existingProduct, nil
}
