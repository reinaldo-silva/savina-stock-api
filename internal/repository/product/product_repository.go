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

func (r *GormProductRepository) GetAll(
	page int,
	pageSize int,
	nameFilter string,
	categoryIDs []uint) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	query := r.db.Preload("Categories")

	if nameFilter != "" {
		query = query.Where("name ILIKE ?", "%"+nameFilter+"%")
	}

	if len(categoryIDs) > 0 {
		query = query.Joins("JOIN product_categories pc ON pc.product_id = products.id").
			Where("pc.category_id IN ?", categoryIDs)
	}

	if err := query.Model(&domain.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("id ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *GormProductRepository) Create(p domain.Product) error {
	return r.db.Create(&p).Error
}

func (r *GormProductRepository) FindBySlug(slug string) (*domain.Product, error) {
	var product domain.Product
	result := r.db.Where("slug = ?", slug).Preload("Categories").First(&product)
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

func (r *GormProductRepository) ClearProductCategories(productID uint) error {
	return r.db.Model(&domain.Product{ID: productID}).Association("Categories").Clear()
}

func (r *GormProductRepository) UpdateProductCategories(product *domain.Product) error {
	return r.db.Model(product).Association("Categories").Replace(product.Categories)
}
