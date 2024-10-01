package usecase_category

import (
	"fmt"

	"github.com/reinaldo-silva/savina-stock/internal/domain/category"
)

type CategoryUseCase struct {
	repo category.CategoryRepository
}

func NewCategoryUseCase(repo category.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{repo: repo}
}

func (uc *CategoryUseCase) CreateCategory(category *category.Category) (*category.Category, error) {

	if category.Name == "" {
		return nil, fmt.Errorf("category name is required")
	}

	err := uc.repo.Create(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}
