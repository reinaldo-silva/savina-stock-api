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
