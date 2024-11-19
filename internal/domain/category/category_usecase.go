package category

import (
	"fmt"
)

type CategoryUseCase struct {
	repo CategoryRepository
}

func NewCategoryUseCase(repo CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{repo: repo}
}

func (uc *CategoryUseCase) CreateCategory(category *Category) (*Category, error) {
	if category.Name == "" {
		return nil, fmt.Errorf("category name is required")
	}

	err := uc.repo.Create(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (uc *CategoryUseCase) GetAllCategories() ([]Category, error) {
	categories, err := uc.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (uc *CategoryUseCase) GetCategoryByID(id uint) (*Category, error) {
	category, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("category not found")
	}
	return category, nil
}

func (uc *CategoryUseCase) UpdateCategory(updatedCategory *Category) error {
	existingCategory, err := uc.repo.GetByID(updatedCategory.ID)
	if err != nil {
		return fmt.Errorf("category not found")
	}

	if updatedCategory.Name == "" {
		return fmt.Errorf("category name is required")
	}

	existingCategory.Name = updatedCategory.Name
	err = uc.repo.Update(existingCategory)
	if err != nil {
		return err
	}
	return nil
}

func (uc *CategoryUseCase) DeleteCategory(id uint) error {
	_, err := uc.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("category not found")
	}

	hasProducts, err := uc.repo.HasProducts(id)
	if err != nil {
		return fmt.Errorf("error checking related products: %v", err)
	}
	if hasProducts {
		return fmt.Errorf("cannot delete category, products are associated with it")
	}

	err = uc.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
