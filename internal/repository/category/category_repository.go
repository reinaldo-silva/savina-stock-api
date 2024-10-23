package category_repository

import (
	"fmt"

	"github.com/reinaldo-silva/savina-stock/internal/domain/category"
	"gorm.io/gorm"
)

type CategoryRepositoryImpl struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) category.CategoryRepository {
	return &CategoryRepositoryImpl{db: db}
}

func (repo *CategoryRepositoryImpl) Create(category *category.Category) error {
	return repo.db.Create(category).Error
}

func (repo *CategoryRepositoryImpl) GetAll() ([]category.Category, error) {
	var categories []category.Category

	err := repo.db.Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepositoryImpl) FindById(id uint) (*category.Category, error) {
	var category category.Category
	if err := r.db.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("category not found")
		}
		return nil, err
	}
	return &category, nil
}
