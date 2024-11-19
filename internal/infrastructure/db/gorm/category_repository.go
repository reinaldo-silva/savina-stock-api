package gorm

import (
	"fmt"

	"github.com/reinaldo-silva/savina-stock/internal/domain/category"
	"gorm.io/gorm"
)

type GormCategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) category.CategoryRepository {
	return &GormCategoryRepository{db: db}
}

func (repo *GormCategoryRepository) Create(category *category.Category) error {
	return repo.db.Create(category).Error
}

func (repo *GormCategoryRepository) GetAll() ([]category.Category, error) {
	var categories []category.Category
	err := repo.db.Order("id ASC").Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (repo *GormCategoryRepository) GetByID(id uint) (*category.Category, error) {
	var category category.Category
	if err := repo.db.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("category not found")
		}
		return nil, err
	}
	return &category, nil
}

func (repo *GormCategoryRepository) Update(category *category.Category) error {
	return repo.db.Save(category).Error
}

func (repo *GormCategoryRepository) Delete(id uint) error {
	if err := repo.db.Delete(&category.Category{}, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("category not found")
		}
		return err
	}
	return nil
}

func (repo *GormCategoryRepository) HasProducts(categoryID uint) (bool, error) {
	var count int64
	err := repo.db.Table("product_categories").Where("category_id = ?", categoryID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
