package category_repository

import (
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
